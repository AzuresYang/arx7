/*
 * @Author: rayou
 * @Date: 2019-04-27 16:01:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-27 22:16:50
 */

package controller

import (
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func MonitorInfoHandler(response http.ResponseWriter, request *http.Request) {
	log.Info("Get R")
	fmt.Fprintf(response, "Hello, Welcome to go web programming...")
}

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
