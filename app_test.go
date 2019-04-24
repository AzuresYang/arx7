/*
 * @Author: rayou
 * @Date: 2019-04-15 20:45:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 22:28:38
 */
package arx7

import (
	"fmt"
	"testing"
	"time"

	"github.com/AzuresYang/arx7/app/arxmaster"
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
	start_info := buildSpiderCfg()
	fmt.Printf("get config:%+v\n", start_info)
	if start_info == nil {
		return
	}
	nodes := []string{"127.0.0.1:9888", "127.0.0.1:9889"}
	// send_bytes, _ := json.Marshal(start_info)
	// ret := arxdeployment.SendMessageToSpider(start_info.NodeAddrs, message.MSG_REQ_GET_SPIDER_INFO, []byte(""), "echo")
	arxdeployment.InitRedis(start_info)
	arxdeployment.SendStartToMaster(start_info, nodes)
	time.Sleep(20 * time.Second)
	t.Log("done")

}
