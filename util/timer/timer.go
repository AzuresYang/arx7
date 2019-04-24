/*
 * @Author: rayou
 * @Date: 2019-04-24 21:16:40
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 21:28:58
 * @breif :一个简单的定时器
 */

package timer

import "time"

type Timer struct {
	ifStop   bool
	task     func()
	interval time.Duration
}

func New() *Timer {
	return &Timer{
		ifStop:   true,
		interval: 1 * time.Second,
	}
}

// 直接运行即可
func (self *Timer) RunTask(interval time.Duration, f func()) {
	self.interval = interval
	self.ifStop = false
	self.task = f
	// 定时执行任务
	ticker := time.NewTicker(self.interval)
	go func() {
		for !self.ifStop {
			select {
			case <-ticker.C:
				self.task()
			}
		}
	}()
}

func (self *Timer) IfStop() bool {
	return self.ifStop
}

func (self *Timer) Stop() {
	self.ifStop = true
}
