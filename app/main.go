/*
 * @Author: rayou
 * @Date: 2019-04-22 19:11:44
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-22 23:26:20
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
					Value: "file",
					Usage: "set info out put.[file, screen]",
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
					Value: "file",
					Usage: "set info out put.[file, screen]",
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

	// 用linux的启动方式，这里就不做设置日志的保存文件了
	masterSvr := arxmaster.NewArxMaster()
	port := ctx.String("port")
	err := masterSvr.Init(port)
	if err != nil {
		fmt.Printf("Init master fail:%s\n", err.Error())
		fmt.Println("----start arx7 master fail!-----")
		return
	}
	fmt.Println("----start arx7 master succ----")
	masterSvr.Run()
}

func startAsSpider(ctx *cli.Context) {
	Init()
	// 设置日志级别
	log_level := ctx.String("log-level")
	setLogLevel(log_level)

	// 用linux的启动方式，这里就不做设置日志的保存文件了
	spider := spider.NewSpider()
	port := ctx.String("port")
	err := spider.Init(port)
	if err != nil {
		fmt.Printf("Init spider fail:%s\n", err.Error())
		fmt.Println("----start arx7 master fail!-----")
		return
	}
	fmt.Println("----start arx7 spider succ----")
	spider.Run()
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
