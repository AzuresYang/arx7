/*
 * @Author: rayou
 * @Date: 2019-04-02 19:51:55
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-21 21:17:10
 */
package spider

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/arxmonitor/monitorHandler"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/spider/crawlerEngine"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
	"github.com/AzuresYang/arx7/app/status"
	"github.com/AzuresYang/arx7/config"
	"github.com/AzuresYang/arx7/runtime"
	log "github.com/sirupsen/logrus"
)

type Spider struct {
	name          string
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
	self.name = listenport
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
	runtime.G_CrawlerCfg.TaskConf = cfg.TaskConf
	log.Infof("[%s_%s] start crawler", self.name, code_info)

	// 之前有的话，先停止
	if self.ce != nil {
		self.ce.Stop()
		return errors.New("crawler is running"), status.ERR_START_SPIDER_FAIL_RUNNING
	}
	// 初始化监控组件
	monitorHandler.InitMonitorHandler(cfg.MasterAddr, cfg.TaskConf.TaskId)

	self.ce = crawlerEngine.NewCrawlerEngine(int(cfg.TaskConf.CrawlerTreadNum))
	// 初始化引擎
	err := self.ce.Init(&cfg.TaskConf)
	if err != nil {
		log.Errorf("[%s_%s]init cralwerEngine fail:%s", self.name, code_info, err.Error())
		return err, status.ERR_START_SPIDER_FAIL
	}
	self.ce.Run()
	return nil, 0
}

func (self *Spider) GetSupportCmds() []uint32 {
	cmds := []uint32{
		message.MSG_REQ_STAET_SPIDER,
		message.MSG_REQ_STOP_SPIDER,
		message.MSG_REQ_GET_SPIDER_INFO,
		message.MSG_REG_ECHO,
		message.MSG_REG_ECHO_REDIS,
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
		self.handlerStopSpider(ctx)
	case message.MSG_REQ_GET_SPIDER_INFO:
		self.handlerGetStatus(ctx)
	case message.MSG_REG_ECHO:
		self.handlerEcho(ctx)
	case message.MSG_REG_ECHO_REDIS:
		self.handlerEchoRedis(ctx)
	}
}

func (self *Spider) handlerStartSpider(ctx *arxlet.ConnContext) error {
	code_info := self.name + "::" + "Spider.handlerStartSpider"
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
	self.crawlerConfig.MasterAddr = fmt.Sprintf("%s:%s",
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
	code_info := self.name + "::" + "Spider.handlerStartSpider"
	resp := &message.ResponseMsg{
		Status: 0,
		Msg:    "succ",
	}
	self.StopCrawler()
	log.Infof("[%s]stop crawler succ", code_info)
	return responseConn(ctx, resp)
}

func (self *Spider) handlerGetStatus(ctx *arxlet.ConnContext) error {
	code_info := self.name + "::" + "Spider.handlerGetStatus"
	log.Infof("[%s]recv event get spider status", code_info)
	resp := &message.ResponseMsg{
		Status: 0,
		Msg:    "",
	}
	// ce为空或者不在运行中，则返回stopped的结果
	if self.ce == nil {
		resp.Msg = "crawler stopped"
	} else if self.ce.IfStop() {
		resp.Msg = "crawler stopped"
	} else {
		conf_byte, err := json.Marshal(self.crawlerConfig.TaskConf)
		if err != nil {
			log.Errorf("[%s]generate task info fail:%s", code_info, err.Error())
			resp.Status = 1
			resp.Msg = "generate task info fail"
			return responseConn(ctx, resp)
		}
		resp.Msg = string(conf_byte)
	}
	return responseConn(ctx, resp)
}

func (self *Spider) handlerEcho(ctx *arxlet.ConnContext) error {
	code_info := self.name + "::" + "Spider.handlerEcho"
	t := time.Now()
	date := t.Format("2006-01-02 15:04:05")
	echo := fmt.Sprintf("[%s:%s]echo succ", self.name, date)
	resp := &message.ResponseMsg{
		Status: 0,
		Msg:    echo,
	}
	log.Infof("[%s]echo msg", code_info)
	return responseConn(ctx, resp)
}

// 向redis设置一个值，测试docker下的网路联通性
func (self *Spider) handlerEchoRedis(ctx *arxlet.ConnContext) error {
	code_info := self.name + "::" + "Spider.handlerEchoRedis"
	t := time.Now()
	date := t.Format("2006-01-02 15:04:05")
	resp := &message.ResponseMsg{
		Status: 0,
		Msg:    "",
	}
	echo := fmt.Sprintf("[%s:%s]echo redis", self.name, date)
	err := request.RequestMgr.SetKeyValue(self.crawlerConfig.TaskConf.RedisAddr,
		self.crawlerConfig.TaskConf.RedisPassword,
		"echo",
		echo)
	echo_ret := "succ"
	if err != nil {
		log.Errorf("[%s]echo redis fail", code_info)
		echo_ret = "fail"
	}
	echo = fmt.Sprintf("[%s:%s]echo redis:%s", self.name, date, echo_ret)
	resp.Msg = echo
	log.Infof("[%s]echo redis msg", code_info)
	return responseConn(ctx, resp)
}

func (self *Spider) handlerGetSpiderInfo(ctx *arxlet.ConnContext) error {
	code_info := self.name + "::" + "Spider.handlerGetSpiderInfo"
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
