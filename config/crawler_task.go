/*
 * @Author: rayou
 * @Date: 2019-04-03 20:37:41
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-03 20:40:23
 */
package config

type CrawlerTask struct {
	TaskName        string // 任务名
	TaskId          int    // 任务ID
	CrawlerCapacity uint32 // 爬虫数量
	RedisAddress    string
	RedisPort       uint32
	RedisQueueName  string // redis中的队列名
}
