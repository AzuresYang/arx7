/*
 * @Author: rayou
 * @Date: 2019-04-13 11:25:29
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-15 21:15:55
 */
package arxmaster

import (
	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/arxmaster/arxscheduler"
	"github.com/AzuresYang/arx7/app/message"
	log "github.com/sirupsen/logrus"
)

type ArxMaster struct {
	server     *arxlet.BaseTcpServer
	listenPort string
	scheduler  arxscheduler.ArxScheduler
}

func NewArxMaster() *ArxMaster {
	master := &ArxMaster{
		server: arxlet.NewBaseTcpServer(),
	}
	return master
}
func (self *ArxMaster) Init(listenport string) error {
	self.listenPort = listenport
	err := self.server.Init(listenport)
	if err != nil {
		log.Errorf("ArxMaster init server fail:%s", err.Error())
		return err
	}
	// 记录端口
	self.scheduler.Init(listenport)
	// 注册命名字

	self.server.RegisterHandler(&self.scheduler)
	return nil
}

func (self *ArxMaster) Run() {
	self.server.Run()
}

func (self *ArxMaster) Stop() {
}

func (self *ArxMaster) HandlerEvent(ctx *arxlet.ConnContext) {
	code_info := "ArxMaster.HandlerEvent"
	log.WithFields(log.Fields{
		"line": code_info,
		"addr": ctx.From.RemoteAddr().String(),
		"cmd":  ctx.Msg.Cmd,
	}).Info("recv event.")
	switch ctx.Msg.Cmd {
	case message.MSG_REQ_STAET_SPIDER:
		log.Infof("[%s] recv start spider cmd.msg:%s\n", code_info, string(ctx.Msg.Data))
	case message.MSG_MONITOR_INFO:
	}
}
