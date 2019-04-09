/*
 * @Author: rayou
 * @Date: 2019-04-02 19:51:55
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-09 22:47:25
 */
package spider

import (
	"github.com/AzuresYang/arx7/app/pipeline/output"
	"github.com/AzuresYang/arx7/app/spider/crawler"
	"github.com/AzuresYang/arx7/app/spider/downloader"
	log "github.com/sirupsen/logrus"
)

type Spider struct {
	to_stop chan int
	to_end  chan int
}

var crawler_pool crawler.CrawlerPool = crawler.NewCrawlerPool(5)

func (spider *Spider) Start() {
	log.Info("spider start")
	// request.RequestMgr.Init(20)
	// 包括crawler 池的初始化 ？？？
	craw_num := 5
	for i := 0; i < craw_num; i++ {
		cw := crawler_pool.Get()
		pipe := &output.OutputFile{}
		dl := &downloader.SimpleDownloader{}
		cw.Init(dl, pipe)
		go cw.Run()
	}
}

func (spider *Spider) Stop() {
	if crawler_pool != nil {
		crawler_pool.Stop()
	}
}
