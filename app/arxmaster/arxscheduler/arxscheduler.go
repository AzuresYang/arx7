/*
 * @Author: rayou
 * @Date: 2019-04-15 17:22:22
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-21 21:06:49
 */

package arxscheduler

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/AzuresYang/arx7/app/arxlet"
	"github.com/AzuresYang/arx7/app/message"
	"github.com/AzuresYang/arx7/app/status"
	"github.com/AzuresYang/arx7/config"
	log "github.com/sirupsen/logrus"
)

type ArxScheduler struct {
	MasterListenPort string
}

func NewArxScheduler(masterPort string) *ArxScheduler {
	return &ArxScheduler{
		MasterListenPort: masterPort,
	}
}

func (self *ArxScheduler) Init(masterPort string) {
	self.MasterListenPort = masterPort
}

func (self *ArxScheduler) GetSupportCmds() []uint32 {
	cmds := []uint32{
		message.MSG_ARXCMD_START_SPIDER,
	}
	return cmds
}
func (self *ArxScheduler) HandlerEvent(ctx *arxlet.ConnContext) {
	code_info := "ArxSheduler.HandlerEvent"
	log.WithFields(log.Fields{
		"line": code_info,
		"addr": ctx.From.RemoteAddr().String(),
		"cmd":  ctx.Msg.Cmd,
	}).Info("recv event.")
	switch ctx.Msg.Cmd {
	case message.MSG_ARXCMD_START_SPIDER:
		self.startSpider(ctx)
	case message.MSG_MONITOR_INFO:
		log.Info("recv monitor_info")
	}
}

func (self *ArxScheduler) startSpider(ctx *arxlet.ConnContext) error {
	code_info := "ArxScheduler.startSpider"
	// 获取启动信息
	start_info := message.SpiderStartMsg{}
	json.Unmarshal(ctx.Msg.Data, &start_info)

	// 生成启动任务
	cfg := self.generateCrawlerTask(&start_info)
	// resp := message.NewResponseMsg(0, "succ")
	cfg_data, err := json.Marshal(cfg)
	if err != nil {
		log.Errorf("[%s] serialize crawlerTask fail:%s", code_info, err.Error())
		resp := message.NewResponseMsg(status.ERR_UNSERIALIZE_FAIL, "generate crawler task error.")
		return responseConn(ctx, resp)
	}
	log.Infof("[%s] start spider:%+v", code_info, start_info.NodeAddrs)
	start_result := self.doStartSpider(start_info.NodeAddrs, cfg_data)
	log.Infof("[%s]start spider ret:\n%+v", code_info, start_result)

	var start_fail_num uint32 = 0
	for _, ret := range start_result {
		if ret != "" {
			start_fail_num++
		}
	}
	resp := message.NewResponseMsg(start_fail_num, "")
	resp.Data, _ = json.Marshal(start_result)
	return responseConn(ctx, resp)
}

// 向所有列表发送启动消息，
// 将所有列表的回复组装成一个map[address]string回复
func (self *ArxScheduler) doStartSpider(nodeAddrs []string, cfg_data []byte) map[string]string {
	code_info := "ArxScheduler.doStartSpider"
	var wg sync.WaitGroup
	var start_result map[string]string = make(map[string]string, len(nodeAddrs))
	for _, addr := range nodeAddrs {
		wg.Add(1)
		// 线程发送启动消息
		// 1.是否连接到节点
		// 2.是否有回复
		// 3.是否成功
		go func(addr string) {
			defer wg.Done()
			log.Debugf("[%s] start send task to node:%s\n", code_info, addr)
			err, spider_conn := arxlet.SendTcpMsgTimeoutWithConn(message.MSG_REQ_STAET_SPIDER, cfg_data, addr, 2*time.Second)
			if err != nil {
				log.Errorf("[%s] connect to spider fail:%s", code_info, err.Error())
				start_result[addr] = "connect to spider fail:" + err.Error()
				return
			}
			log.Debugf("[%s] node:%s, waiting response......\n", code_info, addr)
			var resp *message.ResponseMsg
			err, resp = arxlet.ParseResponseFromConn(spider_conn)
			log.Debugf("[%s] node:%s, get response \n", code_info, addr)
			if err != nil {
				log.Errorf("[%s] parseResponse Error:%s", code_info, err.Error())
				err_msg := "lost start info:" + err.Error()
				start_result[addr] = err_msg
				return
			}
			log.Infof("[%s] node[%s], start spider.ret:%d, msg:%s", code_info, addr, resp.Status, resp.Msg)
			start_ret_msg := ""
			if resp.Status != 0 {
				start_ret_msg = resp.Msg
			}
			start_result[addr] = start_ret_msg
		}(addr)
	}
	wg.Wait()
	return start_result
}

func (self *ArxScheduler) generateCrawlerTask(msg *message.SpiderStartMsg) *config.CrawlerTask {
	// cfg := &config.CrawlerTask{
	// 	TaskName:         "test",
	// 	TaskId:           11,
	// 	CrawlerTreadNum:  3,
	// 	RedisAddress:     "",
	// 	RedisPort:        8888,
	// 	RedisAccount:     "auth",
	// 	RedisPassword:    "",
	// 	RedisQueueName:   "test",
	// 	MasterListenPort: self.MasterListenPort,
	// }
	msg.Cfg.MasterListenPort = self.MasterListenPort
	return &msg.Cfg
}

func responseConn(ctx *arxlet.ConnContext, msg *message.ResponseMsg) error {
	code_info := "responseConn"
	send_bytes, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("[%s] marshal responseMsg fail:%s", code_info, err.Error())
		return err
	}
	_, err = ctx.From.Write(send_bytes)
	ctx.From.Close()
	return err
}
