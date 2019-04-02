/*
 * @Author: rayou
 * @Date: 2019-03-26 21:46:02
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-27 19:05:20
 */
package request

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type HttpMethod string

var (
	GET  HttpMethod = "GET"
	POST HttpMethod = "POST"
)

const (
	DEFALUT_CONNECT_TIMEOUT  = 10 * time.Second // 默认链接超时
	DEFALUT_TRY_TIMES        = 3                // 重试次数
	DEFALUT_DOWNLOAD_TIMEOUT = 2 * time.Minute  // 默认下载超时
)

type ArxRequest struct {
	SpiderName      string        // 要使用的spider,适配于解析规则
	Url             string        // url
	Header          http.Header   // http头
	Method          string        // 请求方法， 使用大写
	EnableCookie    bool          // 是否使用cookie
	PostData        string        // post数据
	Priority        int           // 该请求的优先级, 数字越大，优先级越高
	TryTimes        int           // 重连次数
	ConnTimeout     time.Duration // 链接超时
	DownloadTimeout time.Duration // 下载超时
	IDDownloader    int           // 下载器id
}

func (self *ArxRequest) Prepare() error {
	URL, err := url.Parse(self.Url)
	if err != nil {
		return err
	}
	self.Url = URL.String()
	if self.Header == nil {
		self.Header = make(http.Header)
	}
	if self.Method == "" {
		self.Method = "GET"
	} else {
		self.Method = strings.ToUpper(self.Method)
	}

	if self.Priority < 0 {
		self.Priority = 0
	}
	if self.TryTimes < 0 {
		self.TryTimes = DEFALUT_TRY_TIMES
	}

	if self.ConnTimeout < 0 {
		self.ConnTimeout = DEFALUT_CONNECT_TIMEOUT
	}

	if self.DownloadTimeout < 0 {
		self.DownloadTimeout = DEFALUT_DOWNLOAD_TIMEOUT
	}
	return nil
}

func NewArxRequest(url string) ArxRequest {
	request := ArxRequest{}
	request.Url = url
	request.Method = "GET"
	request.Header = make(http.Header)
	request.EnableCookie = false
	request.TryTimes = DEFALUT_TRY_TIMES
	request.ConnTimeout = DEFALUT_CONNECT_TIMEOUT
	request.DownloadTimeout = DEFALUT_DOWNLOAD_TIMEOUT
	return request
}

func (self *ArxRequest) Serialize() string {
	json_byte, _ := json.Marshal(self)
	return string(json_byte[:])
}

// 反序列化
func UnSerialize(s string) (*ArxRequest, error) {
	req := new(ArxRequest)
	return req, json.Unmarshal([]byte(s), req)
}
