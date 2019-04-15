/*
 * @Author: rayou
 * @Date: 2019-04-14 16:39:24
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-15 21:01:10
 * @breif : 负责通信，包括接受数据和发送数据
 */

package arxlet

import (
	"encoding/json"
	"net"
	"time"

	"github.com/AzuresYang/arx7/app/message"
	log "github.com/sirupsen/logrus"
)

const (
	default_dail_timeout time.Duration = 2 * time.Second // 链接超时时间
)

// 负责发送，优化的话，可以做成发送队列
type ArxSender struct {
}

func SendTcpMsgTimeoutWithConn(cmd uint32, data []byte, address string, dail_timeout time.Duration) (error, net.Conn) {
	msg := message.NewCommMsg(cmd, data)
	msg.DataLen = uint64(len(msg.Data))
	send_bytes, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("[ArxSender.SendTcpMsg] marshal msg json fail:%s", err.Error())
		return err, nil
	}
	conn, err := net.DialTimeout("tcp", address, dail_timeout)
	if err != nil {
		log.Errorf("dial server[%s] error:%s", address, err.Error())
		return err, nil
	}
	log.Trace("[ArxSender.SendTcpMsgTimeout]dail to svr succ:", address)
	_, err = conn.Write(send_bytes)
	if err != nil {
		log.Errorf("[ArxSender.SendTcpMsgTimeout]write to adress[%s]error:%s", address, err.Error())
		return err, conn
	}
	return nil, conn
}

// address格式：
// (1)ip:port
// (2):port   链接到本机
// @breif 发送一个CommMsg类型的消息， 带超时时间
func SendTcpMsgTimeout(cmd uint32, data []byte, address string, dail_timeout time.Duration) error {
	err, conn := SendTcpMsgTimeoutWithConn(cmd, data, address, dail_timeout)
	if conn != nil {
		conn.Close()
	}
	return err
}

// @breif 发送一个CommMsg类型的消息， 使用默认超时时间
func SendTcpMsg(cmd uint32, data []byte, address string) error {
	return SendTcpMsgTimeout(cmd, data, address, default_dail_timeout)
}

func ParseResponseFromConn(c net.Conn) (error, *message.ResponseMsg) {
	var buf = make([]byte, max_buffer_size)
	n, err := c.Read(buf)
	if err != nil {
		log.Errorf("[Arxlet.ParseAsMsgFromConn]conn read data error:%s\n", err.Error())
		return err, nil
	}
	resp := new(message.ResponseMsg)
	err = json.Unmarshal(buf[:n], resp)
	if err != nil {
		log.Errorf("[Arxlet.ParseAsMsgFromConn]conn data cant unmarshal to responseMsg:%s\n", err.Error())
		return err, nil
	}
	return nil, resp
}
