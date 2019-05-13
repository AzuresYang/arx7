/*
 * @Author: rayou
 * @Date: 2019-04-21 11:31:53
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 23:20:30
 */
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/arxdeployment"
	"github.com/AzuresYang/arx7/config"
	log "github.com/sirupsen/logrus"
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
				cli.Uint64Flag{
					Name:  "port",
					Value: 31000,
					Usage: "the open listen port",
				},
			},
			Action: CmdDeploymentSpider,
		},
		// 扩缩容
		{
			Name:      "scale",
			ShortName: "s",
			Usage:     "scale spider nodes",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "spidername",
					Value: "default-spider",
					Usage: " spider task name",
				},
				cli.Uint64Flag{
					Name:  "nodes",
					Value: 1,
					Usage: "target num of the nodes",
				},
			},
			Action: scaleSpider,
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
			Action: deleteSpider,
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
			Action: echoSpider,
		},
	}
	app.Run(os.Args)
}

// 启动spider
func startSpider(ctx *cli.Context) {
	file := ctx.String("config")
	conf := &config.SpiderStartConfig{}
	err := config.ReadConfigFromFileJson(file, conf)
	if err != nil {
		fmt.Printf("read config file fail:[%s]%s\n", file, err.Error())
		return
	}
	spider_name := ctx.String("spidername")
	arxdeployment.DoStartSpider(spider_name, conf)
}

// 命令行测试
func getPod(ctx *cli.Context) {
	cmd := "kubectl get po -o wide"
	ret, err := exec_shell(cmd)
	checkErr(err, "get pods")
	fmt.Println(ret)
}

// 发布程序
func CmdDeploymentSpider(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	image := ctx.String("image")
	// node := ctx.Uint64("nodes")
	port := ctx.Uint64("port")
	fmt.Printf("deployment spider start....\n")
	fmt.Printf("spider[%s], image[%s],port[%d]\n", spider_name, image, port)
	if len(image) <= 0 {
		fmt.Println("image is null.")
		return
	}
	arxdeployment.DoDeploymentSpider(spider_name, image, port)
	// 这里应该加一个检查是否发布成功的东西
	// fmt.Printf("%s, %s,%d", spider_name, image, node)
}

// 扩缩容Spider
func scaleSpider(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	node := ctx.Uint64("nodes")
	fmt.Printf("scale spider [%s] nodes  to %d\n", spider_name, node)
	arxdeployment.DoScaleSpider(spider_name, fmt.Sprintf("%d", node))
}

// 删除Spider
func deleteSpider(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	fmt.Printf("delete spider [%s]......\n", spider_name)
	// 停止所有pod
	arxdeployment.DoDeleteSpider(spider_name)
}

//  停止 stop
func stopSpider(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	fmt.Printf("start stop spider:%s\n", spider_name)
	nodes := arxdeployment.GetSpiderNodes(spider_name)
	if nodes == nil {
		fmt.Println("get nodes addr fail.")
		return
	}
	ret := arxdeployment.SendMessageToSpider(nodes, message.MSG_REQ_STOP_SPIDER, []byte(""), "stop spider")
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}
}

// 获取spider状态
func getSpiderStatus(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	nodes := arxdeployment.GetSpiderNodes(spider_name)
	ret := arxdeployment.SendMessageToSpider(nodes, message.MSG_REQ_GET_SPIDER_INFO, []byte(""), "get spider status")
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}
}

// echo spider
func echoSpider(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	nodes := arxdeployment.GetSpiderNodes(spider_name)
	cmd := message.MSG_REG_ECHO
	if ctx.Uint64("model") == 2 {
		cmd = message.MSG_REG_ECHO_REDIS
	}
	ret := arxdeployment.SendMessageToSpider(nodes, cmd, []byte(""), "echo spider")
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}
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
		MaxGetRequestNullTimeSecond: 1,
		MasterListenPort:            "31001",
		FastDfsAddr:                 "fastDfsIp:Port",
	}
	conf := &config.SpiderStartConfig{
		TaskConf:   *task_conf,
		ProcerName: "default",
		Urls:       []string{"http://www.baidu.com", "http://www.xbiquge.la/paihangbang/"},
	}
	file_path := dir + "/" + file
	err := config.WriteConfigToFileJson(dir, file, conf)
	if err != nil {
		fmt.Printf("create config fail:%s\n", err.Error())
		return
	}
	fmt.Printf("create config file succ:%s\n", file_path)
}

func exec_shell(s string) (string, error) {
	//函数返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", s)
	log.Infof("bash cmd:%s", s)
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
