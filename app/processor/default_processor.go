/*
 * @Author: rayou
 * @Date: 2019-03-30 09:34:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-10 02:47:36
 */
package processor

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/AzuresYang/arx7/app/spider/downloader/request"

	"github.com/AzuresYang/arx7/app/spider/downloader/context"
	log "github.com/sirupsen/logrus"
)

type DefaultProcessor struct {
	EnableCookie bool // 所有请求是否使用cookie记录
	// Namespace       func(self *Spider) string                                  // 命名空间，用于输出文件、路径的命名
	// SubNamespace    func(self *Spider, dataCell map[string]interface{}) string // 次级命名，用于输出文件、路径的命名，可依赖具体数据内容
	// 以下字段系统自动赋值
	Id int
}

func NewDefaultProcessor() DefaultProcessor {
	procer := DefaultProcessor{
		Id: 5,
	}
	return procer
}
func (self *DefaultProcessor) GetName() string {
	return "default"
}

func findNewRequest(ctx *context.CommContext) {
	flysnowRegexp := regexp.MustCompile(`<a href="http://www.xbiquge.la/([\d]+/[\d]+)/">(.+)</a></li>`)
	response := ctx.Response
	body, _ := ioutil.ReadAll(response.Body)
	bodystr := string(body)
	log.Infof("download first page")
	lines := strings.Split(bodystr, "\n")
	max_num := 50
	for _, line := range lines {
		params := flysnowRegexp.FindStringSubmatch(line)
		if params != nil {
			max_num--
			if max_num <= 0 {
				return
			}
			new_url := fmt.Sprintf("http://www.xbiquge.la/%s", params[1])
			name := params[2]
			new_req := ctx.Request.Clone()
			new_req.Url = new_url
			new_req.TempStrMap = make(map[string]string)
			new_req.TempStrMap["title"] = name
			new_req.TempStrMap["is_stop"] = "true"
			request.RequestMgr.AddNeedGrabRequest(new_req, 1*time.Second)
			log.Infof("new req：[%s|%s]", name, new_url)
		}
	}
	scanner := bufio.NewScanner(ctx.Response.Body)
	for scanner.Scan() {
		line := scanner.Text()
		log.Infof("scanner read line:%s", line)
	}

}
func saveXiaoShuo(ctx *context.CommContext) {
	file_dir := "E:/crawler_data/"
	title := ctx.Request.TempStrMap["title"]
	file_name := file_dir + fmt.Sprintf("%s.txt", title)
	f, err := os.Create(file_name)
	if err != nil {
		log.Errorf("[procer|process]:create file[%s] error.", file_name)
		return
	}
	response := ctx.Response
	body, _ := ioutil.ReadAll(response.Body)
	bodystr := string(body)
	f.WriteString(bodystr)
	defer f.Close()
	log.Infof("获得小说：%s", file_name)
}
func (self *DefaultProcessor) Process(ctx *context.CommContext) (ret int, msg string) {
	// log.Info("ready pro context:" + ctx.Request.Url)
	log.WithFields(log.Fields{
		"URL": ctx.Request.Url,
		// "Data:": bodystr,
	}).Debug("get data .ready proce")
	ret = 0
	msg = "procer succ"
	is_stop := ctx.Request.TempStrMap["is_stop"]
	if is_stop == "" {
		findNewRequest(ctx)
	} else {
		saveXiaoShuo(ctx)
	}

	// response := ctx.Response
	// body, _ := ioutil.ReadAll(response.Body)
	// bodystr := string(body)
	// f.WriteString(bodystr)
	// defer f.Close()

	return ret, msg
}

func (self *DefaultProcessor) Free() {
	log.Info("default processor, use Freee Method")
	return
}

func (self *DefaultProcessor) GetOneProcessor() Processor {
	procer := DefaultProcessor{
		Id: 5,
	}
	return &procer
}
