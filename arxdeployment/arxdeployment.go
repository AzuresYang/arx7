/*
 * @Author: rayou
 * @Date: 2019-04-21 11:31:53
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-10 01:28:48
 */
package arxdeployment

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
	"github.com/AzuresYang/arx7/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// 从kuberntet中获取service关联到的节点
func GetSpiderNodes(svc_name string) []string {
	// 获取port的方式

	// return []string{"132.232.43.251:31001", "132.232.43.251:31002"}
	cmd := fmt.Sprintf("kubectl get service | grep -w \"%s\" | awk -F \" \" '{print $4}'", svc_name)
	ret_line, err := exec_shell(cmd)
	fmt.Printf("service ret:%s", ret_line)
	if err != nil {
		fmt.Printf("get node info %s fail:%s\n", svc_name, err.Error())
		return nil
	}
	if len(ret_line) <= 0 {
		fmt.Printf("not found service:%s\n", svc_name)
		return nil
	}
	temp := strings.Split(ret_line, "/")
	// 按照规则是80:port/TCP的样式, 两次切割都不需要符合
	if len(temp) != 2 {
		fmt.Printf("[%s] get node info fail:%s\n", svc_name, ret_line)
		return nil
	}
	if len(strings.Split(temp[0], ":")) != 2 {
		fmt.Printf("[%s] get node info fail:%s\n", svc_name, ret_line)
		return nil
	}
	port := strings.Split(temp[0], ":")[1]

	// 获取ip地址
	cmd = fmt.Sprintf("kubectl get po -o wide | grep  \"%s-\" | awk -F \" \" '{print $7}'", svc_name)
	ret_line, err = exec_shell(cmd)
	fmt.Printf("service ret:%s\n", ret_line)
	if err != nil {
		fmt.Printf("[%s]get pod info fail:%s\n", svc_name, err.Error())
		return nil
	}
	if len(ret_line) <= 0 {
		fmt.Printf("[%s] not found deploy instance .ip is zero.\n", svc_name)
		return nil
	}
	ips := strings.Split(ret_line, "\n")
	nodes := []string{}
	for _, ip := range ips {
		if len(ip) <= 0 {
			continue
		}
		if len(strings.Split(ip, ".")) != 4 {
			fmt.Printf("[%s]get pod ip fail:%s\n", svc_name, ret_line)
			return nil
		}
		addr := fmt.Sprintf("%s:%s", ip, port)
		addr = fmt.Sprintf("%s:31001", ip)
		nodes = append(nodes, addr)
	}
	log.Infof("[%s]get nodes:%+v\n", svc_name, nodes)
	return nodes
}

func DoStartSpider(spider_name string, conf *config.SpiderStartConfig) (error, string) {
	err := InitRedis(conf)
	if err != nil {
		fmt.Printf("redis[%s] init fail:%s\n", conf.TaskConf.RedisAddr, err.Error())
		return err, "Init Redis fail"
	}
	ret_msg := "init redis succ" + "\n"
	nodes := GetSpiderNodes(spider_name)
	if len(nodes) <= 0 {
		ret_msg += "No found spider nodes.please ensure had deployment spider." + "\n"
		fmt.Printf("not found spider nodes.please ensure had deployed spider task:%s\n", conf.TaskConf.TaskName)
		return errors.New("No found spider nodes.please ensure had deployment spider."), ret_msg
	}
	// fmt.Printf("config:%+v\n", conf)
	nerr, resp := SendStartToMaster(conf, nodes)
	if nerr != nil {
		ret_msg += nerr.Error() + "\n"
		return nerr, ret_msg
	}
	var start_ret map[string]string
	json.Unmarshal(resp.Data, &start_ret)
	for node, msg := range start_ret {
		ret_msg = fmt.Sprintf("%s\n[%s]:%s", ret_msg, node, msg)
	}
	fmt.Println("start spider done.")
	return nil, ret_msg
}

// 开放权限，调试用
func InitRedis(conf *config.SpiderStartConfig) error {
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
	for _, s := range conf.Urls {
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

func SendStartToMaster(conf *config.SpiderStartConfig, nodes []string) (error, *message.ResponseMsg) {
	fmt.Println("sending start info to arxmaster......")
	start_info := message.SpiderStartMsg{}
	start_info.Cfg = conf.TaskConf
	start_info.NodeAddrs = nodes
	send_bytes, nerr := json.Marshal(start_info)
	if nerr != nil {
		fmt.Printf("marshal start info to master fail:%s\n", nerr.Error())
		return errors.New("marshal start info to master fail"), nil
	}
	// fmt.Printf("send to master: config :%s\n", string(send_bytes))
	masterAddr := "127.0.0.1:" + conf.TaskConf.MasterListenPort
	err, cnn := arxlet.SendTcpMsgTimeoutWithConn(message.MSG_ARXCMD_START_SPIDER, send_bytes, masterAddr, 10*time.Second)
	if err != nil {
		return errors.New(fmt.Sprintf("connect master fail:%s", err.Error())), nil
	}
	rerr, resp_msg := arxlet.ParseResponseFromConn(cnn)
	if rerr != nil {
		return errors.New(fmt.Sprintf("get master response fail:%s", rerr.Error())), nil
	}
	fmt.Println("start info:")
	fmt.Printf("succ:%d, fail:%d\n", uint32(len(nodes))-resp_msg.Status, resp_msg.Status)
	var ret map[string]string
	err = json.Unmarshal(resp_msg.Data, &ret)
	if err != nil {
		fmt.Println("Unmarshal master resp data fail")
		return errors.New("Unmarshal master resp data fail"), nil
	}
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}
	fmt.Println("start spider down")
	return nil, resp_msg
}

func DoDeploymentSpider(spider_name string, image string, port uint64) (error, string) {
	cmd := fmt.Sprintf("kubectl run %s --image=%s --port %d", spider_name, image, port)
	ret_msg := ""
	ret, err := exec_shell(cmd)
	if err != nil {
		fmt.Printf("[error]run image fail.%s\n", err.Error())
		return err, ""
	}
	ret_msg += ret + "\n"
	// 服务暴露
	cmd = fmt.Sprintf("kubectl expose deploy %s --port %d --target-port %d --type NodePort", spider_name, port, port)
	ret, err = exec_shell(cmd)
	ret_msg += ret
	if err != nil {
		return err, ret_msg
	}
	fmt.Println(ret)
	cmd = fmt.Sprintf("kubectl get po -o wide | grep -w %s", spider_name)
	ret, err = exec_shell(cmd)
	fmt.Println("deployment ret:")
	fmt.Println(ret)
	return nil, ret_msg
}

func DoScaleSpider(spider_name string, node string) (error, string) {
	cmd := fmt.Sprintf("kubectl scale deployment %s --replicas=%s", spider_name, node)
	ret, err := exec_shell(cmd)
	if err != nil {
		fmt.Printf("[error]%s\n", err.Error())
		return err, ""
	}
	fmt.Println(ret)
	return nil, ret
}

func DoDeleteSpider(spider_name string) (error, string) {

	cmd := fmt.Sprintf("kubectl scale deployment %s --replicas=0", spider_name)
	ret, err := exec_shell(cmd)
	if err != nil {
		fmt.Printf("[error]stop all pods fail.%s\n", err.Error())
		return err, ""
	}
	ret_msg := ret
	fmt.Printf("stop all pod ret:%s\n", ret)
	cmd = fmt.Sprintf("kubectl delete deployment %s", spider_name)
	ret, err = exec_shell(cmd)
	if err != nil {
		ret_msg += "\n" + "[error]" + err.Error() + "\n" + ret
	} else {
		ret_msg += "\n" + ret
	}
	fmt.Println(ret)
	cmd = fmt.Sprintf("kubectl delete service %s", spider_name)
	ret, err = exec_shell(cmd)
	fmt.Println(ret)
	if err != nil {
		ret_msg += "\n" + "[error]" + err.Error() + "\n" + ret
	} else {
		ret_msg += "\n" + ret
	}
	fmt.Println("delete down.")
	return nil, ret_msg
}

// 获取spider状态
func getSpiderStatus(ctx *cli.Context) {
	spider_name := ctx.String("spidername")
	nodes := GetSpiderNodes(spider_name)
	ret := SendMessageToSpider(nodes, message.MSG_REQ_GET_SPIDER_INFO, []byte(""), "get spider status")
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}
}

func DoStopSpider(nodes []string) map[string]string {
	return SendMessageToSpider(nodes, message.MSG_REQ_STOP_SPIDER, []byte(""), "stop spider")
}

func DoGetSpiderStatusByNodes(nodes []string) map[string]string {
	return SendMessageToSpider(nodes, message.MSG_REQ_GET_SPIDER_INFO, []byte(""), "get spider status")
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

func DoGetPod() string {
	cmd := "kubectl get po -o wide | grep -v READY"
	ret, _ := exec_shell(cmd)
	return ret
}

func DoGetSpiderPod(spider_name string) string {
	cmd := fmt.Sprintf("kubectl get po -o wide | grep \"%s\"", spider_name)
	ret, _ := exec_shell(cmd)
	return ret
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
