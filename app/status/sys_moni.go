/*
 * @Author: rayou
 * @Date: 2019-04-21 21:05:08
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-24 21:49:36
 */

package status

const (
	MONI_SYS_HEART_APP    uint32 = 101  // 程序心跳
	MONI_SYS_HEART_ENGINE uint32 = 102  // 引擎心跳
	MONI_SYS_DOWNLOAD     uint32 = 1002 // 下载结果，0成功，其余为失败码
	MONI_SYS_REQUEST_ADD  uint32 = 1004 // 加入到了一个新链接
	MONI_SYS_REQUEST_GET  uint32 = 1005 // 获取到一个新链接
)
