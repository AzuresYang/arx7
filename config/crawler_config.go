package config

import (
	"strings"

	"github.com/henrylee2cn/pholcus/logs/logs"
	"github.com/henrylee2cn/pholcus/runtime/cache"
)

// 软件信息。
const (
	VERSION   string = "v1.0.0"                                      // 软件版本号
	AUTHOR    string = "AzuresYang"                                 // 软件作者
	NAME      string = "ARX 网络爬虫"                              // 软件名
	FULL_NAME string = NAME + "_" + VERSION + " （by " + AUTHOR + "）" // 软件全称
	TAG       string = "arx"                                     // 软件标识符
)

// 默认配置。
const (
	WORK_ROOT      string = TAG + "_pkg"                    // 运行时的目录名称
	CONFIG         string = WORK_ROOT + "/config.json"       // 配置文件路径
	CACHE_DIR      string = WORK_ROOT + "/cache"            // 缓存文件目录
	LOG            string = WORK_ROOT + "/logs/arx_crawler.log" // 日志文件路径
	// LOG_ASYNC      bool   = true                            // 是否异步输出日志
	PHANTOMJS_TEMP string = CACHE_DIR                       // Surfer-Phantom下载器：js文件临时目录
)

// 来自配置文件的配置项。
var (
	CRAWLS_CAP int = setting.DefaultInt("crawlcap", crawlcap) // 蜘蛛池最大容量
	// DATA_CHAN_CAP            int    = setting.DefaultInt("datachancap", datachancap)                               // 收集器容量
	PHANTOMJS                string = setting.String("phantomjs")                                          // Surfer-Phantom下载器：phantomjs程序路径
	PROXY                    string = setting.String("proxylib")                                           // 代理IP文件路径
	SPIDER_DIR               string = setting.String("spiderdir")                                          // 动态规则目录
	FILE_DIR                 string = setting.String("fileoutdir")                                         // 文件（图片、HTML等）结果的输出目录
	TEXT_DIR                 string = setting.String("textoutdir")                                         // excel或csv输出方式下，文本结果的输出目录
	DB_NAME                  string = setting.String("dbname")                                             // 数据库名称
	MGO_CONN_STR             string = setting.String("mgo::connstring")                                    // mongodb连接字符串
	MGO_CONN_CAP             int    = setting.DefaultInt("mgo::conncap", mgoconncap)                       // mongodb连接池容量
	MGO_CONN_GC_SECOND       int64  = setting.DefaultInt64("mgo::conngcsecond", mgoconngcsecond)           // mongodb连接池GC时间，单位秒
	MYSQL_CONN_STR           string = setting.String("mysql::connstring")                                  // mysql连接字符串
	MYSQL_CONN_CAP           int    = setting.DefaultInt("mysql::conncap", mysqlconncap)                   // mysql连接池容量
	MYSQL_MAX_ALLOWED_PACKET int    = setting.DefaultInt("mysql::maxallowedpacket", mysqlmaxallowedpacket) // mysql通信缓冲区的最大长度

	KAFKA_BORKERS string = setting.DefaultString("kafka::brokers", kafkabrokers) //kafka brokers

	LOG_CAP            int64 = setting.DefaultInt64("log::cap", logcap)          // 日志缓存的容量
	LOG_LEVEL          int   = logLevel(setting.String("log::level"))            // 全局日志打印级别（亦是日志文件输出级别）
	LOG_CONSOLE_LEVEL  int   = logLevel(setting.String("log::consolelevel"))     // 日志在控制台的显示级别
	LOG_FEEDBACK_LEVEL int   = logLevel(setting.String("log::feedbacklevel"))    // 客户端反馈至服务端的日志级别
	LOG_LINEINFO       bool  = setting.DefaultBool("log::lineinfo", loglineinfo) // 日志是否打印行信息                                  // 客户端反馈至服务端的日志级别
	LOG_SAVE           bool  = setting.DefaultBool("log::save", logsave)         // 是否保存所有日志到本地文件
)


func logLevel(l string) int {
	switch strings.ToLower(l) {
	case "app":
		return logs.LevelApp
	case "emergency":
		return logs.LevelEmergency
	case "alert":
		return logs.LevelAlert
	case "critical":
		return logs.LevelCritical
	case "error":
		return logs.LevelError
	case "warning":
		return logs.LevelWarning
	case "notice":
		return logs.LevelNotice
	case "informational":
		return logs.LevelInformational
	case "info":
		return logs.LevelInformational
	case "debug":
		return logs.LevelDebug
	}
	return -10
}
