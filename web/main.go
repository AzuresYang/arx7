/*
 * @Author: rayou
 * @Date: 2019-04-18 10:57:50
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-30 12:33:20
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
	fmt.Println("get conn")
	liumiaocn := person{Id: 1001, Name: "liumiaocn", Country: "China"}

	tmpl, err := template.ParseFiles("./template/index.html")
	if err != nil {
		fmt.Println("Error happened..")
	}
	tmpl.Execute(response, liumiaocn)
}

func Vue(response http.ResponseWriter, request *http.Request) {
	fmt.Println("get vue con")
	http.ServeFile(response, request, "./template/index.html")
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
	http.HandleFunc("/vue", Vue)
	http.HandleFunc("/get/monitor", controller.MonitorInfoHandler)
	// http.HandleFunc("/login/",loginHandler)
	// http.HandleFunc("/ajax/",ajaxHandler)
	http.HandleFunc("/", Hello)
	http.ListenAndServe(":8888", nil)
	// r := gin.Default()
	// r.LoadHTMLGlob("template/*.html")              // 添加入口index.html
	// r.LoadHTMLFiles("static/*/*")              // 添加资源路径
	// r.Static("/static", "./dist/static")       // 添加资源路径
	// r.StaticFile("/index/", "dist/index.html") //前端接口
	// r.Run(":8888")
}
