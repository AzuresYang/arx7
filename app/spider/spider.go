/*
 * @Author: rayou
 * @Date: 2019-04-02 19:51:55
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-02 20:02:29
 */
package spider

import "github.com/AzuresYang/arx7/app/spider/downloader/request"

type Spider struct {
	to_stop chan int
	to_end  chan int
}

func (spider *Spider) Init() {
	request.RequestManager.Init(20)
	// 包括crawler 池的初始化 ？？？

}

func (spider *Spider) Stop() {

}

// 开始运行
func (spider *Spider) Run() {
	spider.Init()
}

func (spider *Spider) run() {
}
