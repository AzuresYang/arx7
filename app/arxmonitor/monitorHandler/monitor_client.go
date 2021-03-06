/*
 * @Author: rayou
 * @Date: 2019-04-07 17:15:03
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 22:23:34
 */

package monitorHandler

import (
	"container/list"
	"encoding/json"
	"sync"
	"time"

	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/arxmonitor"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/util/httpUtil"
	log "github.com/sirupsen/logrus"
)

const (
	default_msg_send_num           int           = 5               // 默认最大监控数量发送
	default_msg_send_time          int64         = 300             // 默认发送时间
	default_msg_send_time_interval time.Duration = 5 * time.Second // 测试的时候用5s
	default_max_msg_num                          = 1000            // chan 最多存多少个监控数据
)

type (
	// MonitorHandler interface {
	// 	init(ip string, port uint32, svcid uint32)       // 初始化服务器地址
	// 	AddOne(metric uint32)                            // 监控项目加1
	// 	AddOneWithClassfy(metric uint32, classfy uint32) // 监控项目加1， 带分类ID
	// 	Set(metric uint32, uint32 num)
	// 	SetWithClassfy(metric uint32, classfy uint32, uint32 num)
	// }
	monitorHandler struct {
		SvcId               uint32
		MasterAddr          string
		LocalIp             string
		m                   sync.Mutex
		MsgList             *list.List
		Msgs                chan *arxmonitor.MonitorMsg
		MaxMsgSendNum       int // 达到这个数字之后就立刻开始发送监控数据
		MsgSendTimeInterval time.Duration
		LastMsgSendTime     int64
		If_Stop             bool
	}
)

var monitor_handler *monitorHandler = &monitorHandler{
	MsgList:             list.New(),
	MaxMsgSendNum:       default_msg_send_num,
	Msgs:                make(chan *arxmonitor.MonitorMsg, default_max_msg_num),
	MsgSendTimeInterval: default_msg_send_time_interval,
	LastMsgSendTime:     0,
}

func InitMonitorHandler(masterAddr string, svcid uint32) error {
	monitor_handler.MasterAddr = masterAddr
	monitor_handler.SvcId = svcid
	local_ip, err := httpUtil.GetLocalIp()
	if err != nil {
		log.Error("cant get local ip")
		local_ip = "127.0.0.1"
	}
	monitor_handler.LocalIp = local_ip
	log.Infof("[MonitorHandler]start send monitor msg")
	monitor_handler.If_Stop = false
	go loopSendMsg()
	return nil
}

func loopSendMsg() {
	log.Info("start send monitor msg")
	for !monitor_handler.If_Stop {
		collectMsg()
		sendMsg()
	}
	log.Info("end send monitor msg")
}

// 定期或者达到一定数量就发送监控消息
func collectMsg() {
	click := time.After(monitor_handler.MsgSendTimeInterval)
	for {
		select {
		case msg := <-monitor_handler.Msgs:
			monitor_handler.m.Lock()
			monitor_handler.MsgList.PushBack(msg)
			if monitor_handler.MsgList.Len() >= monitor_handler.MaxMsgSendNum {
				monitor_handler.m.Unlock()
				return
			}
			monitor_handler.m.Unlock()
		case <-click:
			return
		}
	}

}

func sendMsg() {
	monitor_handler.m.Lock()
	// 没有数据则不需要发送
	if monitor_handler.MsgList.Len() <= 0 {
		monitor_handler.m.Unlock()
		return
	}
	msg_pkg := arxmonitor.NewMonitorMsgPkg(monitor_handler.LocalIp, monitor_handler.MsgList.Len())
	for i := monitor_handler.MsgList.Front(); i != nil; i = i.Next() {
		msg := i.Value.(*arxmonitor.MonitorMsg)
		msg_pkg.Msgs = append(msg_pkg.Msgs, *msg)
	}
	// 初始化清空监控数据列表
	monitor_handler.MsgList.Init()
	monitor_handler.m.Unlock()
	go doSendMsg(msg_pkg)
	// 清空所有
}

// 发送监控数据包
func doSendMsg(pkg *arxmonitor.MonitorMsgPkg) {
	code_info := "MonitorHandler.doSendMsg"
	// for i, _ := range pkg.Msgs {
	// 	fmt.Printf("[%d]\n", i)
	// }
	// fmt.Printf("msg list num:%d", monitor_handler.MsgList.Len())
	data, err := json.Marshal(pkg)
	if err != nil {
		log.Errorf("[%s] serialize msg pkg fail.", code_info)
		return
	}
	err = arxlet.SendTcpMsgTimeout(message.MSG_MONITOR_INFO, data, monitor_handler.MasterAddr, 5*time.Second)
	if err != nil {
		log.Errorf("[%s] send msg pkg fail:%s", code_info, err.Error())
		return
	}
	log.Tracef("[%s] send monitor msg pkg succ.", code_info)
}

// func (self *monitorHandler) init(ip string, port uint32, svcid uint32) {
// 	self.SvcId = svcid
// 	self.Ip = ip
// 	self.Port = port

// }

func Stop() {
	monitor_handler.If_Stop = true
}
func IfStop() bool {
	return monitor_handler.If_Stop
}
func AddOne(metric uint32) {
	msg := arxmonitor.NewMonitorMsg(monitor_handler.SvcId, metric, arxmonitor.MONITORMSG_ADD)
	addMonitorMsg(msg)
}

func AddOneWithClassfy(metric uint32, classfy uint32) {
	msg := arxmonitor.NewMonitorMsg(monitor_handler.SvcId, metric, arxmonitor.MONITORMSG_ADD)
	msg.Classfy = classfy
	addMonitorMsg(msg)
}

func Set(metric uint32, num uint32) {
	msg := arxmonitor.NewMonitorMsg(monitor_handler.SvcId, metric, arxmonitor.MONITORMSG_SET)
	msg.Value = num
	addMonitorMsg(msg)
}

func SetWithClassfy(metric uint32, classfy uint32, num uint32) {
	msg := arxmonitor.NewMonitorMsg(monitor_handler.SvcId, metric, arxmonitor.MONITORMSG_SET)
	msg.Classfy = classfy
	msg.Value = num
	addMonitorMsg(msg)
}

// 没有启动的时候不会添加监控数据
func addMonitorMsg(msg *arxmonitor.MonitorMsg) {
	if monitor_handler.If_Stop {
		return
	}
	monitor_handler.Msgs <- msg
	/*
		// 到达发送数量或者发送时间
		if (monitor_handler.MsgList.Len() >= monitor_handler.MaxMsgSendNum-1) ||
			(time.Now().Unix()-monitor_handler.LastMsgSendTime >= monitor_handler.MsgSendTime) {
			// 要先解锁, 不然没有办法发送
			monitor_handler.m.Unlock()
			go sendMonitorMsg()
		}
	*/
}
