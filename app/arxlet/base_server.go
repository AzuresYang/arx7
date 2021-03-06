/*
 * @Author: rayou
 * @Date: 2019-04-14 17:25:13
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-22 22:34:23
 */

package arxlet

import (
	"encoding/json"
	"net"
	"reflect"

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
	GetSupportCmds() []uint32
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

// 注册处理的handler
// cmd的注册方式，如果已经存在相同的命令，可以报个错误
func (self *BaseTcpServer) RegisterHandler(cmds []uint32, handler ConnectHandler) {
	// cmds := handler.GetSupportCmds()
	for _, cmd := range cmds {
		self.handlers[cmd] = handler
		log.Debugf("[BaseTcpServer.RegisterHandler]register cmd[%d] to %s", cmd, reflect.TypeOf(handler).String())
	}
}

// 开始监听端口
func (self *BaseTcpServer) Init(listenport string) error {
	address := ":" + listenport
	localAddress, _ := net.ResolveTCPAddr("tcp4", address) //定义一个本机IP和端口。
	var err error
	self.tcpListener, err = net.ListenTCP("tcp", localAddress) //在刚定义好的地址上进监听请求。
	if err != nil {
		log.Error("tcp server listen tcp err：", err.Error())
		return err
	}
	log.Infof("[BaseTcpServer]init listen addrss succ:%s\n", address)
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

// 收到消息后，根据msg.Cmd来找到对应的handler,进行处理
func (self *BaseTcpServer) dispatchEvent(c net.Conn) {
	code_info := "BaseTcpServer.dispatchEvent"
	log.Tracef("[%s] get conn:%s", code_info, c.RemoteAddr().String())
	var buf = make([]byte, max_buffer_size)
	n, err := c.Read(buf)
	if err != nil {
		log.Errorf("[%s]conn read data error:%s\n", code_info, err.Error())
		return
	}
	ctx := &ConnContext{
		From: c,
	}
	err = json.Unmarshal(buf[:n], &ctx.Msg)
	if err != nil {
		log.Errorf("[%s]conn data cant unmarshal to CommMsg:%s\n", code_info, err.Error())
		return
	}
	handler := self.handlers[ctx.Msg.Cmd]
	if handler == nil {
		log.Errorf("[%s] cannt found handler:%d.register cmd:%+v", code_info, ctx.Msg.Cmd, self.handlers)
	}
	handler.HandlerEvent(ctx)
	// 记得关闭连接
	c.Close()
}
