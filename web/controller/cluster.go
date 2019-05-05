/*
 * @Author: rayou
 * @Date: 2019-04-27 16:01:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-04 18:46:52
 */
package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
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

func ClusterHandlerScalePods(response http.ResponseWriter, request *http.Request) {
	log.Infof("Conn scale pods. Form:%#v", request.Form)

	pod1 := ClusterInfo{
		SpiderName: "scale",
		NodeStatus: "1/1",
		RunStatus:  "Running",
		Age:        "6d",
		NodeAddr:   "123456",
		Desc:       "is running",
	}
	cluster_infos := []ClusterInfo{pod1}
	err := responseJson(response, 0, "succ", cluster_infos)
	if err != nil {
		log.Errorf("[monitorInfoHandler] response err:%s", err.Error())
	}
}
