/*
 * @Author: rayou
 * @Date: 2019-04-15 20:45:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-21 21:10:28
 */
package app

import (
	"fmt"
	"testing"
	"time"

	"github.com/AzuresYang/arx7/app/arxdeployment"
	"github.com/AzuresYang/arx7/app/arxmaster"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/spider"
	log "github.com/sirupsen/logrus"
)

var MasterSvr *arxmaster.ArxMaster
var SpiderClient *spider.Spider
var SpiderClient2 *spider.Spider

func Init() {
	log.SetLevel(log.TraceLevel)
	MasterSvr = arxmaster.NewArxMaster()
	MasterSvr.Init("8888")
	SpiderClient = spider.NewSpider()
	SpiderClient.Init("9888")
	SpiderClient2 = spider.NewSpider()
	SpiderClient2.Init("9889")
}
func TestStartSpider(t *testing.T) {
	Init()
	// log.SetLevel(log.TraceLevel)
	// MasterSvr = arxmaster.NewArxMaster()
	// MasterSvr.Init("8888")
	// SpiderClient = spider.NewSpider()
	// SpiderClient.Init("9888")
	go MasterSvr.Run()
	go SpiderClient.Run()
	go SpiderClient2.Run()
	go time.Sleep(1 * time.Second)
	start_info := message.SpiderStartMsg{}
	start_info.NodeAddrs = []string{
		":9888",
		":9889",
	}
	// send_bytes, _ := json.Marshal(start_info)
	ret := arxdeployment.SendMessageToSpider(start_info.NodeAddrs, message.MSG_REQ_GET_SPIDER_INFO, []byte(""), "echo")
	for node, msg := range ret {
		fmt.Printf("node:%s,    start result:%s\n", node, msg)
	}

	time.Sleep(1 * time.Second)
	t.Log("done")

}
