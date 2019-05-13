/*
 * @Author: rayou
 * @Date: 2019-03-30 09:34:25
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-09 23:30:21
 */
package biqu

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"time"

	"github.com/AzuresYang/arx7/app/arxmonitor/monitorHandler"
	"github.com/AzuresYang/arx7/app/pipeline"
	"github.com/AzuresYang/arx7/app/pipeline/output"
	"github.com/AzuresYang/arx7/app/processor"
	"github.com/AzuresYang/arx7/app/spider/downloader/context"
	"github.com/AzuresYang/arx7/app/spider/downloader/request"
	"github.com/AzuresYang/arx7/app/status"
	log "github.com/sirupsen/logrus"
)

type BiquProcessor struct {
	EnableCookie bool // 所有请求是否使用cookie记录
	// Namespace       func(self *Spider) string                                  // 命名空间，用于输出文件、路径的命名
	// SubNamespace    func(self *Spider, dataCell map[string]interface{}) string // 次级命名，用于输出文件、路径的命名，可依赖具体数据内容
	// 以下字段系统自动赋值
	Id int
}

func NewProcessor() BiquProcessor {
	procer := BiquProcessor{
		Id: 2,
	}
	return procer
}

// 实现自动注册
var (
	nouse_procer BiquProcessor = NewProcessor()
	_            error         = processor.Manager.Register(&nouse_procer)
)

func (self *BiquProcessor) GetName() string {
	return "biqu"
}

func findNewRequest(ctx *context.CommContext) *pipeline.CollectData {
	flysnowRegexp := regexp.MustCompile(`<a href="http://www.xbiquge.la/([\d]+/[\d]+)/">(.+)</a></li>`)
	response := ctx.Response
	body, _ := ioutil.ReadAll(response.Body)
	bodystr := string(body)
	log.Infof("download first page")
	lines := strings.Split(bodystr, "\n")
	max_num := 20
	get_num := 0
	batch_num := 30 // 提交这个数量后，停止一下，否则机器容易卡
	for _, line := range lines {
		params := flysnowRegexp.FindStringSubmatch(line)
		if params != nil {
			max_num--
			if max_num <= 0 {
				return nil
			}
			batch_num--
			if batch_num < 0 {
				time.Sleep(1 * time.Second)
				batch_num = 0
			}
			new_url := fmt.Sprintf("http://www.xbiquge.la/%s", params[1])
			name := params[2]
			new_req := ctx.Request.Clone()
			new_req.Url = new_url
			new_req.TempStrMap = make(map[string]string)
			new_req.TempStrMap["title"] = name
			new_req.TempStrMap["stage"] = "1"
			request.RequestMgr.AddNeedGrabRequest(new_req)
			get_num++
			log.Tracef("new req：[%s|%s]", name, new_url)
		}
	}
	log.Infof("笔趣阁首页榜小说数量:%d", get_num)
	scanner := bufio.NewScanner(ctx.Response.Body)
	for scanner.Scan() {
		line := scanner.Text()
		log.Infof("scanner read line:%s", line)
	}
	return nil
}
func saveNovelHomePage(ctx *context.CommContext) *pipeline.CollectData {
	file_dir := "crawler_data/"
	title := ctx.Request.TempStrMap["title"]
	file_name := file_dir + fmt.Sprintf("%s.txt", title)

	response := ctx.Response
	body, _ := ioutil.ReadAll(response.Body)
	bodystr := string(body)
	// f, err := os.Create(file_name)
	// if err != nil {
	// 	log.Errorf("[procer|process]:create file[%s] error.", file_name)
	// }
	// f.WriteString(bodystr)
	// defer f.Close()

	// novel_dir := "biqu_data/" + title
	// os.MkdirAll(novel_dir, os.ModePerm)
	log.Infof("保存小说首页：%s", file_name)
	// 存储接下来的链接
	// 需要抓取的章节太多，分配提交
	max_num := 1000
	batch_num := 40 // 提交这个数量后，停止一下，否则机器容易卡
	lines := strings.Split(bodystr, "\n")
	novel_reg := regexp.MustCompile(`<dd><a href='(.+)' >(.+)</a></dd>`)
	for _, line := range lines {
		batch_num--
		if batch_num < 0 {
			time.Sleep(1 * time.Second)
			batch_num = 0
		}
		params := novel_reg.FindStringSubmatch(line)
		if params != nil {
			new_req := ctx.Request.Clone()
			new_req.Url = fmt.Sprintf("http://www.xbiquge.la/%s", params[1])
			new_req.TempStrMap["stage"] = "2"
			new_req.TempStrMap["title"] = title
			new_req.TempStrMap["chapter"] = params[2]
			new_req.TempStrMap["novel_dir"] = "biqu_data/" + title
			request.RequestMgr.AddNeedGrabRequest(new_req)
			max_num--
			if max_num <= 0 {
				break
			}
		}
	}
	// return nil
	monitorHandler.AddOne(status.MONI_APP_NOVEL_NUM)
	dfs_data := output.NewCollectFastDfsData()
	dfs_data.Add("crawler_data", fmt.Sprintf("%s.txt", title), body)
	return dfs_data.ToCollectData()
}

func saveNovelContent(ctx *context.CommContext) *pipeline.CollectData {
	file_dir := ctx.Request.TempStrMap["novel_dir"]
	chapter := ctx.Request.TempStrMap["chapter"]
	file_name := file_dir + "/" + fmt.Sprintf("%s.txt", chapter)
	response := ctx.Response
	body, _ := ioutil.ReadAll(response.Body)

	// f, err := os.Create(file_name)
	// if err != nil {
	// 	log.Errorf("[procer|process]:create file[%s] error.", file_name)
	// }
	// bodystr := string(body)
	// f.WriteString(bodystr)
	// defer f.Close()
	log.Infof("小说章节[%s:%s],目录:%s", ctx.Request.TempStrMap["title"], chapter, file_name)
	monitorHandler.AddOne(status.MONI_APP_CHAPTER_NUM)
	dfs_data := output.NewCollectFastDfsData()
	dfs_data.Add(file_dir, fmt.Sprintf("%s.txt", chapter), body)
	return dfs_data.ToCollectData()
	// return nil
}
func (self *BiquProcessor) Process(ctx *context.CommContext) *pipeline.CollectData {
	// // log.Info("ready pro context:" + ctx.Request.Url)
	// log.WithFields(log.Fields{
	// 	"URL": ctx.Request.Url,
	// 	// "Data:": bodystr,
	// }).Trace("get data .ready proce")
	switch ctx.Request.TempStrMap["stage"] {
	case "":
		return findNewRequest(ctx)
	case "1":
		return saveNovelHomePage(ctx)
	case "2":
		return saveNovelContent(ctx)
	}

	// response := ctx.Response
	// body, _ := ioutil.ReadAll(response.Body)
	// bodystr := string(body)
	// f.WriteString(bodystr)
	// defer f.Close()

	return nil
}

func (self *BiquProcessor) Free() {
	log.Info("default processor, use Freee Method")
	return
}

func (self *BiquProcessor) GetOneProcessor() processor.Processor {
	procer := BiquProcessor{
		Id: 5,
	}
	return &procer
}
