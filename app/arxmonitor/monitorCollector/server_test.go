package monitorCollector

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/AzuresYang/arx7/app/arxmonitor"
	"github.com/AzuresYang/arx7/db"
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
func testInsert() {
	// 构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{userName, ":", password, "@tcp(", ip, ":", port, ")/", dbName, "?charset=utf8"}, "")
	// 打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	// 设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	// 设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	stmt, err := DB.Prepare("INSERT INTO monitor_data (`svcid`, `metric`, `classfy`, `value`, `ip`, `time`) VALUES(?,?,?,?,?,FROM_UNIXTIME(?))")
	checkErr(err)
	res, err := stmt.Exec(5542, 1001, 1, 100, "127.0.0.1", time.Now().Unix())
	if res != nil {
		fmt.Println(res.LastInsertId())
	}
	stmt.Close()
	checkErr(err)
}

func TestServer(t *testing.T) {
	server := &monitorServer{}
	cfg := db.NewMysqlConfig()
	err := server.Start(cfg)
	if err != nil {
		t.Error(err.Error())
	}
	for {
		if server.IfStop() {
			break
		}
	}
	t.Log("done")
}

func buildMsgPkg() *arxmonitor.MonitorMsgPkg {
	pkg := arxmonitor.NewMonitorMsgPkg("127", 50)
	msg := arxmonitor.NewMonitorMsg(5542, 1, arxmonitor.MONITORMSG_ADD)
	for i := 0; i < 1; i++ {
		pkg.AddMsg(msg)
	}
	msg = arxmonitor.NewMonitorMsg(5542, 3, arxmonitor.MONITORMSG_ADD)
	for i := 0; i < 3; i++ {
		pkg.AddMsg(msg)
	}

	msg = arxmonitor.NewMonitorMsg(5542, 5, arxmonitor.MONITORMSG_ADD)
	for i := 0; i < 5; i++ {
		pkg.AddMsg(msg)
	}

	msg = arxmonitor.NewMonitorMsg(5542, 10, arxmonitor.MONITORMSG_SET)
	for i := 0; i < 7; i++ {
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
