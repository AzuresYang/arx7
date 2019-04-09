/*
 * @Author: rayou
 * @Date: 2019-04-07 17:15:03
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-09 22:28:05
 */

package monitor

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/AzuresYang/arx7/util/httpUtil"
	log "github.com/sirupsen/logrus"
)

const (
	default_msg_send_num           int           = 50  // 默认最大监控数量发送
	default_msg_send_time          int64         = 300 // 默认发送时间
	default_msg_send_time_duration time.Duration = 5 * time.Second
	default_max_msg_num                          = 1000 // chan 最多存多少个监控数据
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
		Ip                  string
		Port                uint32
		LocalIp             string
		m                   sync.Mutex
		MsgList             *list.List
		Msgs                chan *MonitorMsg
		MaxMsgSendNum       int   // 达到这个数字之后就立刻开始发送监控数据
		MsgSendTime         int64 // 多久开始发送一次监控数据
		MsgSendTimeDuration time.Duration
		LastMsgSendTime     int64
		If_Stop             bool
	}
)

var monitor_handler *monitorHandler

func InitMonitorHandler(ip string, port uint32, svcid uint32) error {
	monitor_handler = &monitorHandler{
		Ip:                  ip,
		Port:                port,
		SvcId:               svcid,
		MsgList:             list.New(),
		MaxMsgSendNum:       default_msg_send_num,
		Msgs:                make(chan *MonitorMsg, default_max_msg_num),
		MsgSendTime:         default_msg_send_time,
		MsgSendTimeDuration: default_msg_send_time_duration,
		LastMsgSendTime:     0,
	}
	local_ip, err := httpUtil.GetLocalIp()
	if err != nil {
		log.Error("cant get local ip")
		local_ip = "127.0.0.1"
	}
	monitor_handler.LocalIp = local_ip
	fmt.Println("start send monitor msg")
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

func collectMsg() {
	click := time.After(monitor_handler.MsgSendTimeDuration)
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
	msg_pkg := newMonitorMsgPkg(monitor_handler.LocalIp, monitor_handler.MsgList.Len())
	for i := monitor_handler.MsgList.Front(); i != nil; i = i.Next() {
		msg := i.Value.(*MonitorMsg)
		msg_pkg.Msgs = append(msg_pkg.Msgs, msg)
	}
	// 初始化清空监控数据列表
	monitor_handler.MsgList.Init()
	monitor_handler.m.Unlock()
	go doSendMsg(msg_pkg)
	// 清空所有
}

// 发送监控数据包
func doSendMsg(pkg *MonitorMsgPkg) {
	for i, _ := range pkg.Msgs {
		fmt.Printf("[%d]\n", i)
	}
	fmt.Printf("msg list num:%d", monitor_handler.MsgList.Len())
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
	msg := newMonitorMsg(monitor_handler.SvcId, metric, MONITORMSG_ADD)
	addMonitorMsg(msg)
}

func AddOneWithClassfy(metric uint32, classfy uint32) {
	msg := newMonitorMsg(monitor_handler.SvcId, metric, MONITORMSG_ADD)
	msg.Classfy = classfy
	addMonitorMsg(msg)
}

func Set(metric uint32, num uint32) {
	msg := newMonitorMsg(monitor_handler.SvcId, metric, MONITORMSG_SET)
	msg.Value = num
	addMonitorMsg(msg)
}

func SetWithClassfy(metric uint32, classfy uint32, num uint32) {
	msg := newMonitorMsg(monitor_handler.SvcId, metric, MONITORMSG_SET)
	msg.Classfy = classfy
	msg.Value = num
	addMonitorMsg(msg)
}

func addMonitorMsg(msg *MonitorMsg) {
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
