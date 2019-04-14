/*
 * @Author: rayou
 * @Date: 2019-04-13 22:59:09
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-14 16:53:56
 */
package message

import (
	"fmt"
	"testing"
)

func TestMsg(t *testing.T) {
	msg := &CommMsg{
		Cmd: 1001,
		Id:  12346,
	}
	data := []byte("hello msg")
	msg.DataLen = uint64(len(data))
	// msg.Data = make([]byte, msg.Size)
	// for i, _ := range data {
	// 	msg.Data[i] = data[i]
	// }
	msg.Data = data
	fmt.Printf("before msg:%+v\n", msg)
	jbyte := msg.Serialize()
	new_msg, _ := UnSerialize(jbyte)
	fmt.Printf("after msg:%+v", new_msg)
	t.Log("donw")
}
