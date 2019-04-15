/*
 * @Author: rayou
 * @Date: 2019-04-15 17:22:22
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-15 21:15:43
 */

package arxscheduler

import (
	"encoding/json"
	"net"
	"time"

	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/status"
	"github.com/AzuresYang/arx7/config"
	log "github.com/sirupsen/logrus"
)

type ArxScheduler struct {
	MasterListenPort string
}

func NewArxScheduler(masterPort string) *ArxScheduler {
	return &ArxScheduler{
		MasterListenPort: masterPort,
	}
}

func (self *ArxScheduler) Init(masterPort string) {
	self.MasterListenPort = masterPort
}

func (self *ArxScheduler) GetSupportCmds() []uint32 {
	cmds := []uint32{
		message.MSG_ARXCMD_START_SPIDER,
		message.MSG_ARXCMD_STOP_STOP,
		message.MSG_ARXCMD_SCALE_STOP,
		message.MSG_ARXCMD_DELTE_TASK,
		message.MSG_ARXCMD_GET_SPIDER_INFO,
	}
	return cmds
}
func (self *ArxScheduler) HandlerEvent(ctx *arxlet.ConnContext) {
	code_info := "ArxSheduler.HandlerEvent"
	log.WithFields(log.Fields{
		"line": code_info,
		"addr": ctx.From.RemoteAddr().String(),
		"cmd":  ctx.Msg.Cmd,
	}).Info("recv event.")
	switch ctx.Msg.Cmd {
	case message.MSG_ARXCMD_START_SPIDER:
		self.startSpider(ctx)
	case message.MSG_MONITOR_INFO:
		log.Info("recv monitor_info")
	}
}

func (self *ArxScheduler) startSpider(ctx *arxlet.ConnContext) error {
	code_info := "ArxScheduler.startSpider"
	cfg := self.generateCrawlerTask(ctx)
	// resp := message.NewResponseMsg(0, "succ")
	data, err := json.Marshal(cfg)
	if err != nil {
		log.Errorf("[%s] serialize crawlerTask fail:%s", code_info, err.Error())
		resp := message.NewResponseMsg(status.ERR_UNSERIALIZE_FAIL, "generate crawler task error.")
		return responseConn(ctx, resp)
	}
	var spider_conn net.Conn
	// 这个端口是测试用的客户端地址
	err, spider_conn = arxlet.SendTcpMsgTimeoutWithConn(message.MSG_REQ_STAET_SPIDER, data, ":9888", 2*time.Second)
	if err != nil {
		log.Errorf("[%s] connect to spider fail:%s", code_info, err.Error())
		resp := message.NewResponseMsg(status.ERR_UNSERIALIZE_FAIL, "connect to spider fail.")
		return responseConn(ctx, resp)
	}
	var resp *message.ResponseMsg
	err, resp = arxlet.ParseResponseFromConn(spider_conn)
	if err != nil {
		log.Errorf("[%s] parseResponse Error:%s", code_info, err.Error())
		return err
	}
	log.Infof("[%s] start spider.ret:%d, msg:%s", code_info, resp.Status, resp.Msg)
	return nil
}

func (self *ArxScheduler) generateCrawlerTask(ctx *arxlet.ConnContext) *config.CrawlerTask {
	cfg := &config.CrawlerTask{
		TaskName:         "test",
		TaskId:           11,
		CrawlerTreadNum:  3,
		RedisAddress:     "",
		RedisPort:        8888,
		RedisAccount:     "auth",
		RedisPassword:    "",
		RedisQueueName:   "test",
		MasterListenPort: self.MasterListenPort,
	}
	return cfg
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
