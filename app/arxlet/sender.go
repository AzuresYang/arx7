/*
 * @Author: rayou
 * @Date: 2019-04-14 16:39:24
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-14 21:46:56
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

// address格式：ip:port
// :port则链接到本机
func SendTcpMsgTimeout(cmd uint32, data []byte, address string, dail_timeout time.Duration) error {
	msg := message.NewCommMsg(cmd, data)
	msg.DataLen = uint64(len(msg.Data))
	send_bytes, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("[ArxSender.SendTcpMsg] marshal msg json fail:%s", err.Error())
		return err
	}
	log.Trace("ready conn server")
	conn, err := net.DialTimeout("tcp", address, dail_timeout)
	if err != nil {
		log.Errorf("dial server[%s] error:%s", address, err.Error())
		return err
	}
	log.Trace("[ArxSender.SendTcpMsgTimeout]conn server succ:", address)
	defer conn.Close()
	_, err = conn.Write(send_bytes)
	if err != nil {
		log.Errorf("[ArxSender.SendTcpMsgTimeout]write to adress[%s]error:%s", address, err.Error())
		return err
	}
	return nil
}

func SendTcpMsg(cmd uint32, data []byte, address string) error {
	return SendTcpMsgTimeout(cmd, data, address, default_dail_timeout)
}
