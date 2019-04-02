/*
 * @Author: rayou
 * @Date: 2019-03-26 23:20:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-26 23:29:35
 */

package downloader

import (
	"github.com/AzuresYang/arx7/app/processor"
	"github.com/AzuresYang/arx7/app/spider/downloader/context"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
)

// 定义一个下载器接口， 接受一个请求， 解析器， 生成解析上下文
type Downloader interface {
	Download(*processor.Processor, *request.ArxRequest) *context.CommContext
}
