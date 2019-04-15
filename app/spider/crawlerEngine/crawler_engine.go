/*
 * @Author: rayou
 * @Date: 2019-04-14 20:42:23
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-15 21:28:39
 */
package crawlerEngine

import (
	"github.com/AzuresYang/arx7/app/pipeline/output"
	"github.com/AzuresYang/arx7/app/spider/downloader"
	"github.com/AzuresYang/arx7/app/status"
	log "github.com/sirupsen/logrus"
)

type CrawlerEngine struct {
	crawlerPool CrawlerPool
	crawlerNum  int
	state       int
}

func NewCrawlerEngine(crawlerNum int) *CrawlerEngine {
	engine := &CrawlerEngine{
		crawlerPool: NewCrawlerPool(uint32(crawlerNum)),
		crawlerNum:  crawlerNum,
		state:       status.STOP,
	}
	return engine
}

func (self *CrawlerEngine) Init() error {
	// 初始化一些东西，比如requestMgr的
	return nil
}

func (self *CrawlerEngine) Run() {
	log.Info("engin run")
	self.state = status.RUN
	// request.RequestMgr.Init(20)
	// 包括crawler 池的初始化 ？？？
	for i := 0; i < self.crawlerNum; i++ {
		cw := self.crawlerPool.Get()
		pipe := &output.OutputFile{}
		dl := &downloader.SimpleDownloader{}
		cw.Init(dl, pipe)
		go cw.Run()
	}
}

func (self *CrawlerEngine) Stop() {
	self.state = status.STOP
	if self.crawlerPool != nil {
		self.crawlerPool.Stop()
	}
}

func (self *CrawlerEngine) IfStop() bool {
	if self.state == status.STOP {
		return true
	}
	return false
}
