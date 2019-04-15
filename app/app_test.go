/*
 * @Author: rayou
 * @Date: 2019-04-15 20:45:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-15 21:35:04
 */
package app

import (
	"testing"
	"time"

	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/arxmaster"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/spider"
	log "github.com/sirupsen/logrus"
)

var MasterSvr *arxmaster.ArxMaster
var SpiderClient *spider.Spider

func Init() {
	log.SetLevel(log.TraceLevel)
	MasterSvr = arxmaster.NewArxMaster()
	MasterSvr.Init("8888")
	SpiderClient = spider.NewSpider()
	SpiderClient.Init("9888")
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
	time.Sleep(1 * time.Second)
	err := arxlet.SendTcpMsg(message.MSG_ARXCMD_START_SPIDER, []byte(""), ":8888")
	if err != nil {
		t.Errorf("send msg to master fail:%s", err.Error())
	} else {
		log.Info("send msg succ")
	}
	time.Sleep(1 * time.Second)
	t.Log("done")

}
