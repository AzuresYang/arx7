/*
 * @Author: rayou
 * @Date: 2019-04-16 23:43:27
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-16 23:44:07
 */
package stringUtil

import (
	"hash/crc32"
)

// 计算字符串hash值
func Hash(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}
