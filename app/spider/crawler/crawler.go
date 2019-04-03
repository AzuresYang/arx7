/*
 * @Author: rayou
 * @Date: 2019-03-25 22:21:15
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-03 23:00:25
 */
package crawler

import (
	"github.com/AzuresYang/arx7/app/pipeline"
	"github.com/AzuresYang/arx7/app/processor"
	"github.com/AzuresYang/arx7/app/spider/downloader"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
	log "github.com/sirupsen/logrus"
	"github.com/AzuresYang/arx7/config"
	
)

type (
	Crawler interface {
		// Init( /*一个spider分析器，配置文件*/ ) error // 初始化
		Run()       // 运行
		Stop()      // 停止运行
		GetId() int // 获取ID
	}

	crawler struct {
		// 一个采集规则分析器spider,下载器 downloader, 存储pipleLine
		procer     processor.Processor
		downloader downloader.Downloader
		pipeline   pipeline.Pipeline
		id         int // id
		if_stop    bool
		pause      [2]int64 //[距离下个请求的最短时常， 距离下个请求的最长时长]
	}
)

// 新建一个Crawler
func NewCrawler(id int) Crawler {
	return &crawler{
		id:      id,
		if_stop: false,
	}
}

func (self *crawler) Init() error {
	pipe_err := self.pipelien.Init()
	if pipe_err != nil {
		log.Errorlf("crawler[%d]init pipeline[%s] fail.errmsg:%s", self.GetId(), self.pipeline.GetName(), pipe_err.Error())
		return pipe_err
	}
	procer_err = self.procer.Init()
	if procer_err != nil {
		log.Errorlf("crawler[%d]init processor[%s] fail.errmsg:%s", self.GetId(), self.procer.GetName(), procer_err.Error())
		return procer_err
	}
	return nil
}

func (self *crawler) Run() {
	// self.Init()+
	for !self.if_stop {
		req := request.RequestMgr.GetRequest(config.CrawlerCfg.RequestGetTimeOut)
		// 没有新链接的时候，等一段时间继续获取
		if req == nil{
			log.Error("Get reqeust time out. Stop Crawler")
			time.Sleep(config.DEFAULT_REQ_IS_NULL_WAITTIME)
			continue
		}
		go self.processRequest(req) // 下载链接
		self.sleep()
	}

}

func (self *crawler) Stop() {
	self.if_stop = true
}

func (self *crawler) GetId() int {
	return self.id
}

func (self *crawler) processRequest(req *request.ArxRequest){

}

func (self *crawler) sleep(){
	sleep_time := self.pause[0] + rand.Int63n(self.pause[1])
	time.Sleep(time.Duration(sleep_time) * time.Millisecond)
}