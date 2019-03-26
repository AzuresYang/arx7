package request

import "time"

type RequestManager struct {
	wait_push_req_queue chan *ArxRequest //等待推送到主服务的req无锁队列

}

func (self *RequestManager) Init(max_queue_len int) {
	self.wait_push_req_queue = make(chan *ArxRequest, max_queue_len)
}

func (self *RequestManager) Start() {

}

func (self *RequestManager) Stop() {

}

// 获取一个请求， 等可以连上redis之后，从redis中获取
func (self *RequestManager) GetRequest() (req *ArxRequest) {
	click := time.After(10 * time.Second)
	select {
	case req = <-self.wait_push_req_queue:
		return req
	case <-click:
		return nil
	}
}

// 无锁队列入列
func (self *RequestManager) AddNeedGrabRequest(req *ArxRequest, timeout time.Duration) bool {
	click := time.After(timeout)
	select {
	case self.wait_push_req_queue <- req:
		return true
	case <-click:
		return false
	}
}
