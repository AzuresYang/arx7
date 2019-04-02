/*
 * @Author: rayou
 * @Date: 2019-03-25 23:09:36
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-26 23:35:12
 */
package context

import (
	"net/http"

	"github.com/AzuresYang/arx7/app/spider/downloader/request"
)

// Downloader下载内容构成的上下文
type CommContext struct {
	Processor *Processor          // 使用的解析器
	Request   *request.ArxRequest // 原始请求
	Response  *http.Response      // http的响应流
	ErrMsg    string              // 错误描述， 可以改
	Stat      int                 // 状态码，0 是成功
}

// 后续可以考虑使用syn里面的池子方法来构建context上下文的池子
