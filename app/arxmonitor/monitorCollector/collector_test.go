package monitorCollector

import (
	"database/sql"
	"fmt"
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
	password = "mysql5722"
	ip       = "127.0.0.1"
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

func TestCollector(t *testing.T) {
	connectDb()
	if DB == nil {
		t.Error("DB is nil")
	}
	log.SetLevel(log.TraceLevel)
	server := New(5)
	cfg := config.NewMysqlConfig()
	err := server.Start(cfg)
	if err != nil {
		t.Error(err.Error())
	}
	pkg := buildMsgPkg()
	server.AddMonitorPkg(pkg)
	time.Sleep(3 * time.Second)
	server.Stop()
	time.Sleep(3 * time.Second)
	t.Log("done")
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
	pkg := buildMsgPkg()
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
