/*
 * @Author: rayou
 * @Date: 2019-03-25 22:21:15
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-25 22:49:43
 */
package crawler

type (
	Crawler interface {
		Init( /*一个spider分析器，配置文件*/ ) Crawler // 初始化
		Run()                                // 运行
		Stop()                               // 停止运行
		GetId()                              // 获取ID
	}

	
	crawler struct {
		// 一个采集规则分析器spider， 一个request控制器request-controller,下载器 downloader, 存储pipleLine
		spider      int
		request_mgr int
		downloader  int
		pipeline    int
		id          int      // id
		pause       [2]int64 //[距离下个请求的最短时常， 距离下个请求的最长时长]
	}
)
