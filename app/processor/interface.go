/*
 * @Author: rayou
 * @Date: 2019-03-25 22:50:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-02 23:14:48
 */
package processor

import (
	"github.com/AzuresYang/arx7/app/pipeline"
	"github.com/AzuresYang/arx7/app/spider/downloader/context"
)

// 处理下载数据， 对需要解析不同页面，可以自己实现不同的处理器
type Processor interface {
	GetName() string // 处理器的代称
	// 获取一个处理器
	// 这里这么设计的原因是：对于不同的解析器, 下载好页面之后，需要找到处理这个页面的处理器
	// 一种是直接使用Process()方法， 一种是生成一个处理器后，由处理器进行解析。
	// 目前采用第二种方法， 返回的是值还是指针，这个有实现者决定
	GetOneProcessor() Processor
	Free()
	// 返回错误码和处理消息
	Process(*context.CommContext) *pipeline.CollectData
}
