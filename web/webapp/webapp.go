/*
 * @Author: rayou
 * @Date: 2019-05-01 19:13:41
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-10 01:36:59
 */
/*
 * @Author: rayou
 * @Date: 2019-04-18 10:57:50
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-05 14:51:30
 */

package webapp

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/AzuresYang/arx7/config"
	"github.com/AzuresYang/arx7/web/controller"
)

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

func Start(conf *config.MasterConfig) {
	controller.DbService.Init(&conf.MysqlConf)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./static/fonts"))))
	http.HandleFunc("/index/", Index)
	http.HandleFunc("/get/monitor", controller.MonitorInfoHandler)
	http.HandleFunc("/get/pods", controller.ClusterHandlerGetPods)
	http.HandleFunc("/cluster/spiderstatus", controller.ClusterHandlerGetSpiderStatus)
	http.HandleFunc("/cluster/deployment", controller.ClusterHandlerDeployment)
	http.HandleFunc("/cluster/start", controller.ClusterHandlerStartSpider)
	http.HandleFunc("/cluster/stop", controller.ClusterHandlerStopSpider)
	http.HandleFunc("/cluster/scale", controller.ClusterHandlerScalePods)
	http.HandleFunc("/cluster/delete", controller.ClusterHandlerDeleteSpider)

	fmt.Printf("WebApp listenAddr:%s\n", conf.WebListenAddr)
	http.ListenAndServe(conf.WebListenAddr, nil)
}
