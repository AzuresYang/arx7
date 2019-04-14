/*
 * @Author: rayou
 * @Date: 2019-04-13 11:25:29
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-13 11:33:01
 */
package monitorCollector

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"
)

type monitorServer struct {
}

func (self *monitorServer) handlerEvent(c net.Conn) {
	var buf = make([]byte, 30)
	n, err := c.Read(buf)
	if err != nil {
		log.Debugf("conn read error:%s\n", err.Error())
		return
	}
	log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
}

func (self *monitorServer) Run() {
	code_info := "MonitorServer"
	log.Info("MonitorServer running")
	localAddress, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:8080") //定义一个本机IP和端口。
	var tcpListener, err = net.ListenTCP("tcp", localAddress)       //在刚定义好的地址上进监听请求。
	if err != nil {
		fmt.Println("监听出错：", err)
		return
	}
	defer func() { //担心return之前忘记关闭连接，因此在defer中先约定好关它。
		tcpListener.Close()
	}()
	for {
		c, err := tcpListener.Accept()
		if err != nil {
			log.Errorf("[%s]tcp listen err:%s", code_info, err.Error())
			continue
		}
		go self.handlerEvent(c)
	}
}
