/*
 * @Author: rayou
 * @Date: 2019-04-14 13:06:51
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-14 13:25:12
 */
package message

// 通信使用的命令字
const (
	MSG_CONN_END uint32 = 1000 // 停止交互链接

	// spider通信消息
	MSG_REQ_STAET_SPIDER    uint32 = 1001 // 启动爬虫
	MSG_RSP_STAET_SPIDER    uint32 = 1002 // 启动爬虫
	MSG_REQ_STOP_SPIDER     uint32 = 1003 // 停止爬虫
	MSG_RSP_STOP_SPIDER     uint32 = 1004 // 停止爬虫
	MSG_REQ_DELETE_TASK     uint32 = 1005 // 删除爬虫任务
	MSG_RSP_DELETE_TASK     uint32 = 1006 // 删除爬虫任务
	MSG_REQ_SCALE_SPIDER    uint32 = 1007 // 爬虫节点扩缩容
	MSG_RSP_SCALE_SPIDER    uint32 = 1008 // 爬虫节点扩缩容
	MSG_REQ_GET_SPIDER_INFO uint32 = 1009 // 查看爬虫节点任务
	MSG_RSP_GET_SPIDER_INFO uint32 = 1010 // 查看爬虫节点任务

	// 监控信息
	MSG_MONITOR_INFO uint32 = 2001

	// CMD命令行操作信息
	MSG_ARXCMD_START_SPIDER    uint32 = 3001
	MSG_ARXCMD_STOP_STOP       uint32 = 3002
	MSG_ARXCMD_SCALE_STOP      uint32 = 3003 // 扩缩容
	MSG_ARXCMD_DELTE_TASK      uint32 = 3004 // 删除任务
	MSG_ARXCMD_GET_SPIDER_INFO uint32 = 3005 // 获取爬虫节点信息
)
