/*
 * @Author: rayou
 * @Date: 2019-04-18 10:57:50
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-05 14:51:30
 */

package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/AzuresYang/arx7/config"
	"github.com/AzuresYang/arx7/web/controller"
)

// "	"github.com/AzuresYang/arx7/app/arxlet""
func Hello(response http.ResponseWriter, request *http.Request) {
	type person struct {
		Id      int
		Name    string
		Country string
	}
	fmt.Println("get conn")
	liumiaocn := person{Id: 1001, Name: "liumiaocn", Country: "China"}

	tmpl, err := template.ParseFiles("./template/user.tpl")
	if err != nil {
		fmt.Println("Error happened..")
	}
	tmpl.Execute(response, liumiaocn)
}

func Index(response http.ResponseWriter, request *http.Request) {
	type person struct {
		Id      int
		Name    string
		Country string
	}
	fmt.Println("get indx conn")
	liumiaocn := person{Id: 1001, Name: "liumiaocn", Country: "China"}

	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		fmt.Println("Error happened..")
	}
	tmpl.Execute(response, liumiaocn)
}

func main() {
	conf := &config.MasterConfig{}
	err := config.ReadConfigFromFileJson("F:\\master.json", conf)
	if err != nil {
		fmt.Printf("read config fail:%s\n", err.Error())
		return
	}
	controller.DbService.Init(&conf.MysqlConf)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	// http.Handle("/css/", http.FileServer(http.Dir("template")))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./static/fonts"))))
	// http.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
	http.HandleFunc("/index/", Index)
	http.HandleFunc("/get/monitor", controller.MonitorInfoHandler)
	http.HandleFunc("/get/pods", controller.ClusterHandlerGetPods)
	http.HandleFunc("/cluster/spiderstatus", controller.ClusterHandlerGetSpiderStatus)
	http.HandleFunc("/cluster/deployment", controller.ClusterHandlerDeployment)
	http.HandleFunc("/cluster/start", controller.ClusterHandlerStartSpider)
	http.HandleFunc("/cluster/scale", controller.ClusterHandlerScalePods)
	http.HandleFunc("/cluster/delete", controller.ClusterHandlerDeleteSpider)
	http.HandleFunc("/", Hello)
	http.ListenAndServe(":8888", nil)
}
