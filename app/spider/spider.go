/*
 * @Author: rayou
 * @Date: 2019-03-25 22:50:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-26 23:22:41
 */
package spider

import "sync"

type Spider struct {
	EnableCookie bool // 所有请求是否使用cookie记录
	// Namespace       func(self *Spider) string                                  // 命名空间，用于输出文件、路径的命名
	// SubNamespace    func(self *Spider, dataCell map[string]interface{}) string // 次级命名，用于输出文件、路径的命名，可依赖具体数据内容
	// 以下字段系统自动赋值
	id      int
	subName string // 由Keyin转换为的二级标识名
	status  int    // 执行状态
	lock    sync.RWMutex
	once    sync.Once
}
