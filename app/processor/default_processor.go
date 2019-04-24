/*
 * @Author: rayou
 * @Date: 2019-03-30 09:34:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-22 23:53:07
 */
package processor

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/AzuresYang/arx7/app/pipeline"
	"github.com/AzuresYang/arx7/app/pipeline/output"
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

func findNewRequest(ctx *context.CommContext) *pipeline.CollectData {
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
				return nil
			}
			new_url := fmt.Sprintf("http://www.xbiquge.la/%s", params[1])
			name := params[2]
			new_req := ctx.Request.Clone()
			new_req.Url = new_url
			new_req.TempStrMap = make(map[string]string)
			new_req.TempStrMap["title"] = name
			new_req.TempStrMap["is_stop"] = "true"
			request.RequestMgr.AddNeedGrabRequest(new_req)
			log.Infof("new req：[%s|%s]", name, new_url)
		}
	}
	scanner := bufio.NewScanner(ctx.Response.Body)
	for scanner.Scan() {
		line := scanner.Text()
		log.Infof("scanner read line:%s", line)
	}
	return nil
}
func saveXiaoShuo(ctx *context.CommContext) *pipeline.CollectData {
	file_dir := "E:/crawler_data/"
	title := ctx.Request.TempStrMap["title"]
	file_name := file_dir + fmt.Sprintf("%s.txt", title)
	f, err := os.Create(file_name)
	if err != nil {
		log.Errorf("[procer|process]:create file[%s] error.", file_name)
	}
	response := ctx.Response
	body, _ := ioutil.ReadAll(response.Body)
	dfs_data := output.NewCollectFastDfsData()
	dfs_data.Add("crawler_data", fmt.Sprintf("%s.txt", title), body)
	bodystr := string(body)
	f.WriteString(bodystr)
	defer f.Close()
	log.Infof("获得小说：%s", file_name)
	return dfs_data.ToCollectData()
}
func (self *DefaultProcessor) Process(ctx *context.CommContext) *pipeline.CollectData {
	// log.Info("ready pro context:" + ctx.Request.Url)
	log.WithFields(log.Fields{
		"URL": ctx.Request.Url,
		// "Data:": bodystr,
	}).Debug("get data .ready proce")
	is_stop := ctx.Request.TempStrMap["is_stop"]
	if is_stop == "" {
		return findNewRequest(ctx)
	} else {
		return saveXiaoShuo(ctx)
	}

	// response := ctx.Response
	// body, _ := ioutil.ReadAll(response.Body)
	// bodystr := string(body)
	// f.WriteString(bodystr)
	// defer f.Close()

	return nil
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
