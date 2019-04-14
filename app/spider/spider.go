/*
 * @Author: rayou
 * @Date: 2019-04-02 19:51:55
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-14 21:52:02
 */
package spider

import (
	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/spider/crawlerEngine"
	log "github.com/sirupsen/logrus"
)

type Spider struct {
	ce     *crawlerEngine.CrawlerEngine
	server *arxlet.BaseTcpServer
}

func NewSpider() *Spider {
	spider := &Spider{
		server: arxlet.NewBaseTcpServer(),
	}

	cmds := []uint32{message.MSG_REQ_STAET_SPIDER, message.MSG_REQ_STOP_SPIDER}
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

func (self *Spider) HandlerEvent(ctx *arxlet.ConnContext) {
	code_info := "Spider.HandlerEvent"
	log.WithFields(log.Fields{
		"line": code_info,
		"addr": ctx.From.RemoteAddr().String(),
		"cmd":  ctx.Msg.Cmd,
	}).Info("recv event.")
	switch ctx.Msg.Cmd {
	case message.MSG_REQ_STAET_SPIDER:
		log.Infof("[%s] recv start spider cmd.msg:%s\n", code_info, string(ctx.Msg.Data))
	}
}
