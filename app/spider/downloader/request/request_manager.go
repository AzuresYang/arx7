package request

import "time"
import "fmt"

type RequestManager struct {
	wait_push_req_queue chan *ArxRequest //等待推送到主服务的req无锁队列

}

var RequestMgr = RequestManager{}

func (self *RequestManager) Init(max_queue_len int) {
	self.wait_push_req_queue = make(chan *ArxRequest, max_queue_len)
}

func (self *RequestManager) Start() {

}

func (self *RequestManager) Stop() {

}

// 获取一个请求， 等可以连上redis之后，从redis中获取
func (self *RequestManager) GetRequest(timeout time.Duration) (req *ArxRequest) {
	click := time.After(timeout)
	select {
	case req = <-self.wait_push_req_queue:
		return req
	case <-click:
		return nil
	}
}

// 无锁队列入列， 需要维护一个去重的URL队列
func (self *RequestManager) AddNeedGrabRequest(req *ArxRequest, timeout time.Duration) bool {
	self.wait_push_req_queue <- req
	fmt.Printf("[requestmgr]get new url:%s", req.Url)
	return true
	// click := time.After(timeout)
	// select {
	// case
	// 	return true
	// case <-click:
	// 	return false
	// }
}

// 下载失败的链接处理
func (self *RequestManager) AddDownLoadFailReqeust(req *ArxRequest, msg string) bool {
	return true
}

// 记录下载成功链接
func (self *RequestManager) AddDownloadSuccReq(req *ArxRequest) {

}

func (self *RequestManager) AddProcessReq(req *ArxRequest) {

}

// 处理成功链接
