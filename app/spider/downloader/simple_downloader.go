/*
 * @Author: rayou
 * @Date: 2019-03-27 10:11:58
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-10 00:20:15
 */

package downloader

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"fmt"
	"io"
	http "net/http"

	"github.com/AzuresYang/arx7/app/processor"
	"github.com/AzuresYang/arx7/app/spider/downloader/context"
	arxrequest "github.com/AzuresYang/arx7/app/spider/downloader/request"
	httpUtil "github.com/AzuresYang/arx7/util/httpUtil"
	log "github.com/sirupsen/logrus"
)

type SimpleDownloader struct {
}

func (self *SimpleDownloader) Download(procer processor.Processor, req *arxrequest.ArxRequest) *context.CommContext {
	if err := req.Prepare(); err != nil {
		log.Error("request prepare fail.msg:", err.Error())
		return nil
	}

	var ctx *context.CommContext
	response, err := self.doRequest(req)
	if err != nil {
		log.WithFields(log.Fields{
			"URL":    req.Url,
			"Method": req.Method,
			"errmsg": err.Error(),
		}).Error("download fail.")
	} else {

		log.WithFields(log.Fields{
			"URL":    req.Url,
			"Method": req.Method,
		}).Debug("download succ.")

		ctx = context.GetNewContext()
		ctx.Request = req
		ctx.Response = response
		ctx.Status = response.StatusCode
		// 根据返回的类型转换字节流
		switch response.Header.Get("Content-Encoding") {
		case "gzip":
			var gzipReader *gzip.Reader
			gzipReader, err = gzip.NewReader(response.Body)
			if err == nil {
				response.Body = gzipReader
			}
		case "deflate":
			response.Body = flate.NewReader(response.Body)
		case "zlib":
			var readCloser io.ReadCloser
			readCloser, err = zlib.NewReader(response.Body)
			if err == nil {
				response.Body = readCloser
			}
		}
		// log.Infof("ctx is :%#v", ctx)
		return ctx
	}
	return ctx
}

func (self *SimpleDownloader) doRequest(req *arxrequest.ArxRequest) (response *http.Response, err error) {
	switch req.Method {
	case "GET":
		log.WithFields(log.Fields{
			"URL":    req.Url,
			"Header": httpUtil.GetHeaderString(&req.Header),
		}).Trace()
		client := &http.Client{}
		request, _ := http.NewRequest(req.Method, req.Url, nil)
		request.Header = req.Header
		response, _ := client.Do(request)
		if response.StatusCode < 400 {
			return response, nil
		} else { // 大于400的错误码都是有问题的
			err = errors.New(fmt.Sprintf("download fail.code:%d", response.StatusCode))
			// 做一个错误记录request.
			return nil, err
		}

	case "POST":
		log.WithFields(log.Fields{
			"URL": req.Url,
		}).Error("now not support post")
		return nil, errors.New("unsupport post method")
	}
	return nil, errors.New("could go here")
}
