/*
 * @Author: rayou
 * @Date: 2019-04-21 11:31:53
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-22 01:30:54
 */
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
	"github.com/AzuresYang/arx7/config"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "arxdepl"
	app.Usage = "arx7 command-line app by azureyang"
	app.Author = "AzureYang"
	app.Email = "AzureYang@xxx.com"
	app.Commands = []cli.Command{
		// 开始
		{
			Name:      "start",
			ShortName: "start",
			Usage:     "start spider task",
			Action:    startSpider,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "spidername",
					Value: "default-spider",
					Usage: "spider task name",
				},
				cli.StringFlag{
					Name:  "config",
					Value: "spider.json",
					Usage: "spider config json file",
				},
				cli.StringFlag{
					Name:  "port",
					Value: "31001",
					Usage: "arx master port",
				},
			},
		},
		{
			// 生成默认配置
			Name:      "genconf",
			ShortName: "g",
			Usage:     "generate config file",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "file",
					Value: "spider.json",
					Usage: "config file name",
				},
				cli.StringFlag{
					Name:  "dir",
					Value: "./",
					Usage: "path to config file",
				},
			},
			Action: createDefaultConf,
		},
		// 部署
		{
			Name:      "deployment",
			ShortName: "dep",
			Usage:     "deployment spider",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "spidername",
					Value: "default-spider",
					Usage: " spider task name",
				},
				cli.StringFlag{
					Name:  "image",
					Value: "",
					Usage: "spider image",
				},
				cli.Uint64Flag{
					Name:  "nodes",
					Value: 1,
					Usage: "num of the nodes",
				},
			},
			Action: deploymentSpider,
		},
		// 删除
		{
			Name:      "delete",
			ShortName: "d",
			Usage:     "delete spider task",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "spidername",
					Value: "default-spider",
					Usage: "spider task name",
				},
			},
			Action: deploymentSpider,
		},
		// 停止
		{
			Name:      "stop",
			ShortName: "st",
			Usage:     "stop spider",
			Action:    stopSpider,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "spidername",
					Value: "default-spider",
					Usage: "spider task name",
				},
			},
		},
		// 获取spider状态
		{
			Name:   "status",
			Usage:  "get spider status",
			Action: getSpiderStatus,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "spidername",
					Value: "default-spider",
					Usage: "spider task name",
				},
			},
		},
		// 获取spider状态
		{
			Name:   "pod",
			Usage:  "get pods",
			Action: getPod,
		},
		// echo
		{
			Name:      "echo",
			ShortName: "e",
			Usage:     "echo spider",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "spidername",
					Value: "default-spider",
					Usage: "spider task name",
				},
				cli.Uint64Flag{
					Name:  "model",
					Value: 1,
					Usage: "1 is echo, 2 is echo redis",
				},
			},
			Action: deploymentSpider,
		},
	}
	app.Run(os.Args)
}

// 从kuberntet中获取service关联到的节点
func getSpiderNodes(svc_name string) []string {
	nodes := []string{"127.0.0.1:31002"}
	return nodes
}

func startSpider(ctx *cli.Context) {
	file := ctx.String("config")
	conf := &config.SpiderStartConfig{}
	err := config.ReadConfigFromFileJson(file, conf)
	if err != nil {
		fmt.Printf("read config file fail:[%s]%s\n", file, err.Error())
		return
	}
	err = initRedis(conf)
	if err != nil {
		fmt.Printf("redis[%s] init fail:%s\n", conf.TaskConf.RedisAddr, err.Error())
		return
	}
	spider_name := ctx.String("spidername")
	nodes := getSpiderNodes(spider_name)
	if len(nodes) <= 0 {
		fmt.Printf("not found spider nodes.please ensure had deployed spider task:%s\n", conf.TaskConf.TaskName)
		return
	}
	fmt.Printf("config:%+v\n", conf)
}

func initRedis(conf *config.SpiderStartConfig) error {
	fmt.Println("init redis.....")
	reqMgr := request.NewRequestManager()
	err := reqMgr.Init(&conf.TaskConf)
	if err != nil {
		return errors.New("Connect redis fail")
	}
	err = reqMgr.ClearRedis(conf.TaskConf.RedisAddr, conf.TaskConf.RedisPassword)
	if err != nil {
		return errors.New("init redis fail")
	}
	reqs := []string{
		"http://www.xbiquge.la/paihangbang/",
	}
	for _, s := range reqs {
		req := request.NewArxRequest(s)
		req.ProcerName = conf.ProcerName
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")
		reqMgr.AddNeedGrabRequest(req)
	}
	fmt.Printf("init redis succ!\n")
	return nil
}

func sendStartToMaster(conf *config.SpiderStartConfig, nodes []string) error {
	fmt.Println("sending start info to arxmaster......")
	start_info := message.SpiderStartMsg{}
	start_info.NodeAddrs = nodes
	send_bytes, _ := json.Marshal(start_info)
	masterAddr := "127.0.0.1:" + conf.TaskConf.MasterListenPort
	err, cnn := arxlet.SendTcpMsgTimeoutWithConn(message.MSG_ARXCMD_START_SPIDER, send_bytes, masterAddr, 10*time.Second)
	if err != nil {
		return errors.New(fmt.Sprintf("connect master fail:%s", err.Error()))
	}
	rerr, resp_msg := arxlet.ParseResponseFromConn(cnn)
	if rerr != nil {
		return errors.New(fmt.Sprintf("get master response fail:%s", rerr.Error()))
	}
	fmt.Println("start info:")
	fmt.Printf("succ:%d, fail:%d\n", uint32(len(nodes))-resp_msg.Status, resp_msg.Status)
	var ret map[string]string
	err = json.Unmarshal(resp_msg.Data, &ret)
	if err != nil {
		fmt.Println("Unmarshal master resp data fail")
		return nil
	}
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}
	fmt.Println("start spider down")
	return nil
}

// 命令行测试
func getPod(ctx *cli.Context) {
	cmd := "kubectl get pods"
	ret, err := exec_shell(cmd)
	checkErr(err, "get pods")
	fmt.Println(ret)
}

// 发布程序
func deploymentSpider(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	image := ctx.String("image")
	node := ctx.Uint64("nodes")
	fmt.Printf("%s, %s,%d", spider_name, image, node)
}

// 删除Spider
func deleteSpider(ctx *cli.Context) {

}

//  停止 stop
func stopSpider(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	nodes := getSpiderNodes(spider_name)
	ret := SendMessageToSpider(nodes, message.MSG_REQ_STOP_SPIDER, []byte(""), "stop spider")
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}
}

// 获取spider状态
func getSpiderStatus(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	nodes := getSpiderNodes(spider_name)
	ret := SendMessageToSpider(nodes, message.MSG_REQ_GET_SPIDER_INFO, []byte(""), "get spider status")
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}
}

// echo spider
func echoSpider(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	nodes := getSpiderNodes(spider_name)
	cmd := message.MSG_REG_ECHO
	if ctx.Uint64("model") == 2 {
		cmd = message.MSG_REG_ECHO_REDIS
	}
	ret := SendMessageToSpider(nodes, cmd, []byte(""), "echo spider")
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}
}

// 向spider发送操作信息
func SendMessageToSpider(nodes []string, cmd uint32, data []byte, comment string) map[string]string {
	fmt.Printf("start %s spider....\n", comment)
	// message.MSG_REQ_STOP_SPIDER
	var wg sync.WaitGroup
	var resp_result map[string]string = make(map[string]string, len(nodes))
	for _, addr := range nodes {
		wg.Add(1)
		go func(addr string) {
			defer wg.Done()
			fmt.Printf("node[%s], %s......\n", addr, comment)
			err, spider_conn := arxlet.SendTcpMsgTimeoutWithConn(cmd, data, addr, 2*time.Second)
			if err != nil {
				resp_result[addr] = "connect to spider fail:" + err.Error()
				return
			}
			var resp *message.ResponseMsg
			err, resp = arxlet.ParseResponseFromConn(spider_conn)
			if err != nil {
				fmt.Printf("node[%s] parseResponse Error:%s\n", addr, err.Error())
				err_msg := "error response:" + err.Error()
				resp_result[addr] = err_msg
				return
			}
			resp_result[addr] = resp.Msg
		}(addr)
	}
	wg.Wait()
	return resp_result
}

func createDefaultConf(ctx *cli.Context) {
	dir := ctx.String("dir")
	file := ctx.String("file")
	task_conf := &config.CrawlerTask{
		TaskName:                    "TaskName-Spider",
		TaskId:                      0,
		CrawlerTreadNum:             1,
		RedisAddr:                   "RedisIp:Port",
		RedisPassword:               "12345",
		MaxGetRequestNullTimeSecond: 1 * time.Minute,
		MasterListenPort:            "31001",
		FastDfsAddr:                 "fastDfsIp:Port",
	}
	conf := &config.SpiderStartConfig{
		TaskConf:   *task_conf,
		ProcerName: "default",
		Urls:       []string{"http://www.baidu.com", "http://www.xbiquge.la/paihangbang/"},
	}
	file_path := dir + "/" + file
	err := config.WriteConfigToFileJson(dir, file_path, conf)
	if err != nil {
		fmt.Printf("create config fail:%s\n", err.Error())
		return
	}
	fmt.Printf("create config file succ:%s\n", file_path)
}

func exec_shell(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)

	//读取io.Writer类型的cmd.Stdout，再通过bytes.Buffer(缓冲byte类型的缓冲器)将byte类型转化为string类型(out.String():这是bytes类型提供的接口)
	var out bytes.Buffer
	cmd.Stdout = &out

	//Run执行c包含的命令，并阻塞直到完成。  这里stdout被取出，cmd.Wait()无法正确获取stdin,stdout,stderr，则阻塞在那了
	err := cmd.Run()
	return out.String(), err
}

func checkErr(err error, info string) {
	if err != nil {
		fmt.Printf("%s:%s\n", info, err.Error())
		// panic(me(err, info))
		os.Exit(1)
	}
}