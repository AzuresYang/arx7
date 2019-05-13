package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// 软件信息。
const (
	VERSION   string = "v1.0.0"                                      // 软件版本号
	AUTHOR    string = "AzuresYang"                                  // 软件作者
	NAME      string = "ARX 网络爬虫"                                    // 软件名
	FULL_NAME string = NAME + "_" + VERSION + " （by " + AUTHOR + "）" // 软件全称
	TAG       string = "arx"                                         // 软件标识符
)

// 默认配置。
const (
	WORK_ROOT      string = "./" + TAG + "_pkg" // 运行时的目录名称
	CONFIG_DIR     string = WORK_ROOT + string(filepath.Separator) + "cfg" + string(filepath.Separator)
	CONFIG         string = "config.json"                       // 配置文件路径
	CRAWLER_CONFIG string = "crawler_config.json"               // 配置文件路径
	CACHE_DIR      string = WORK_ROOT + "/cache"                // 缓存文件目录
	LOG            string = WORK_ROOT + "/logs/arx_crawler.log" // 日志文件路径
	// LOG_ASYNC      bool   = true                            // 是否异步输出日志
	PHANTOMJS_TEMP               string        = CACHE_DIR // Surfer-Phantom下载器：js文件临时目录
	DEFAULT_REQ_GET_TIMEOUT      time.Duration = 2 * time.Second
	DEFAULT_REQ_IS_NULL_WAITTIME time.Duration = 500 * time.Millisecond
	DEFAULT_REQ_MAX_NULL_TIME    uint32        = 100 // 获取req为空的时间最长时间，超过这个时间，则爬虫停止爬取
)

type CrawlerConfig struct {
	ConfigDir      string // 配置文件所在路径
	ConfigFileName string // 配置文件名
	MasterAddr     string // master地址
	TaskConf       CrawlerTask

	RequestGetTimeOut uint32
}

func BuildDefaultCrawlerConfig() *CrawlerConfig {
	conf := &CrawlerConfig{}
	conf.ConfigDir = CONFIG_DIR
	conf.ConfigFileName = CRAWLER_CONFIG
	conf.RequestGetTimeOut = uint32(DEFAULT_REQ_GET_TIMEOUT.Seconds())
	// 爬取任务默认配置
	conf.TaskConf = CrawlerTask{
		TaskName:                    "-1",
		TaskId:                      0,
		CrawlerTreadNum:             1,
		MaxGetRequestNullTimeSecond: DEFAULT_REQ_MAX_NULL_TIME,
	}
	return conf
}

func (self *CrawlerConfig) Serialize() string {
	json_byte, _ := json.Marshal(self)
	return string(json_byte[:])
}

func WriteToFile(conf *CrawlerConfig) error {
	json_byte, err := json.Marshal(conf)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Clean(conf.ConfigDir), 0777)
	if err != nil {
		return err
	}
	file_name := conf.ConfigDir + conf.ConfigFileName
	// file_name := conf.ConfigFileName
	fmt.Printf("ready to open file:%s\n", file_name)
	var f *os.File
	// 退出时关闭文件
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	// 文件不存在时创建一个
	if _, err := os.Stat(file_name); os.IsNotExist(err) {
		f, err = os.Create(file_name)
	} else {
		f, err = os.OpenFile(file_name, os.O_WRONLY, 0777)
	}
	content := string(json_byte[:])
	if err != nil {
		return err
	}
	_, w_err := io.WriteString(f, content)
	return w_err
}

func GetFromFile(file string) *CrawlerConfig {
	return nil
}
