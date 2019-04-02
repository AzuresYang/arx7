/*
 * @Author: rayou
 * @Date: 2019-04-02 20:57:36
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-02 22:15:25
 */
package crawler

import (
	"sync"
	"time"

	"github.com/AzuresYang/arx7/app/status"
)

// crawler池
type (
	CrawlerPool interface {
		Set(size uint32) uint32
		Get() Crawler
		Free(Crawler)
		Stop()
	}

	crawlerpool struct {
		capacity uint32
		count    uint32
		can_use  chan Crawler
		pool     []Crawler
		status   int
		sync.RWMutex
	}
)

func NewCrawlerPool(size uint32) CrawlerPool {
	cp := new(crawlerpool)
	cp.Set(size)
	return cp
}

func (self *crawlerpool) Set(size uint32) uint32 {
	self.Lock()
	defer self.Unlock()
	status = status.RUN
	var pool_size uint32 = 1
	if size > 0 {
		pool_size = size
	}
	self.capacity = pool_size
	self.count = 0
	self.can_use = make(chan Crawler, pool_size)
	for _, crawler := range self.pool {
		if self.count < self.capacity {
			self.can_use <- crawler
			self.count++
		}
	}
	return pool_size
}

func (self *crawlerpool) Get() Crawler {
	var crawler Crawler
	// 谁先调用，谁先获取
	self.Lock()
	defer self.Unlock()
	for {
		if self.status == status.STOP {
			return nil
		}
		select {
		case crawler = <-self.can_use:
			return crawler
		default:
			if self.count < self.capacity {
				crawler = NewCrawler(int(self.count))
				self.pool = append(self.pool, crawler)
				return crawler
			}
		}
		// 迟0.5秒后才获取， 太快获取也是空
		time.Sleep(500 * time.Millisecond)
	}
}

func (self *crawlerpool) Free(crawler Crawler) {
	self.RLock()
	if self.status != status.STOP {
		self.can_use <- crawler
	}
	self.RUnlock()
}

// 终止池子中所有的crawler， 停止提供crawler
func (self *crawlerpool) Stop() {
	self.Lock()
	defer self.Unlock()
	if self.status == status.STOP
	{
		return
	}
	self.status = status.STOP
	close(self.can_use)
	self.can_use = nil
	for _, crawler := range self.pool{
		crawler.Stop()
	}
}
