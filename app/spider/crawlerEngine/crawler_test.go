package crawlerEngine

import (
	"testing"

	// "github.com/AzuresYang/arx7/app/pipeline"
	"github.com/AzuresYang/arx7/app/processor"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
	// log "github.com/sirupsen/logrus"
	// "github.com/AzuresYang/arx7/config"
	// "github.com/AzuresYang/arx7/util/record"
)

func buildReq(procer processor.Processor) {
	reqs := []string{
		"http://www.baidu.com",
		"https://www.bilibili.com/read/douga?from=category_0",
		"https://blog.csdn.net/EasternUnbeaten/article/details/72355127",
		"https://www.bilibili.com/read/cv2385487?from=category_2",
	}
	for _, s := range reqs {
		req := request.NewArxRequest(s)
		req.ProcerName = procer.GetName()
		request.RequestMgr.AddNeedGrabRequest(req)
	}
}
func TestCrawler(t *testing.T) {
	// request.RequestMgr.Init(10)
	// procer := processor.NewDefaultProcessor()
	// processor.Manager.Register(&procer)
	// buildReq(&procer)
	// dl := &downloader.SimpleDownloader{}
	// craw := &crawler{
	// 	id:      12,
	// 	if_stop: false,
	// 	pause:   [2]int64{1, 2},
	// }
	// craw.SetDownloader(dl)
	// craw.Run()
	// t.Log("done")
}
