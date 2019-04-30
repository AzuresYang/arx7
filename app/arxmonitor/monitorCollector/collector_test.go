package monitorCollector

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	// "time"

	"github.com/AzuresYang/arx7/app/arxmonitor"
	"github.com/AzuresYang/arx7/config"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

//数据库配置
const (
	userName = "root"
	password = "Mysql@#9420"
	ip       = "193.112.68.221"
	port     = "3306"
	dbName   = "monitor_info"
)

var DB *sql.DB

func connectDb() {
	// 构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	var err error
	// 打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, err = sql.Open("mysql", path)
	if err != nil {
		fmt.Printf("init db fail:%s\n", err.Error())
	}
	// 设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	// 设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func TestInsert(t *testing.T) {
	connectDb()
	stmt, err := DB.Prepare("INSERT INTO monitor_data (`svcid`, `metric`, `classfy`, `value`, `ip`, `time`) VALUES(?,?,?,?,?,?)")
	checkErr(err)
	res, err := stmt.Exec(5542, 1001, 1, 100, "127.0.0.1", 12345)
	if res != nil {
		fmt.Println(res.LastInsertId())
	}
	stmt.Close()
	checkErr(err)
	t.Log("donw")
}

func TestQueryUpdate(t *testing.T) {
	connectDb()
	stmt, err := DB.Prepare("SELECT id,value from  monitor_data where svcid=? and metric=? and classfy=? and ip =? and time=?")
	checkErr(err)
	rows, _ := stmt.Query(5542, 1, 0, "127", 12345678)
	defer rows.Close()
	i := 0
	var id int
	var value int
	for rows.Next() {
		if err := rows.Scan(&id, &value); err != nil {
			fmt.Printf("rows error:%s\n", err.Error())
		}
		fmt.Printf("query found[%d]%+v, value:%d\n", i, id, value)
		i++
	}
	update_sql := "update monitor_data set value=? where id=?"
	ustmt, err := DB.Prepare(update_sql)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = ustmt.Exec(value+1, id)
	if err != nil {
		fmt.Printf("update msg date fail:%s\n", err.Error())
	}
	stmt.Close()
	ustmt.Close()
	t.Log("donw")
}

func TestQueryMonitor(t *testing.T) {
	connectDb()
	if DB == nil {
		t.Error("DB is nil")
	}
	log.SetLevel(log.TraceLevel)
	// sql := "select time, value from monitor_data where svcid=? and metric=? and time in(?,?)"
}

func TestCollector(t *testing.T) {
	connectDb()
	if DB == nil {
		t.Error("DB is nil")
	}
	log.SetLevel(log.TraceLevel)
	server := New(5)
	conf := &config.MasterConfig{}
	err := config.ReadConfigFromFileJson("F:\\master.json", conf)
	if err != nil {
		fmt.Printf("read config fail:%s\n", err.Error())
		return
	}
	err = server.Start(&conf.MysqlConf)
	if err != nil {
		t.Error(err.Error())
		return
	}
	pkg := buildMoniMsgPkg()
	server.AddMonitorPkg(pkg)
	time.Sleep(5 * time.Second)
	server.Stop()
	//time.Sleep(3 * time.Second)
	t.Log("done")
}

func buildMoniMsgPkg() *arxmonitor.MonitorMsgPkg {
	bitPause := []int64{30, 30}
	// sleep_time := self.pause[0] + rand.Int63n(self.pause[1])
	var timeLen int64 = 60 * 60 * 2 // 2个小时
	var interSm int64 = 10
	toBeCharge := "2019-04-30 00:00:00"                             //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "2006-01-02 15:04:05"                             //转化所需模板
	loc, _ := time.LoadLocation("Local")                            //重要：获取时区
	theTime, _ := time.ParseInLocation(timeLayout, toBeCharge, loc) //使用模板在对应时区转化为time.time类型
	start := theTime.Unix()
	fmt.Println(start)
	end := start + timeLen
	pkg := arxmonitor.NewMonitorMsgPkg("127", 50)
	msg := arxmonitor.NewMonitorMsg(5542, 3001, arxmonitor.MONITORMSG_ADD)
	var lastValue int32 = 5
	valuePause := []int32{20, 600}
	for start <= end {
		// fmt.Printf("start:%d, end:%d\n", start, end)
		next_time := bitPause[0] + rand.Int63n(bitPause[1]) + start
		next_pause_value := valuePause[0] + rand.Int31n(valuePause[1])
		var cha int32 = next_pause_value - lastValue
		var i int64 = start
		for i < next_time {
			// fmt.Printf("i:%d, next_time:%d\n", i, next_time)
			msg.Time = i
			rand_value := float32(cha) * rand.Float32()
			cha -= int32(rand_value)
			lastValue += int32(rand_value)
			msg.Value = uint32(lastValue)
			pkg.AddMsg(msg)
			i += interSm
		}
		start = next_time
	}
	return pkg
}

func buildMsgPkg() *arxmonitor.MonitorMsgPkg {
	var now int64 = 12345678
	pkg := arxmonitor.NewMonitorMsgPkg("127", 50)
	msg := arxmonitor.NewMonitorMsg(5542, 1, arxmonitor.MONITORMSG_ADD)

	for i := 0; i < 2; i++ {
		pkg.AddMsg(msg)
	}
	msg = arxmonitor.NewMonitorMsg(5542, 3, arxmonitor.MONITORMSG_ADD)
	msg.Time = now
	for i := 0; i < 2; i++ {
		msg.Time = now
		pkg.AddMsg(msg)
	}

	msg = arxmonitor.NewMonitorMsg(5542, 5, arxmonitor.MONITORMSG_ADD)
	msg.Time = now
	for i := 0; i < 2; i++ {

		pkg.AddMsg(msg)
	}

	msg = arxmonitor.NewMonitorMsg(5542, 10, arxmonitor.MONITORMSG_SET)
	msg.Time = now
	for i := 0; i < 2; i++ {
		pkg.AddMsg(msg)
	}
	return pkg
}
func TestIntegrationMsgPkg(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	pkg := buildMoniMsgPkg()
	pkg = refactorMonitorPkg(pkg)
	fmt.Printf("Ip:%s\n", pkg.Ip)
	for i, msg := range pkg.Msgs {
		fmt.Printf("[%d]%+v\n", i, msg)
	}

}

func showDb() {
	rows, _ := DB.Query("select * from monitor_data;")
	for rows.Next() {
		var id int
		var svcid int
		var metric int
		var classfy int
		var value int
		var ip string
		var t string
		rows.Scan(&id, &svcid, &metric, &classfy, &value, &ip, &t)
		fmt.Printf("id[%d],svcid[%d],metric[%d],classfy[%d],value[%d],ip[%s],time[%s]\n", id, svcid, metric, classfy, value, ip, t)
	}

}

func deletedb() {
	_, err := DB.Exec("delete from monitor_data")
	if err == nil {
		fmt.Println("delete monitor_data succ")
	} else {
		fmt.Printf("delete db error:%s\n", err.Error())
	}
}
func TestSaver(t *testing.T) {
	connectDb()
	if DB == nil {
		t.Error("DB is nil")
	}
	log.SetLevel(log.TraceLevel)
	msgpkgs := make(chan *arxmonitor.MonitorMsgPkg, 5)
	save := &msgsaver{
		Db:     DB,
		ifStop: false,
	}
	// deletedb()
	showDb()
	pkg := buildMsgPkg()
	msgpkgs <- pkg
	log.Info("ready save msg")
	go save.doSaveMonitorMsg(msgpkgs)
	time.Sleep(2 * time.Second)
	showDb()
	t.Log("test saver succ")
}

func TestShowDb(t *testing.T) {
	connectDb()
	if DB == nil {
		t.Error("DB is nil")
	}
	log.SetLevel(log.TraceLevel)
	showDb()
}
