/*
 * @Author: rayou
 * @Date: 2019-04-11 19:00:44
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-11 22:52:23
 */

package monitorCollector

import (
	"database/sql"
	"strings"
	"time"

	"github.com/AzuresYang/arx7/app/arxmonitor"
	"github.com/AzuresYang/arx7/db"
	log "github.com/sirupsen/logrus"
)

const (
	DEFAULT_MONITOR_TIME_INTERVAL = 60 // 监控时间精度，单位 s
)

type monitorServer struct {
	db      *sql.DB
	ifStop  bool
	msgpkgs chan *arxmonitor.MonitorMsgPkg
}

func (self *monitorServer) Start(db_cfg *db.MysqlConfig) error {
	path := strings.Join([]string{db_cfg.UserName, ":", db_cfg.Password, "@tcp(", db_cfg.Ip, ":",
		string(db_cfg.Port), ")/", db_cfg.DbName, "?charset=", db_cfg.Charset}, "")
	var err error
	self.db, err = sql.Open("mysql", path)
	if err != nil {
		log.Errorf("[MonitorServer.Start] Open mysql error:%s", err.Error())
		return err
	}
	// 设置链接最大空闲时间
	self.db.SetConnMaxLifetime(100 * time.Second)
	// 设置上数据库最大闲置连接数
	self.db.SetMaxIdleConns(db_cfg.MaxIdleConns)
	go self.Run()
	return nil
}
func (self *monitorServer) Run() {
	stmt, err := self.db.Prepare("INSERT INTO monitor_data (`svcid`, `metric`, `classfy`, `value`, `ip`, `time`)" +
		" VALUES(?,?,?,?,?,FROM_UNIXTIME(?))")
	if err != nil {
		log.Errorf("[MonitorServer] Prepare Insert statment fail:%s", err.Error())
		return
	}
	res, err := stmt.Exec(5542, 1001, 1, 100, "127.0.0.1", time.Now().Unix())
	stmt.Close()
	if err != nil {
		log.Errorf("[MonitorServer]insert monitor data err:%s", err.Error())
	}
	if res != nil {
		log.Trace("insert monitor data succ")
	}
	self.Stop()
}

func (self *monitorServer) Stop() {
	self.ifStop = true
}

func (self *monitorServer) IfStop() bool {
	return self.ifStop
}

// 对监控数据包进行重构，所有时间相同的
// 统一成分钟的格式进行存储
func refactorMonitorPkg(pkg *arxmonitor.MonitorMsgPkg) *arxmonitor.MonitorMsgPkg {
	ret_pkg := pkg
	for i, msg := range pkg.Msgs {
		pkg.Msgs[i].Time = msg.Time - msg.Time%DEFAULT_MONITOR_TIME_INTERVAL
	}
	ret_pkg = integrationMonitorMsgPkg(pkg)
	return ret_pkg
}

// 整合监控数据包
func integrationMonitorMsgPkg(pkg *arxmonitor.MonitorMsgPkg) *arxmonitor.MonitorMsgPkg {
	ret_pkg := &arxmonitor.MonitorMsgPkg{
		Ip:   pkg.Ip,
		Msgs: make([]arxmonitor.MonitorMsg, 0, len(pkg.Msgs)),
	}
	// 相同时间的SET 类型的监控数据，只取一个就好了
	set := make(map[int]bool)
	for i, msg := range pkg.Msgs {
		if set[i] {
			// 之前已经处理过这个，不再处理
			continue
		}
		// 找到相同的时间的监控数据，把下标记下来，之后不要重复处理
		idxs := findSameTimeMonitorMsg(&msg, pkg.Msgs)
		log.Tracef("[MonitorServer.Integration]same msg:%+v", idxs)
		for _, idx := range idxs {
			set[idx] = true
		}
		switch msg.MsgType {
		case arxmonitor.MONITORMSG_ADD:
			new_msg := integrationAddMsg(idxs, pkg.Msgs)
			ret_pkg.AddMsg(new_msg)
		case arxmonitor.MONITORMSG_SET:
			new_msg := integrationSetMsg(idxs, pkg.Msgs)
			ret_pkg.AddMsg(new_msg)
		default:
			log.Errorf("[MonitorServer] found unknown monitormsg type:%#v", msg.MsgType)
		}
	}
	return ret_pkg
}

func integrationAddMsg(idxs []int, msgs []arxmonitor.MonitorMsg) *arxmonitor.MonitorMsg {
	new_msg := msgs[idxs[0]]
	for i, idx := range idxs {
		if i == 0 {
			continue
		}
		new_msg.Value += msgs[idx].Value
	}
	return &new_msg
}

func integrationSetMsg(idxs []int, msgs []arxmonitor.MonitorMsg) *arxmonitor.MonitorMsg {
	return &msgs[idxs[0]]
}

func findSameTimeMonitorMsg(targetmsg *arxmonitor.MonitorMsg, msgs []arxmonitor.MonitorMsg) []int {
	var idx []int
	for i, m := range msgs {
		if arxmonitor.IsEqualMsgButValue(targetmsg, &m) {
			idx = append(idx, i)
		}
	}
	return idx
}
