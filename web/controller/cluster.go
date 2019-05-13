/*
 * @Author: rayou
 * @Date: 2019-04-27 16:01:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-10 01:36:13
 */
package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/AzuresYang/arx7/arxdeployment"
	"github.com/AzuresYang/arx7/config"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

const (
	default_spider_port uint64 = 31001
)

func GetSpiderNodeInfo(spider_name string) []SpiderNodeInfo {
	node_infos := []SpiderNodeInfo{}
	kube_msg := arxdeployment.DoGetSpiderPod(spider_name)
	if !strings.Contains(kube_msg, spider_name) {
		return node_infos
	}
	lines := strings.Split(kube_msg, "\n")
	for _, line := range lines {
		item := strings.Fields(line)
		if len(item) <= 0 {
			continue
		}
		node_infos = append(node_infos, SpiderNodeInfo{
			SpiderName: spider_name,
			NodeStatus: item[1],
			RunStatus:  item[2],
			Age:        item[4],
			NodeAddr:   item[6],
			Desc:       ""})
	}
	return node_infos
}

func GetClusterNodeInfo() []SpiderNodeInfo {
	node_infos := []SpiderNodeInfo{}
	kube_msg := arxdeployment.DoGetPod()
	if len(kube_msg) <= 0 {
		return node_infos
	}
	lines := strings.Split(kube_msg, "\n")
	for _, line := range lines {
		item := strings.Fields(line)
		if len(item) != 7 {
			continue
		}
		node_infos = append(node_infos, SpiderNodeInfo{
			SpiderName: item[0],
			NodeStatus: item[1],
			RunStatus:  item[2],
			Age:        item[4],
			NodeAddr:   item[6],
			Desc:       ""})
	}
	return node_infos
}

func ClusterHandlerGetPods(response http.ResponseWriter, request *http.Request) {
	log.Infof("Get pods Conn. Form:%#v", request.Form)
	spider_infos := GetClusterNodeInfo()
	log.Infof("[Get Pods]%+v", spider_infos)
	responseJson(response, 0, "succ", spider_infos)
}

// 获取爬虫状态
func ClusterHandlerGetSpiderStatus(response http.ResponseWriter, request *http.Request) {
	code_info := "Cluster.GetSpiderStatus"
	request.ParseForm()
	log.Infof("Conn Spider status. Form:%#v", request.Form)
	form := parseFormAsMap(request)
	spider_name := form["spidername"]
	if spider_name == "" {
		log.Errorf("[%s]参数不能为空", code_info)
		responseJson(response, 1, "参数不能为空", "")
		return
	}
	nodes := arxdeployment.GetSpiderNodes(spider_name)
	if len(nodes) <= 0 {
		log.Errorf("[%s]没有该爬虫的部署信息:%s", code_info, spider_name)
		responseJson(response, 0, "没有该爬虫的部署信息："+spider_name, "")
		return
	}
	ret_map := arxdeployment.DoGetSpiderStatusByNodes(nodes)
	ret_msg := "Get Status Ret:"
	for node, msg := range ret_map {
		ret_msg = fmt.Sprintf("%s\n---[%s]:%s---", ret_msg, node, msg)
	}
	log.Infof("[%s]get status ret:%s", code_info, ret_msg)
	spider_infos := GetSpiderNodeInfo(spider_name)
	err := responseJson(response, 0, ret_msg, spider_infos)
	if err != nil {
		log.Errorf("[monitorInfoHandler] response err:%s", err.Error())
	}
}

// 部署爬虫
func ClusterHandlerDeployment(response http.ResponseWriter, request *http.Request) {
	code_info := "Cluster.Deployment"
	request.ParseForm()
	log.Infof("Conn deployment spider. Form:%#v", request.Form)
	form := parseFormAsMap(request)
	spider_name := form["spidername"]
	image := form["image"]
	if spider_name == "" || image == "" {
		log.Errorf("[%s]参数不能为空", code_info)
		responseJson(response, 1, "参数不能为空", "")
		return
	}
	var ret_code uint32 = 0
	err, ret := arxdeployment.DoDeploymentSpider(spider_name, image, default_spider_port)
	if err != nil {
		ret_code = 1
		ret = fmt.Sprintf("[error]%s:---%s", err.Error(), ret)
	}
	log.Infof("[%s]ret:%s", code_info, ret)
	err = responseJson(response, ret_code, ret, "")
	if err != nil {
		log.Errorf("[%s] response err:%s", code_info, err.Error())
	}
}

// 启动爬虫
func ClusterHandlerStartSpider(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	code_info := "CusterHandler.StartSpider"
	log.Infof("Conn Start Spider. Form:%#v", request.Form)

	// 检查
	file, file_header, err := request.FormFile("config")
	if err != nil {
		log.Errorf("[%s]配置文件不能为空", code_info)
		responseJson(response, 2, "配置文件不能为空", "")
		return
	}
	defer file.Close()

	spider_name := request.PostFormValue("spidername")
	if spider_name == "" {
		log.Errorf("[%s]参数不能为空", code_info)
		responseJson(response, 1, "参数不能为空", "")
		return
	}

	nodes := arxdeployment.GetSpiderNodes(spider_name)
	if len(nodes) <= 0 {
		log.Errorf("[%s]没有该爬虫的部署信息:%s", code_info, spider_name)
		responseJson(response, 1, "没有该爬虫的部署信息："+spider_name, "")
		return
	}
	// 读取文件
	file_bytes, ierr := ioutil.ReadAll(file)
	if ierr != nil {
		log.Errorf("[%s]读取上传文件错误:%s", code_info, ierr.Error())
		responseJson(response, 3, "读取配置文件有误，请确认是json文件", "")
		return
	}
	log.Infof("[%s]spider:%s, config:%s | %s", code_info, spider_name, file_header.Filename, string(file_bytes))

	// 检查配置文件合法性
	conf := &config.SpiderStartConfig{}
	err = json.Unmarshal(file_bytes, conf)
	if err != nil {
		log.Errorf("[%s]配置文件不是Spider配置:%s", code_info, err.Error())
		responseJson(response, 4, "配置文件不是Spider配置,请确认是Spider的json配置文件", "")
		return
	}

	var ret_code uint32 = 0
	nerr, ret := arxdeployment.DoStartSpider(spider_name, conf)
	if nerr != nil {
		ret_code = 5
		ret = fmt.Sprintf("[error]%s:---%s", err.Error(), ret)
	}
	log.Infof("[%s]start ret:config:%+v\n", code_info, conf)
	err = responseJson(response, ret_code, ret, "")
	if err != nil {
		log.Errorf("[monitorInfoHandler] response err:%s", err.Error())
	}
}

// 节点扩缩容
func ClusterHandlerScalePods(response http.ResponseWriter, request *http.Request) {
	code_info := "Cluster.Scale"
	request.ParseForm()
	log.Infof("Conn scale pods. Form:%#v", request.Form)
	form := parseFormAsMap(request)
	spider_name := form["spidername"]
	nodes := form["nodes"]
	if spider_name == "" || nodes == "" {
		log.Errorf("[%s]参数不能为空", code_info)
		responseJson(response, 1, "参数不能为空", "")
		return
	}
	var ret_code uint32 = 0
	err, ret := arxdeployment.DoScaleSpider(spider_name, nodes)
	if err != nil {
		ret_code = 1
		ret = fmt.Sprintf("[error]%s:---%s", err.Error(), ret)
	}
	log.Infof("[%s]ret:%s", code_info, ret)
	err = responseJson(response, ret_code, ret, "")
	if err != nil {
		log.Errorf("[monitorInfoHandler] response err:%s", err.Error())
	}
}

// STOP爬虫
func ClusterHandlerStopSpider(response http.ResponseWriter, request *http.Request) {
	code_info := "Cluster.StopSpider"
	request.ParseForm()
	log.Infof("Conn stop spider status. Form:%#v", request.Form)
	form := parseFormAsMap(request)
	spider_name := form["spidername"]
	if spider_name == "" {
		log.Errorf("[%s]参数不能为空", code_info)
		responseJson(response, 1, "参数不能为空", "")
		return
	}
	nodes := arxdeployment.GetSpiderNodes(spider_name)
	if len(nodes) <= 0 {
		log.Errorf("[%s]没有该爬虫的部署信息:%s", code_info, spider_name)
		responseJson(response, 0, "没有该爬虫的部署信息："+spider_name, "")
		return
	}
	ret_map := arxdeployment.DoStopSpider(nodes)
	spider_infos := GetSpiderNodeInfo(spider_name)
	ret_msg := "Stop spider Status Ret:"
	for node, msg := range ret_map {
		ret_msg = fmt.Sprintf("%s\n---[%s]:%s---", ret_msg, node, msg)
		ip := strings.Split(node, ":")[0]
		for i, _ := range spider_infos {
			if spider_infos[i].NodeAddr == ip {
				spider_infos[i].Desc = msg
			}
		}
	}
	log.Infof("[%s]stop spider ret:%s", code_info, ret_msg)
	err := responseJson(response, 0, ret_msg, spider_infos)
	if err != nil {
		log.Errorf("[monitorInfoHandler] response err:%s", err.Error())
	}
}

// 删除爬虫
func ClusterHandlerDeleteSpider(response http.ResponseWriter, request *http.Request) {
	code_info := "Cluster.DeleteSpider"
	request.ParseForm()
	log.Infof("Conn delete. Form:%#v", request.Form)
	form := parseFormAsMap(request)
	spider_name := form["spidername"]
	if spider_name == "" {
		log.Errorf("[%s]参数不能为空", code_info)
		responseJson(response, 1, "参数不能为空", "")
		return
	}
	var ret_code uint32 = 0
	err, ret := arxdeployment.DoDeleteSpider(spider_name)
	if err != nil {
		ret_code = 1
		ret = fmt.Sprintf("[error]%s:---%s", err.Error(), ret)
	}
	log.Infof("[%s]ret:%s", code_info, ret)
	err = responseJson(response, ret_code, ret, "")
	if err != nil {
		log.Errorf("[monitorInfoHandler] response err:%s", err.Error())
	}
}
