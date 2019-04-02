/*
 * @Author: rayou
 * @Date: 2019-03-25 22:50:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-02 22:40:11
 */
package processor

// 处理下载数据， 对需要解析不同页面，可以自己实现不同的处理器
type Processor interface {
	GetProcessorName() string // 处理器的代称
	// 获取一个处理器
	// 这里这么设计的原因是：对于不同的解析器, 下载好页面之后，需要找到处理这个页面的处理器
	// 一种是直接使用Process()方法， 一种是生成，
	GetOneProcessor() Processor
	Free(Processor)
}
