/*
 * @Author: rayou
 * @Date: 2019-04-09 19:50:57
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-09 22:22:32
 */

package monitor

import (
	"encoding/json"
	"time"
)

type MonitorMsgType int

const (
	MONITORMSG_ADD MonitorMsgType = iota
	MONITORMSG_SET
)

// 监控消息
type MonitorMsg struct {
	SvcId   uint32 //服务id
	Metric  uint32 //	监控ID
	Classfy uint32 // 分类ID
	Value   uint32 // 监控值
	Time    int64  // 监控上报时间
	MsgType MonitorMsgType
}

type MonitorMsgPkg struct {
	Ip   string // 监控上报地址
	Msgs []*MonitorMsg
}

func newMonitorMsgPkg(ip string, msg_num int) *MonitorMsgPkg {
	pkg := &MonitorMsgPkg{
		Ip:   ip,
		Msgs: make([]*MonitorMsg, 0, msg_num),
	}
	return pkg
}

func (self *MonitorMsg) Serialize() string {
	json_byte, _ := json.Marshal(self)
	return string(json_byte[:])
}

// 反序列化
func UnSerialize(s string) (*MonitorMsg, error) {
	msg := new(MonitorMsg)
	return msg, json.Unmarshal([]byte(s), msg)
}

func newMonitorMsg(svcid uint32, metric uint32, msg_type MonitorMsgType) *MonitorMsg {
	msg := &MonitorMsg{
		SvcId:   svcid,
		Metric:  metric,
		Classfy: 0,
		Value:   1,
		Time:    time.Now().Unix(),
		MsgType: msg_type,
	}
	return msg
}
