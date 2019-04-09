package spider

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	// "github.com/AzuresYang/arx7/app/pipeline"
	"github.com/AzuresYang/arx7/app/processor"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
	log "github.com/sirupsen/logrus"
	// "github.com/AzuresYang/arx7/config"
	// "github.com/AzuresYang/arx7/util/record"
)

func buildReq(procer processor.Processor) {
	reqs := []string{
		"http://www.xbiquge.la/paihangbang/",
	}
	for _, s := range reqs {
		req := request.NewArxRequest(s)
		req.ProcerName = procer.GetName()
		// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		// req.Header.Set("Host", "image.baidu.com")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/71.0.3578.98 Safari/537.36")

		request.RequestMgr.AddNeedGrabRequest(req, 2*time.Second)
	}
}
func TestSpider(t *testing.T) {
	log.SetLevel(log.TraceLevel)
	request.RequestMgr.Init(10)
	procer := processor.NewDefaultProcessor()
	processor.Manager.Register(&procer)
	fmt.Println("get procer succ")
	buildReq(&procer)
	sp := &Spider{}
	sp.Start()
	time.Sleep(5 * time.Second)
	t.Log("done")
}

func TestRegx(t *testing.T) {
	ctx := `<li>1<a href="http://www.xbiquge.la/8/8226/">通天武尊</a></li>`

	flysnowRegexp := regexp.MustCompile(`<a href="http://www.xbiquge.la/([\d]+/[\d]+)/">(.+)</a></li>`)
	params := flysnowRegexp.FindStringSubmatch(ctx)

	for i, param := range params {
		fmt.Printf("[%d]%s\n", i, param)
	}

}
