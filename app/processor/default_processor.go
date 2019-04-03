/*
 * @Author: rayou
 * @Date: 2019-03-30 09:34:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-02 23:27:59
 */
package processor

import (
	"github.com/AzuresYang/arx7/app/spider/downloader/context"
	log "github.com/sirupsen/logrus"
)

type DefaultProcessor struct {
	EnableCookie bool // 所有请求是否使用cookie记录
	// Namespace       func(self *Spider) string                                  // 命名空间，用于输出文件、路径的命名
	// SubNamespace    func(self *Spider, dataCell map[string]interface{}) string // 次级命名，用于输出文件、路径的命名，可依赖具体数据内容
	// 以下字段系统自动赋值
	Id int
}

func NewDefaultProcessor() DefaultProcessor {
	procer := DefaultProcessor{
		Id: 5,
	}
	return procer
}
func (self *DefaultProcessor) GetName() string {
	return "default"
}

func (self *DefaultProcessor) Process(ctx *context.CommContext) {
	log.Info("ready pro context:" + ctx.Request.Url)
}

func (self *DefaultProcessor) Free() {
	log.Info("default processor, use Freee Method")
	return
}

func (self *DefaultProcessor) GetOneProcessor() Processor {
	procer := DefaultProcessor{
		Id: 5,
	}
	return &procer
}
