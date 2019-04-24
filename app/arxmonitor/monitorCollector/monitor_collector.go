/*
 * @Author: rayou
 * @Date: 2019-04-11 19:00:44
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 23:36:45
 */

package monitorCollector

import (
	"database/sql"
	"strings"
	"time"

	"github.com/AzuresYang/arx7/app/arxmonitor"
	"github.com/AzuresYang/arx7/config"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

const (
	DEFAULT_MONITOR_TIME_INTERVAL = 60 // 监控时间精度，单位 s
)

// 用于保存数据到数据库中，支持开启多个
type msgsaver struct {
	Id     int
	ifStop bool
	Db     *sql.DB
}

type MonitorCollector struct {
	db      *sql.DB
	ifStop  bool
	msgpkgs chan *arxmonitor.MonitorMsgPkg
	savers  []*msgsaver
	dbCfg   *config.MysqlConfig
}

func New(maxPkgs int) *MonitorCollector {
	c := &MonitorCollector{
		ifStop:  true,
		msgpkgs: make(chan *arxmonitor.MonitorMsgPkg, maxPkgs),
	}
	return c
}
func (self *MonitorCollector) Start(db_cfg *config.MysqlConfig) error {
	path := strings.Join([]string{db_cfg.UserName, ":", db_cfg.Password, "@tcp(", db_cfg.Ip, ":",
		string(db_cfg.Port), ")/", db_cfg.DbName, "?charset=", db_cfg.Charset}, "")
	var err error
	self.db, err = sql.Open("mysql", path)
	if err != nil {
		log.Errorf("[MonitorCollector.Start] Open mysql error:%s", err.Error())
		return err
	}
	log.Debugf("[MonitorCollector.Start]config :%+v", db_cfg)
	// 设置链接最大空闲时间
	self.db.SetConnMaxLifetime(100 * time.Second)
	// 设置上数据库最大闲置连接数
	self.db.SetMaxIdleConns(db_cfg.MaxIdleConns)
	for i := 0; i < db_cfg.TaskNum; i++ {
		saver := &msgsaver{
			Db:     self.db,
			Id:     i,
			ifStop: false,
		}
		self.savers = append(self.savers, saver)
	}
	log.Info("[MonitorCollector.Start] start...")
	self.ifStop = false
	go self.Run()
	return nil
}

func (self *MonitorCollector) Stop() {
	self.ifStop = true
	log.Info("[MonitorCollector.Stop] stop...")
	for i, _ := range self.savers {
		self.savers[i].Stop()
	}
}

func (self *MonitorCollector) IfStop() bool {
	return self.ifStop
}

func (self *MonitorCollector) Run() {
	// 启动每一个数据库交互线程
	self.ifStop = false
	log.Infof("[MonitorCollecotr.Run]saver thread is:%d", len(self.savers))
	for i, _ := range self.savers {
		self.savers[i].ifStop = false
		go self.savers[i].doSaveMonitorMsg(self.msgpkgs)
	}
}

func (self *MonitorCollector) AddMonitorPkg(pkg *arxmonitor.MonitorMsgPkg) {
	code_info := "MonitorCollector.AddMonitorPkg"
	log.Tracef("[%s]recv monitor pkg:%+v", code_info, pkg)
	if self.ifStop {
		return
	}
	self.msgpkgs <- pkg
}

func (self *msgsaver) Stop() {
	self.ifStop = true
	log.Debugf("[%d] saver stop:%v", self.Id, self)
}

// 保存监控数据
func (self *msgsaver) doSaveMonitorMsg(c chan *arxmonitor.MonitorMsgPkg) {
	code_info := "MonitorCollector.Saver.doSaveMonitorMsg"
	insert_sql := "INSERT INTO monitor_data (`svcid`, `metric`, `classfy`, `value`, `ip`, `time`)" +
		" VALUES(?,?,?,?,?,?)"
	query_sql := "SELECT id,value from  monitor_data where " +
		"svcid=? and metric=? and classfy=? and ip =? and time=?"
	update_sql := "update monitor_data set value=? where id=?"
	if self.Db == nil {
		log.Errorf("[%s]db is nil.", code_info)
		return
	}

	// defer func() {
	// 	insert_stmt.Close()
	// 	query_stmt.Close()
	// 	update_stmt.Close()
	// }()
	// if ierr != nil || qerr != nil || uerr != nil {
	// 	log.Errorf("[%s] sql prepare fail:insert:%s|query:%s|update:%s", code_info, ierr.Error(), qerr.Error(), uerr.Error())
	// 	return
	// }
	log.Debugf("[%s]waiting monitor msg.", code_info)
	for !self.ifStop {
		select {
		case pkg := <-c:
			log.Debugf("[%s] get monitor pkg.ip:%s, len:%d", code_info, pkg.Ip, len(pkg.Msgs))
			insert_stmt, ierr := self.Db.Prepare(insert_sql)
			query_stmt, qerr := self.Db.Prepare(query_sql)
			update_stmt, uerr := self.Db.Prepare(update_sql)
			if ierr != nil || qerr != nil || uerr != nil {
				log.Errorf("[%s] sql prepare fail:insert:%s|query:%s|update:%s", code_info, ierr.Error(), qerr.Error(), uerr.Error())
				continue
			}
			err := self.addMonitorMsgToDb(pkg, query_stmt, insert_stmt, update_stmt)
			if err != nil {
				log.Errorf("[%s]add monitor msg fail:%s", code_info, err.Error())
			}
			insert_stmt.Close()
			query_stmt.Close()
			update_stmt.Close()
		default:
			if !self.ifStop {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
	log.Debugf("[%s][%d] stop save monitormsg", code_info, self.Id)
}

func (self *msgsaver) addMonitorMsgToDb(pkg *arxmonitor.MonitorMsgPkg, qstmt *sql.Stmt, istmt *sql.Stmt, ustmt *sql.Stmt) error {
	code_info := "MonitorCollector.Saver.addMonitorMsgToDb"
	// 先查一下是否有这个监控数值，有的话累加， 没有的话，设置
	for i, _ := range pkg.Msgs {
		msg := &pkg.Msgs[i]
		log.Tracef("[%s]process msg:%+v", code_info, msg)
		var id int64 = -1
		var value int
		// 先查看是否存在原先监控项
		rows, err := qstmt.Query(msg.SvcId, msg.Metric, msg.Classfy, pkg.Ip, msg.Time)
		if err != nil {
			log.Errorf("[%s]query monitor msg fail:%s", code_info, err.Error())
			continue
		}
		for rows.Next() {
			rows.Scan(&id, &value)
			log.Tracef("[%s]query found metric:%d, id:%d, value:%d", code_info, msg.Metric, id, value)
			break
		}
		// 如果不存在，添加一个
		if id < 0 {
			ret, err := istmt.Exec(msg.SvcId, msg.Metric, msg.Classfy, msg.Value, pkg.Ip, msg.Time)
			if err != nil {
				log.Errorf("[%s] insert monitor msg fail:%s", code_info, err.Error())
			} else if ret != nil {
				last_id, _ := ret.LastInsertId()
				log.Tracef("[%s]insert monitor msg succ.id:%d", code_info, last_id)
			}
		} else {
			// 如果已经存在监控项， 累加监控值，或者设置监控值
			if msg.MsgType == arxmonitor.MONITORMSG_ADD {
				log.Tracef("[%s]should add monitor value[q:%d|m:%d]", code_info, value, msg.Value)
				value += int(msg.Value)

			}
			_, err := ustmt.Exec(value, id)
			if err != nil {
				log.Errorf("[%s] update monitor msg fail:%s", code_info, err.Error())
			} else {
				log.Tracef("[%s] update monitor.id:%d,metric:%d, value:%d", code_info, id, msg.Metric, value)
			}
		}
	}
	log.Tracef("[%s] add monitor pkg done.ip:%s, len:%d", code_info, pkg.Ip, len(pkg.Msgs))
	return nil
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

// 整合累加类型的监控消息
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

// 整合设置类型的监控消息
func integrationSetMsg(idxs []int, msgs []arxmonitor.MonitorMsg) *arxmonitor.MonitorMsg {
	return &msgs[idxs[0]]
}

// 找到相同时间，相同类型的监控数据（除了value不同）
func findSameTimeMonitorMsg(targetmsg *arxmonitor.MonitorMsg, msgs []arxmonitor.MonitorMsg) []int {
	var idx []int
	for i, m := range msgs {
		if arxmonitor.IsEqualMsgButValue(targetmsg, &m) {
			idx = append(idx, i)
		}
	}
	return idx
}
