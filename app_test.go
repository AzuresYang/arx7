/*
 * @Author: rayou
 * @Date: 2019-04-15 20:45:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 23:44:36
 */
package arx7

import (
	"fmt"
	"testing"
	"time"

	"github.com/AzuresYang/arx7/app/arxmaster"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/spider"
	"github.com/AzuresYang/arx7/arxdeployment"
	"github.com/AzuresYang/arx7/config"
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

func buildSpiderCfg() *config.SpiderStartConfig {
	conf := &config.SpiderStartConfig{}
	err := config.ReadConfigFromFileJson("E:\\GoPath\\spider1.json", conf)
	if err != nil {
		fmt.Printf("read config fail:%s\n", err.Error())
		return nil
	}
	//     fmt.Printf("get configd:%+v\n", conf)
	return conf
}
func buildMasterCfg() *config.MasterConfig {
	conf := &config.MasterConfig{}
	err := config.ReadConfigFromFileJson("F:\\master.json", conf)
	if err != nil {
		fmt.Printf("read config fail:%s\n", err.Error())
		return nil
	}
	fmt.Printf("get configd:%+v\n", conf)
	return conf
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
	// go SpiderClient2.Run()
	go time.Sleep(1 * time.Second)
	masterCfg := buildMasterCfg()
	if masterCfg == nil {
		return
	}
	// MasterSvr.StartMonitorCollector(&masterCfg.MysqlConf)

	start_info := buildSpiderCfg()
	// fmt.Printf("get config:%+v\n", start_info)
	if start_info == nil {
		return
	}
	nodes := []string{"127.0.0.1:9888", "127.0.0.1:9889"}
	arxdeployment.InitRedis(start_info)
	arxdeployment.SendStartToMaster(start_info, nodes)
	time.Sleep(20 * time.Second)
	t.Log("done")

}

func TestArxlet(t *testing.T) {
	Init()
	go MasterSvr.Run()
	go SpiderClient.Run()
	go time.Sleep(1 * time.Second)
	masterCfg := buildMasterCfg()
	if masterCfg == nil {
		return
	}
	// MasterSvr.StartMonitorCollector(&masterCfg.MysqlConf)

	start_info := buildSpiderCfg()
	// fmt.Printf("get config:%+v\n", start_info)
	if start_info == nil {
		return
	}
	nodes := []string{"127.0.0.1:9888"}
	// arxdeployment.InitRedis(start_info)
	arxdeployment.SendMessageToSpider(nodes, message.MSG_REG_ECHO, []byte(""), "echo")
	arxdeployment.SendMessageToSpider(nodes, message.MSG_REQ_STOP_SPIDER, []byte(""), "stop spider")
	arxdeployment.SendStartToMaster(start_info, nodes)
	time.Sleep(20 * time.Second)
	t.Log("done")

}
