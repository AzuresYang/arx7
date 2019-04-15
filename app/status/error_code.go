/*
 * @Author: rayou
 * @Date: 2019-04-15 15:08:55
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-15 16:20:13
 */
package status

const (
	ERR_UNSERIALIZE_FAIL uint32 = 101
	ERR_RESP_ERROR uint32 = 102

	ERR_START_SPIDER_FAIL         uint32 = 200 // 启动crawler错误
	ERR_START_SPIDER_FAIL_RUNNING uint32 = 201 // 启动crawler错误，正在运行中
)
