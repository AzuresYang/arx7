/*
 * @Author: rayou
 * @Date: 2019-04-03 20:37:41
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-03 20:40:23
 */
package config

import "time"

type CrawlerTask struct {
	TaskName                    string // 任务名
	TaskId                      int    // 任务ID
	CrawlerNum                  uint32 // 爬虫数量
	RedisAddress                string
	RedisPort                   uint32
	RedisAccount                string        // redis账户名
	RedisPassword               string        // redis账户名
	RedisQueueName              string        // redis中的队列名
	MaxGetRequestNullTimeSecond time.Duration // 长时间内没有新链接时，停止工作crawler的设置， 单位：秒， 为0时表示一直工作
}
