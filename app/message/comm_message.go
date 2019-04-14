/*
 * @Author: rayou
 * @Date: 2019-04-13 22:51:20
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-14 16:57:03
 */
package message

import (
	"encoding/json"
	"time"
)

type CommMsg struct {
	Cmd     uint32 // Cmd，或者说是消息类型
	Id      uint64 // 序列号， 暂时可以不管
	DataLen uint64 // data的长度
	Data    []byte
}

type ResponseMsg struct {
	Status uint32
	Msg    string
	Data   []byte
}

func NewCommMsg(cmd uint32, data []byte) *CommMsg {
	msg := &CommMsg{
		Cmd:     cmd,
		DataLen: uint64(len(data)),
		Data:    data,
		Id:      uint64(time.Now().Unix()),
	}
	return msg
}

func (self *CommMsg) Serialize() []byte {
	json_byte, _ := json.Marshal(self)
	return json_byte
}

// 反序列化
func UnSerializeCommMsg(b []byte) (*CommMsg, error) {
	req := new(CommMsg)
	return req, json.Unmarshal(b, req)
}
