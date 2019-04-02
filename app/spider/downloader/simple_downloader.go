/*
 * @Author: rayou
 * @Date: 2019-03-27 10:11:58
 * @Last Modified by: rayou
 * @Last Modified time: 2019-03-27 19:34:01
 */

package downloader

import (
	"fmt"
	"io/ioutil"
	http "net/http"

	"github.com/AzuresYang/arx7/app/processor"
	"github.com/AzuresYang/arx7/app/spider/downloader/context"
	arxrequest "github.com/AzuresYang/arx7/app/spider/downloader/request"
	httpUtil "github.com/AzuresYang/arx7/util/http"
	log "github.com/sirupsen/logrus"
)

type SimpleDownloader struct {
}

func (self *SimpleDownloader) Download(procer *processor.Processor, req *arxrequest.ArxRequest) *context.CommContext {
	if err := req.Prepare(); err != nil {
		log.Error("request prepare fail.msg:", err.Error())
		return nil
	}
	switch req.Method {
	case "GET":
		log.WithFields(log.Fields{
			"URL":    req.Url,
			"Header": httpUtil.GetHeaderString(&req.Header),
		})
		client := &http.Client{}
		request, _ := http.NewRequest(req.Method, req.Url, nil)
		request.Header = req.Header
		response, _ := client.Do(request)
		if response.StatusCode == 200 {
			body, _ := ioutil.ReadAll(response.Body)
			bodystr := string(body)
			fmt.Println(bodystr)
		}

	case "POST":
		log.WithFields(log.Fields{
			"URL": req.Url,
		}).Error("now not support post")
	}
	return nil
}
