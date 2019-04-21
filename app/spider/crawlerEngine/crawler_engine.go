/*
 * @Author: rayou
 * @Date: 2019-04-14 20:42:23
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-15 21:28:39
 */
package crawlerEngine

import (
	"errors"

	"reflect"

	"github.com/AzuresYang/arx7/app/pipeline/output"
	"github.com/AzuresYang/arx7/app/spider/downloader"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
	"github.com/AzuresYang/arx7/app/status"
	"github.com/AzuresYang/arx7/config"
	log "github.com/sirupsen/logrus"
)

type CrawlerEngine struct {
	crawlerPool   CrawlerPool
	crawlerNum    int
	state         int
	fastDfsOutput output.OutputFastDfs
}

func NewCrawlerEngine(crawlerNum int) *CrawlerEngine {
	if crawlerNum <= 0 {
		crawlerNum = 1
	}
	engine := &CrawlerEngine{
		crawlerPool: NewCrawlerPool(uint32(crawlerNum)),
		crawlerNum:  crawlerNum,
		state:       status.STOP,
	}
	return engine
}

func (self *CrawlerEngine) Init(cfg *config.CrawlerTask) error {
	// 初始化链接管理器
	code_info := reflect.TypeOf(self).String()
	err := request.RequestMgr.Init(cfg)
	if err != nil {
		log.Errorf("[%s]init request Mgr fail fail:%s", code_info, err.Error())
		return errors.New("request mgr init fail")
	}
	self.fastDfsOutput.Reset(cfg.FastDfsAddr, cfg.TaskName)
	err = self.fastDfsOutput.Init()
	if err != nil {
		log.Errorf("[%s]init output to fast dfs fail:%s", code_info, err.Error())
		return errors.New("init output to fast dfs fail")
	}
	log.Info("CralerEngine init succ")
	// 初始化一些东西，比如requestMgr的
	return nil
}

func (self *CrawlerEngine) Run() {
	log.Info("crawler engine run start")
	self.state = status.RUN
	for i := 0; i < self.crawlerNum; i++ {
		cw := self.crawlerPool.Get()
		// pipe := &output.OutputFile{}
		dl := &downloader.SimpleDownloader{}
		cw.Init(dl, &self.fastDfsOutput)
		go cw.Run()
	}
}

func (self *CrawlerEngine) Stop() {
	log.Info("CrawlerEngine ready Stop")
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
