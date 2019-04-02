/*
 * @Author: rayou
 * @Date: 2019-03-25 22:21:15
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-02 22:24:41
 */
package crawler

type (
	Crawler interface {
		// Init( /*一个spider分析器，配置文件*/ ) error // 初始化
		Run()       // 运行
		Stop()      // 停止运行
		GetId() int // 获取ID
	}

	crawler struct {
		// 一个采集规则分析器spider， 一个request控制器request-controller,下载器 downloader, 存储pipleLine
		spider      int
		request_mgr int
		downloader  int
		id          int // id
		if_stop     bool
		pause       [2]int64 //[距离下个请求的最短时常， 距离下个请求的最长时长]
	}
)

// 新建一个Crawler
func NewCrawler(id int) Crawler {
	return &crawler{
		id:      id,
		if_stop: false,
	}
}

func (self *crawler) Init() error {
	return nil
}

func (self *crawler) Run() {

}

func (self *crawler) Stop() {
	self.if_stop = true
}

func (self *crawler) GetId() int {
	return self.id
}
