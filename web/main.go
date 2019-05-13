/*
 * @Author: rayou
 * @Date: 2019-04-18 10:57:50
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-08 19:15:58
 */

package main

import (
	"fmt"

	"github.com/AzuresYang/arx7/config"
	"github.com/AzuresYang/arx7/web/webapp"
)

func main() {
	conf := &config.MasterConfig{}
	err := config.ReadConfigFromFileJson("F:\\master.json", conf)
	if err != nil {
		fmt.Printf("read config fail:%s\n", err.Error())
		return
	}
	webapp.Start(conf)
}
