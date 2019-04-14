/*
 * @Author: rayou
 * @Date: 2019-04-13 21:15:43
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-13 21:41:54
 */

package arxmaster

import (
	"net"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func buildClient() {
	log.Info("ready conn server")
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Error("dial error:", err)
		return
	}
	log.Info("conn server succ")
	defer conn.Close()
	data := "hello tcp"
	var n int
	n, err = conn.Write([]byte(data))
	if err != nil {
		log.Error("write error:", err.Error())
	}
	log.Infof("write byte:%d", n)
	buf := make([]byte, 1024)
	n, err = conn.Read(buf)
	if err != nil {
		log.Error("write error:", err.Error())
	} else {
		log.Infof("get echo:%s", string(buf[:n]))
	}
}
func TestServer(t *testing.T) {
	server := &masterServer{}
	log.SetLevel(log.TraceLevel)
	go server.Run()
	time.Sleep(2 * time.Second)
	go buildClient()
	time.Sleep(2 * time.Second)
	t.Log("donw")
}
