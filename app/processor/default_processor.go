/*
 * @Author: rayou
 * @Date: 2019-03-30 09:34:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-30 09:36:55
 */
package DefalutProcessor

import "sync"

type DefalutProcessor struct {
	EnableCookie bool // 所有请求是否使用cookie记录
	// Namespace       func(self *Spider) string                                  // 命名空间，用于输出文件、路径的命名
	// SubNamespace    func(self *Spider, dataCell map[string]interface{}) string // 次级命名，用于输出文件、路径的命名，可依赖具体数据内容
	// 以下字段系统自动赋值
	id     int
	status int // 执行状态
	lock   sync.RWMutex
	once   sync.Once
}
