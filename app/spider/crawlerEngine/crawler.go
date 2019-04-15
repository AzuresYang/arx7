/*
 * @Author: rayou
 * @Date: 2019-03-25 22:21:15
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-15 15:59:26
 */
package crawlerEngine

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/AzuresYang/arx7/app/pipeline"
	"github.com/AzuresYang/arx7/app/processor"
	"github.com/AzuresYang/arx7/app/spider/downloader"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
	"github.com/AzuresYang/arx7/config"
	"github.com/AzuresYang/arx7/util/record"
	log "github.com/sirupsen/logrus"
)

type (
	Crawler interface {
		Init(dl downloader.Downloader, pipe pipeline.Pipeline) error // 初始化
		Run()                                                        // 运行
		Stop()
		IfStop() bool // 停止运行
		GetId() int   // 获取ID
	}

	crawler struct {
		// 一个采集规则分析器procer,这个应该是不用的,下载器 downloader, 存储pipleLine
		// procer     processor.Processor
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
		id:         id,
		if_stop:    false,
		pause:      [2]int64{1, 2},
		downloader: &downloader.SimpleDownloader{},
	}
}
func (self *crawler) IfStop() bool {
	return self.if_stop
}
func (self *crawler) SetDownloader(dl downloader.Downloader) {
	self.downloader = dl
}

func (self *crawler) Init(dl downloader.Downloader, pipe pipeline.Pipeline) error {
	self.downloader = dl
	self.pipeline = pipe

	pipe_err := self.pipeline.Init()
	if pipe_err != nil {
		log.Error(fmt.Sprintf("crawler[%d]init pipeline[%s] fail.errmsg:%s", self.GetId(), self.pipeline.GetName(), pipe_err.Error()))
		return pipe_err
	}
	/*
		procer_err := self.procer.Init()
		if procer_err != nil {
			log.Errorl(fmt.Sprintf("crawler[%d]init processor[%s] fail.errmsg:%s", self.GetId(), self.procer.GetName(), procer_err.Error()))
			return procer_err
		}
	*/
	return nil
}

func (self *crawler) Stop() {
	self.if_stop = true
}

func (self *crawler) GetId() int {
	return self.id
}

// 主方法
func (self *crawler) Run() {
	// err := self.init()
	// if err != nil {
	//	log.Errorf("crawler init fail.err:%s", err.Error())
	//	return
	//}
	self.if_stop = true
	self.run()
	self.Stop()
}

func (self *crawler) run() {
	// self.Init()+
	// 不断获取链接，下载，处理
	var max_req_null_times int = int(config.CrawlerCfg.TaskConf.MaxGetRequestNullTimeSecond /
		config.CrawlerCfg.RequestGetTimeOut)
	get_req_null_times := 0
	for !self.if_stop {
		req := request.RequestMgr.GetRequest(config.CrawlerCfg.RequestGetTimeOut)
		if req == nil {
			// 太长时间没有新链接的话，自动停止工作
			get_req_null_times += 1
			if max_req_null_times != 0 && get_req_null_times >= max_req_null_times {
				log.Errorf("Get reqeust time out. Stop Crawler:%d", self.id)
				self.if_stop = true
				break
			} else {
				log.Tracef("crawler[%d]get request nil.times:%d|%d", self.id, get_req_null_times, max_req_null_times)
				// 没有新链接的时候，等一段时间继续获取
				time.Sleep(config.DEFAULT_REQ_IS_NULL_WAITTIME)
				continue
			}
		}
		get_req_null_times = 0
		go self.processRequest(req) // 下载链接
		self.sleep()
	}
}

func (self *crawler) processRequest(req *request.ArxRequest) {
	procer := processor.Manager.GetProcessor(req.ProcerName)
	is_download_succ := false
	// 统计下载情况

	defer func() {
		if is_download_succ {
			record.CountAddOne(record.COUNT_DOWNLOAD_SUCC)
		} else {
			record.CountAddOne(record.COUNT_DOWNLOAD_FAIL)
		}
	}()

	if procer == nil {
		log.Error(fmt.Sprintf("req could found procer[%s]", req.ProcerName))
		processor.Manager.PrintAllProcessor("no found procer")
		record.DownloadSuccReq(req, "req could found procer:"+req.ProcerName)
		return
	}
	ctx := self.downloader.Download(procer, req)
	// 处理下载数据
	if ctx == nil {
		log.Error(fmt.Sprintf("download fail.ctx is nil.req:%#v", req))
		record.DownloadFailReq(req, "download fail. ctx is nil")
		return
	}
	log.Info("download succ.URL:" + req.Url)
	ret, msg := procer.Process(ctx)
	log.Info(fmt.Sprintf("process ret[%d|%s], URL:[%s]", ret, msg, req.Url))
}

func (self *crawler) sleep() {
	sleep_time := self.pause[0] + rand.Int63n(self.pause[1])
	time.Sleep(time.Duration(sleep_time) * time.Millisecond)
}
