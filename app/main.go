/*
 * @Author: rayou
 * @Date: 2019-04-22 19:11:44
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 23:42:34
 */
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AzuresYang/arx7/app/arxmaster"
	"github.com/AzuresYang/arx7/app/processor"
	"github.com/AzuresYang/arx7/app/spider"
	"github.com/AzuresYang/arx7/config"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	log.SetLevel(log.DebugLevel)
	app := cli.NewApp()
	app.Name = "arxmaster"
	app.Usage = "arx7 command-line app by azureyang"
	app.Author = "AzureYang"
	app.Email = "AzureYang@xxx.com"
	app.Commands = []cli.Command{
		// 开始
		{
			Name:   "master",
			Usage:  "start as spider master",
			Action: startAsMaster,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Value: "31001",
					Usage: "listen port port",
				},
				cli.StringFlag{
					Name:  "log-level",
					Value: "debug",
					Usage: "set log level.[trace, debug, info, error, fatal]",
				},
				cli.StringFlag{
					Name:  "output",
					Value: "master.log",
					Usage: "set info output file.nil is screen",
				},
				cli.StringFlag{
					Name:  "config",
					Value: "",
					Usage: "master config file",
				},
			},
		},
		// 作为spider启动
		{
			Name:   "spider",
			Usage:  "start as spider",
			Action: startAsSpider,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "port",
					Value: "31002",
					Usage: "listen port port",
				},
				cli.StringFlag{
					Name:  "log-level",
					Value: "debug",
					Usage: "set log level.[trace, debug, info, error, fatal]",
				},
				cli.StringFlag{
					Name:  "output",
					Value: "spider.log",
					Usage: "set info output file.nil is screen",
				},
			},
		},
		// 作为spider启动
		{
			Name:  "config",
			Usage: "generate config",
			Subcommands: []cli.Command{
				{
					Name:   "master",
					Usage:  "generate master default config",
					Action: createMasterConfig,
				},
			},
		},
	}
	app.Run(os.Args)
}

func startAsMaster(ctx *cli.Context) {
	Init()
	// 设置日志级别
	log_level := ctx.String("log-level")
	setLogLevel(log_level)
	file := ctx.String("output")
	if file != "" {
		setOutput(file)
	}
	// 用linux的启动方式，这里就不做设置日志的保存文件了
	masterSvr := arxmaster.NewArxMaster()
	port := ctx.String("port")
	// 优先采用配置里的port
	config_file := ctx.String("config")
	conf := &config.MasterConfig{}
	if config_file != "" {
		err := config.ReadConfigFromFileJson(config_file, conf)
		if err != nil {
			fmt.Printf("[error]read config[%s] fail:%s\n", config_file, err.Error())
		} else {
			port = conf.ListenPort
		}
	}
	// 初始化
	err := masterSvr.Init(port)
	if err != nil {
		fmt.Printf("[error]init master fail:%s\n", err.Error())
		fmt.Println("----start arx7 master fail!-----")
		return
	}
	// 启动监控收集
	if config_file != "" {
		err = masterSvr.StartMonitorCollector(&conf.MysqlConf)
		if err != nil {
			fmt.Printf("[error]start monitor collector fail:%s\n", err.Error())
			return
		}
	}
	fmt.Println("----start arx7 master succ----")
	masterSvr.Run()
}

func startAsSpider(ctx *cli.Context) {
	Init()
	// 设置日志级别
	log_level := ctx.String("log-level")
	setLogLevel(log_level)
	file := ctx.String("output")
	if file != "" {
		setOutput(file)
	}
	// 用linux的启动方式，这里就不做设置日志的保存文件了
	spider := spider.NewSpider()
	port := ctx.String("port")
	err := spider.Init(port)
	if err != nil {
		fmt.Printf("[error]Init spider fail:%s\n", err.Error())
		fmt.Println("----start arx7 master fail!-----")
		return
	}
	fmt.Println("----start arx7 spider succ----")
	spider.Run()
}

// 创建默认配置
func createMasterConfig(ctx *cli.Context) {
	cfg := &config.MasterConfig{}
	cfg.MysqlConf = *config.NewMysqlConfig()
	dir := "./"
	file_name := "master.json"
	err := config.WriteConfigToFileJson(dir, file_name, cfg)
	if err != nil {
		fmt.Printf("[error]create config fail:%s\n", err.Error())
		return
	}
	file_path := dir + "/" + file_name
	fmt.Printf("create config file succ:%s\n", file_path)
}

// 运行前初始化
func Init() {
	os.MkdirAll(filepath.Clean("./arx7/config"), 0777)
	os.MkdirAll(filepath.Clean("./arx7/log"), 0777)
	os.MkdirAll(filepath.Clean("./arx7/data"), 0777)
}

func InitProcer() {
	procer := processor.NewDefaultProcessor()
	processor.Manager.Register(&procer)
}
func setOutput(file string) {
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		fmt.Printf("open file fail%s.\n", file)
		return
	}
	log.SetOutput(f)
}
func setLogLevel(level_str string) {
	level := log.DebugLevel
	level_str = strings.ToLower(level_str)
	switch level_str {
	case "trace":
		level = log.TraceLevel
	case "debug":
		level = log.DebugLevel
	case "info":
		level = log.InfoLevel
	case "error":
		level = log.ErrorLevel
	case "fatal":
		level = log.FatalLevel
	}
	log.SetLevel(level)
}
