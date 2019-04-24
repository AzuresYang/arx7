/*
 * @Author: rayou
 * @Date: 2019-04-13 11:25:29
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 22:22:07
 */
package arxmaster

import (
	"github.com/AzuresYang/arx7/app/arxlet"
	"encoding/json"
	"github.com/AzuresYang/arx7/app/arxmaster/arxscheduler"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/arxmonitor"
	"github.com/AzuresYang/arx7/app/arxmonitor/monitorCollector"
	"github.com/AzuresYang/arx7/db"
	log "github.com/sirupsen/logrus"
)

type ArxMaster struct {
	server     *arxlet.BaseTcpServer
	listenPort string
	scheduler  arxscheduler.ArxScheduler
	mc  *monitorCollector.MonitorCollector // 监控组件
}

func NewArxMaster() *ArxMaster {
	master := &ArxMaster{
		server: arxlet.NewBaseTcpServer(),
		mc : monitorCollector.New(1000),
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
	cmds := self.scheduler.GetSupportCmds()
	self.server.RegisterHandler(cmds, &self.scheduler)
	cmds = self.GetSupportCmds()
	self.server.RegisterHandler(cmds,self)
	return nil
}

// 开始收集监控
func (self *ArxMaster) StartMonitorCollector(db_cfg *db.MysqlConfig){
	self.mc.Start(db_cfg)
}

func (self *ArxMaster) Run() {
	self.server.Run()
}

func (self *ArxMaster) Stop() {
}

func (self *ArxMaster)GetSupportCmds()[]uint32{
	return []uint32{
		message.MSG_MONITOR_INFO,
	}
}
func (self *ArxMaster) HandlerEvent(ctx *arxlet.ConnContext) {
	code_info := "ArxMaster.HandlerEvent"
	client := ctx.From.RemoteAddr().String()
	log.WithFields(log.Fields{
		"line": code_info,
		"addr": client,
		"cmd":  ctx.Msg.Cmd,
	}).Info("recv event.")
	switch ctx.Msg.Cmd {
	case message.MSG_MONITOR_INFO:
		log.Tracef("[%s] [%s]get monitor info.", code_info, client)
		pkg := &arxmonitor.MonitorMsgPkg{}
		err := json.Unmarshal(ctx.Msg.Data, pkg)
		if err != nil{
			log.Errorf("[%s] [%s]unmarshal MonitorPkg fail.",code_info, client)
			return
		}
		self.mc.AddMonitorPkg(pkg)
	}
}