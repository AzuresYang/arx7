/*
 * @Author: rayou
 * @Date: 2019-03-26 21:46:02
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-10 01:33:43
 */
package request

import (
	"encoding/json"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
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

type Temp map[string]interface{}
type ArxRequest struct {
	ProcerName      string        // 要使用的Procer,适配于解析规则
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
	Temp            Temp          // 临时数据，给处理器processor用的
	TempStrMap      map[string]string
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

func NewArxRequest(url string) *ArxRequest {
	request := &ArxRequest{}
	request.Url = url
	request.Method = "GET"
	request.Header = make(http.Header)
	request.EnableCookie = false
	request.TryTimes = DEFALUT_TRY_TIMES
	request.ConnTimeout = DEFALUT_CONNECT_TIMEOUT
	request.DownloadTimeout = DEFALUT_DOWNLOAD_TIMEOUT
	request.TempStrMap = make(map[string]string)
	return request
}

func (self *ArxRequest) Clone() *ArxRequest {
	request := NewArxRequest(self.Url)
	request.ProcerName = string(self.ProcerName)
	request.Method = self.Method
	request.EnableCookie = self.EnableCookie
	request.TryTimes = self.TryTimes
	request.ConnTimeout = self.ConnTimeout
	request.DownloadTimeout = self.DownloadTimeout
	request.IDDownloader = self.IDDownloader
	request.PostData = string(self.PostData)
	request.Priority = self.Priority
	for name, headers := range self.Header {
		name = string(name)
		for _, h := range headers {
			request.Header.Add(name, h)
		}
	}
	for k, v := range self.TempStrMap {
		request.TempStrMap[k] = string(v)
	}
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

// 返回临时缓存数据
func (self Temp) get(key string, defaultValue interface{}) interface{} {
	defer func() {
		if p := recover(); p != nil {
			log.Errorf(" *     Request.Temp.Get(%v): %v", key, p)
		}
	}()

	var (
		err error
		b   = []byte(self[key].(string))
	)

	if reflect.TypeOf(defaultValue).Kind() == reflect.Ptr {
		err = json.Unmarshal(b, defaultValue)
	} else {
		err = json.Unmarshal(b, &defaultValue)
	}
	if err != nil {
		log.Errorf(" *     Request.Temp.Get(%v): %v", key, err)
	}
	return defaultValue
}

func (self Temp) set(key string, value interface{}) Temp {
	b, err := json.Marshal(value)
	if err != nil {
		log.Errorf(" *     Request.Temp.Set(%v): %v", key, err)
	}
	self[key] = string(b)
	return self
}
