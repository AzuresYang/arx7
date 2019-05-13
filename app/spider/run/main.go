/*
 * @Author: rayou
 * @Date: 2019-05-08 12:59:04
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-08 13:17:57
 */

package main

import (
	"fmt"
	"time"

	"github.com/AzuresYang/arx7/app/processor/biqu"
	"github.com/AzuresYang/arx7/app/spider"

	// "github.com/AzuresYang/arx7/app/pipeline"
	"github.com/AzuresYang/arx7/app/arxmonitor/monitorHandler"
	"github.com/AzuresYang/arx7/arxdeployment"
	"github.com/AzuresYang/arx7/config"
	log "github.com/sirupsen/logrus"
	// "github.com/AzuresYang/arx7/util/record"
)

func main() {
	log.SetLevel(log.TraceLevel)
	procer := biqu.NewProcessor()
	procer.GetName()
	// processor.Manager.Register(&procer)
	sp := spider.NewSpider()
	cfg := buildCfg()
	crawler_config := &config.CrawlerConfig{
		TaskConf:   cfg.TaskConf,
		MasterAddr: "127.0.0.1",
	}
	arxdeployment.InitRedis(cfg)
	sp.StartCrawler(crawler_config)
	monitorHandler.Stop()
	ticker := time.NewTicker(5 * time.Second)
	start := time.Now()

	for !sp.IfCrawlerEngineStop() {
		select {
		case <-ticker.C:
			cost := time.Since(start)
			fmt.Printf("spider running....time:%s\n", cost)
		}
	}

}

func buildCfg() *config.SpiderStartConfig {
	cfg := &config.SpiderStartConfig{}
	cfg.TaskConf = config.CrawlerTask{
		TaskName:                    "BiQuSpider",
		TaskId:                      0,
		CrawlerTreadNum:             5,
		RedisAddr:                   "193.112.68.221:6379",
		RedisPassword:               "Redis@2019416",
		FastDfsAddr:                 "http://172.17.87.202:8080",
		MaxGetRequestNullTimeSecond: 10, // 没有链接的停止时间
	}
	cfg.ProcerName = "biqu"
	cfg.Urls = []string{
		"http://www.xbiquge.la/paihangbang/",
	}
	return cfg
}
