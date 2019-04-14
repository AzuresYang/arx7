/*
 * @Author: rayou
 * @Date: 2019-04-14 17:25:13
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-14 21:44:02
 */

package arxlet

import (
	"encoding/json"
	"net"

	"github.com/AzuresYang/arx7/app/message"
	log "github.com/sirupsen/logrus"
)

const (
	max_buffer_size int = 64000
)

type ConnContext struct {
	From net.Conn
	Msg  message.CommMsg
}
type ConnectHandler interface {
	HandlerEvent(ctx *ConnContext)
}

type BaseTcpServer struct {
	tcpListener *net.TCPListener
	handlers    map[uint32]ConnectHandler
}

func NewBaseTcpServer() *BaseTcpServer {
	server := &BaseTcpServer{
		handlers: make(map[uint32]ConnectHandler),
	}
	return server
}
func (self *BaseTcpServer) RegisterHandler(cmds []uint32, handler ConnectHandler) {
	for _, cmd := range cmds {
		self.handlers[cmd] = handler
	}
}

func (self *BaseTcpServer) Init(listenport string) error {
	address := "127.0.0.1:" + listenport
	localAddress, _ := net.ResolveTCPAddr("tcp4", address) //定义一个本机IP和端口。
	var err error
	self.tcpListener, err = net.ListenTCP("tcp", localAddress) //在刚定义好的地址上进监听请求。
	if err != nil {
		log.Error("tcp server listen tcp err：", err.Error())
		return err
	}
	return nil
}

func (self *BaseTcpServer) Run() {
	code_info := "BaseTcpServer.Run"
	defer func() {
		//担心return之前忘记关闭连接，因此在defer中先约定好关它。
		self.tcpListener.Close()
	}()
	for {
		c, err := self.tcpListener.Accept()
		if err != nil {
			log.Errorf("[%s]tcp listen err:%s", code_info, err.Error())
			continue
		}
		go self.dispatchEvent(c)
	}
}

func (self *BaseTcpServer) dispatchEvent(c net.Conn) {
	log.Infof("[dispatch] get conn:%s", c.RemoteAddr().String())
	var buf = make([]byte, max_buffer_size)
	n, err := c.Read(buf)
	if err != nil {
		log.Errorf("conn read data error:%s\n", err.Error())
		return
	}
	ctx := &ConnContext{
		From: c,
	}
	err = json.Unmarshal(buf[:n], &ctx.Msg)
	if err != nil {
		log.Errorf("conn data cant unmarshal to CommMsg:%s\n", err.Error())
		return
	}
	handler := self.handlers[ctx.Msg.Cmd]
	if handler == nil {
		log.Errorf("[dispatchEvent] cannt found handler:%d", ctx.Msg.Cmd)
	}
	handler.HandlerEvent(ctx)
	// 记得关闭连接
	c.Close()
}
