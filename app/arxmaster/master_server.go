/*
 * @Author: rayou
 * @Date: 2019-04-13 11:25:29
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-13 21:40:03
 */
package arxmaster

import (
	"bytes"
	"net"

	log "github.com/sirupsen/logrus"
)

type masterServer struct {
}

func (self *masterServer) handlerEvent(c net.Conn) {
	code_info := "MasterServer"
	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil {
		log.Debugf("conn read error:%s\n", err.Error())
		return
	}
	log.Infof("[MasterServer]read %d bytes, content is %s\n", n, string(buf[:n]))
	echo := []byte("get you conn,echo:")
	var btBuf bytes.Buffer
	btBuf.Write(echo)
	btBuf.Write(buf[:n])
	echo = btBuf.Bytes()
	n, err = c.Write(echo)
	log.Infof("[%s]echo %d bytes, content is %s\n", code_info, n, string(echo[:n]))
}

func (self *masterServer) Run() {
	code_info := "MasterServer"
	log.Info("MasterServer running")
	localAddress, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:8888") //定义一个本机IP和端口。
	var tcpListener, err = net.ListenTCP("tcp", localAddress)       //在刚定义好的地址上进监听请求。
	if err != nil {
		log.Error("监听出错：", err.Error())
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
		log.Debugf("[%s]get conn from:%s", code_info, c.RemoteAddr())
		go self.handlerEvent(c)
	}
}
