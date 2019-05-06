/*
 * @Author: rayou
 * @Date: 2019-04-27 16:01:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-06 23:26:12
 */
package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/AzuresYang/arx7/arxdeployment"
	"github.com/AzuresYang/arx7/config"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

const (
	default_spider_port uint64 = 31001
)

func ClusterHandlerGetPods(response http.ResponseWriter, request *http.Request) {
	log.Infof("Get pods Conn. Form:%#v", request.Form)

	pod1 := ClusterInfo{
		SpiderName: "spider01",
		NodeStatus: "1/1",
		RunStatus:  "Running",
		Age:        "6d",
		NodeAddr:   "127.0.0.1",
		Desc:       "desc",
	}
	pod2 := pod1
	pod2.SpiderName = "spider002"
	cluster_infos := []ClusterInfo{pod1, pod2}
	jdata, _ := json.Marshal(cluster_infos)
	log.Infof("response pods data")
	fmt.Fprintf(response, string(jdata))
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
	err := responseJson(response, 0, ret_msg, ret_msg)
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
	log.Infof("[%s]start ret:config:%+v\n", code_info, conf)
	err = responseJson(response, 0, "succ", "")
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
