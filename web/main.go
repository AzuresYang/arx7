/*
 * @Author: rayou 
 * @Date: 2019-04-18 10:57:50 
 * @Last Modified by:   rayou 
 * @Last Modified time: 2019-04-18 10:57:50 
 */

package main
import (
    "net/http"
)
 
func main() {
    http.Handle("/css/", http.FileServer(http.Dir("template")))
    http.Handle("/js/", http.FileServer(http.Dir("template")))
     
    http.HandleFunc("/index/", adminHandler)
    http.HandleFunc("/login/",loginHandler)
    http.HandleFunc("/ajax/",ajaxHandler)
    http.HandleFunc("/",NotFoundHandler)
    http.ListenAndServe(":8888", nil)
 
}