/*
 * @Author: rayou
 * @Date: 2019-04-29 19:15:59
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-10 01:15:28
 */

package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type ResponseData struct {
	Status uint32      `json:"Status"`
	Msg    string      `json:"Msg"`
	Data   interface{} `json:"Data"`
}

type FormQueryMonitor struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	SvcId     string `json:"svcid"`
	Metric    string `json:"metric"`
	Classfy   string `json:"classfy"`
}

type MonitorData struct {
	Id      int
	Svcid   int
	Metric  int
	Classfy int
	Value   int
	Ip      string
	Time    int64 // 给前端的是这个数据，就这样吧
}

type SpiderNodeInfo struct {
	SpiderName string `json:"SpiderName"`
	NodeStatus string `json:"NodeStatus"`
	RunStatus  string `json:"RunStatus"`
	Age        string `json:"Age"`
	NodeAddr   string `json:"NodeAddr"`
	Desc       string `json:"Desc"`
}

type MonitorDataWrapper struct {
	Data []MonitorData
}

func (self *MonitorDataWrapper) Len() int { // 重写 Len() 方法
	return len(self.Data)
}
func (self *MonitorDataWrapper) Swap(i, j int) { // 重写 Swap() 方法
	self.Data[i], self.Data[j] = self.Data[j], self.Data[i]
}
func (self *MonitorDataWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return self.Data[i].Time < self.Data[j].Time
}

func parseForm(request *http.Request, st interface{}) error {
	request.ParseForm()
	log.Infof("Form:%#v", request.Form)
	var new_form = make(map[string]string)
	for k, v := range request.Form {
		if v[0] == "" {
			return errors.New("表单参数不能为空")
		}
		new_form[k] = v[0]
	}
	log.Infof("Form parse:%#v", new_form)
	jdata, err := json.Marshal(new_form)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jdata, st)
	if err != nil {
		return err
	}
	return nil
}

func parseFormAsMap(request *http.Request) map[string]string {
	request.ParseForm()
	var new_form = make(map[string]string)
	for k, v := range request.Form {
		new_form[k] = v[0]
	}
	return new_form
}

func responseJson(response http.ResponseWriter, status uint32, msg string, data interface{}) error {
	resp := ResponseData{
		Status: status,
		Msg:    msg,
		Data:   data,
	}
	jdata, err := json.Marshal(resp)
	if err != nil {
		return err
	}
	// 允许跨域
	//跨域请求，*代表允许全部类型
	response.Header().Set("Access-Control-Allow-Origin", "*")
	//允许请求方式
	response.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE")
	//用来指定本次预检请求的有效期，单位为秒，在此期间不用发出另一条预检请求
	response.Header().Set("Access-Control-Max-Age", "3600")
	//请求包含的字段内容，如有多个可用哪个逗号分隔如下
	response.Header().Set("Access-Control-Allow-Headers", "content-type,x-requested-with,Authorization, x-ui-request,lang")
	fmt.Fprintf(response, string(jdata))
	return nil
}
