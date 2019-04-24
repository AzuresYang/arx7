/*
 * @Author: rayou
 * @Date: 2019-04-03 20:37:41
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 23:06:42
 */
package config

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type CrawlerTask struct {
	TaskName                    string // 任务名
	TaskId                      uint32 // 任务ID
	CrawlerTreadNum             uint32 // 爬虫线程数量
	RedisAddr                   string
	RedisPassword               string        // redis密码
	MaxGetRequestNullTimeSecond time.Duration // 长时间内没有新链接时，停止工作crawler的设置， 单位：秒， 为0时表示一直工作
	MasterListenPort            string        // SpiderMaster监听的端口
	FastDfsAddr                 string        // 分布式系统地址
}

type SpiderStartConfig struct {
	TaskConf   CrawlerTask // 任务配置
	ProcerName string      // 处理器名称
	Urls       []string    // 原始URL
}

type MasterConfig struct {
	MysqlConf  MysqlConfig
	ListenPort string
}

func WriteConfigToFileJson(dir string, fileName string, conf interface{}) error {
	json_byte, err := json.MarshalIndent(conf, "", "\n")
	if err != nil {
		return errors.New("Marshl CrawlerTask fail")
	}
	file_path := dir + "/" + fileName
	err = os.MkdirAll(filepath.Clean(dir), 0777)
	if err != nil {
		return err
	}
	var f *os.File
	// 退出时关闭文件
	defer func() {
		if f != nil {
			f.Close()
		}
	}()

	// 文件不存在时创建一个
	if _, err := os.Stat(file_path); os.IsNotExist(err) {
		f, err = os.Create(file_path)
	} else {
		f, err = os.OpenFile(file_path, os.O_WRONLY, 0777)
	}
	content := string(json_byte[:])
	if err != nil {
		return err
	}
	_, w_err := io.WriteString(f, content)
	return w_err
}

func ReadConfigFromFileJson(filePth string, conf interface{}) error {
	f, err := os.Open(filePth)
	if err != nil {
		return err
	}
	var data []byte
	data, err = ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, conf)
}
