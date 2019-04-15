/*
 * @Author: rayou
 * @Date: 2019-04-02 19:51:55
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-15 21:32:23
 */
package spider

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/arxmonitor/monitorHandler"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/spider/crawlerEngine"
	"github.com/AzuresYang/arx7/app/status"
	"github.com/AzuresYang/arx7/config"
	log "github.com/sirupsen/logrus"
)

type Spider struct {
	ce            *crawlerEngine.CrawlerEngine
	server        *arxlet.BaseTcpServer
	crawlerConfig config.CrawlerConfig
}

func NewSpider() *Spider {
	spider := &Spider{
		server: arxlet.NewBaseTcpServer(),
	}
	cmds := spider.GetSupportCmds()
	spider.server.RegisterHandler(cmds, spider)
	return spider
}

func (self *Spider) Init(listenport string) error {
	err := self.server.Init(listenport)
	if err != nil {
		log.Errorf("Spider init server fail:%s", err.Error())
		return err
	}
	return nil
}

func (self *Spider) Run() {
	self.server.Run()
}

func (self *Spider) Stop() {
	self.ce.Stop()
}
func (self *Spider) StopCrawler() {
	self.ce.Stop()
	// 为空大法好
	self.ce = nil
}

func (self *Spider) StartCrawler(cfg *config.CrawlerConfig) (error, uint32) {
	code_info := "Spider.StartCrawler"
	log.Infof("[%s] start crawler", code_info)
	// 之前有的话，先停止
	if self.ce != nil {
		self.ce.Stop()
		return errors.New("crawler is running"), status.ERR_START_SPIDER_FAIL_RUNNING
	}
	monitorHandler.InitMonitorHandler(cfg.MasterAddr, cfg.TaskConf.TaskId)
	self.ce = crawlerEngine.NewCrawlerEngine(int(cfg.TaskConf.CrawlerTreadNum))
	// 初始化引擎
	//
	err := self.ce.Init()
	if err != nil {
		log.Errorf("[%s]init cralwerEngine fail:%s", code_info, err.Error())
		return err, status.ERR_START_SPIDER_FAIL
	}
	self.ce.Run()
	return nil, 0
}

func (self *Spider) GetSupportCmds() []uint32 {
	cmds := []uint32{
		message.MSG_REQ_STAET_SPIDER,
		message.MSG_REQ_STOP_SPIDER,
	}
	return cmds
}
func (self *Spider) HandlerEvent(ctx *arxlet.ConnContext) {
	code_info := "Spider.HandlerEvent"
	log.WithFields(log.Fields{
		"line": code_info,
		"addr": ctx.From.RemoteAddr().String(),
		"cmd":  ctx.Msg.Cmd,
	}).Info("recv event.")
	switch ctx.Msg.Cmd {
	case message.MSG_REQ_STAET_SPIDER:
		self.handlerStartSpider(ctx)
	case message.MSG_REQ_STOP_SPIDER:
		self.StopCrawler()
	}
}

func (self *Spider) handlerStartSpider(ctx *arxlet.ConnContext) error {
	code_info := "Spider.handlerStartSpider"
	err := json.Unmarshal(ctx.Msg.Data[:], &self.crawlerConfig.TaskConf)
	resp := &message.ResponseMsg{
		Status: 0,
		Msg:    "succ",
	}
	if err != nil {
		log.Errorf("[handlerStartSpider] unmarshal task conf fail.%s\n", err.Error())
		resp.Status = status.ERR_UNSERIALIZE_FAIL
		resp.Msg = "unserialize task conf fail"
		return responseConn(ctx, resp)
	}
	self.crawlerConfig.MasterAddr = fmt.Sprintf("%s:%d",
		ctx.From.RemoteAddr().String(), self.crawlerConfig.TaskConf.MasterListenPort)
	var ret_code uint32
	err, ret_code = self.StartCrawler(&self.crawlerConfig)
	if err != nil {
		log.Errorf("[%s]start crawler fail:%s", code_info, err.Error())
		resp.Status = ret_code
		resp.Msg = err.Error()
	}
	return responseConn(ctx, resp)
}

func (self *Spider) handlerStopSpider(ctx *arxlet.ConnContext) error {
	code_info := "Spider.handlerStartSpider"
	resp := &message.ResponseMsg{
		Status: 0,
		Msg:    "succ",
	}
	self.StopCrawler()
	log.Infof("[%s]stip crawler succ", code_info)
	return responseConn(ctx, resp)
}

func (self *Spider) handlerGetSpiderInfo(ctx *arxlet.ConnContext) error {
	code_info := "Spider.handlerGetSpiderInfo"
	resp := &message.ResponseMsg{
		Status: 0,
		Msg:    "succ",
	}
	if self.ce == nil || self.ce.IfStop() {
		resp.Status = status.STOP
		resp.Msg = "spider stop"
	} else {
		resp.Status = status.RUN
		resp.Msg = "spider running"
		data, err := json.Marshal(self.crawlerConfig.TaskConf)
		if err != nil {
			resp.Status = 1000
			resp.Msg = "generate spider running info fail."
			return responseConn(ctx, resp)
		}
		resp.Data = data
	}
	log.Infof("[%s]generate spider info succ", code_info)
	return responseConn(ctx, resp)
}

func responseConn(ctx *arxlet.ConnContext, msg *message.ResponseMsg) error {
	code_info := "responseConn"
	send_bytes, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("[%s] marshal responseMsg fail:%s", code_info, err.Error())
		return err
	}
	_, err = ctx.From.Write(send_bytes)
	ctx.From.Close()
	return err
}
