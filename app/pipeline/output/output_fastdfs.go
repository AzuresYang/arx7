/*
 * @Author: rayou
 * @Date: 2019-04-05 22:05:31
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-23 00:08:35
 */

package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/AzuresYang/arx7/app/arxmonitor/monitorHandler"
	"github.com/AzuresYang/arx7/app/pipeline"
	"github.com/AzuresYang/arx7/app/status"
	log "github.com/sirupsen/logrus"
)

type CellFastDfs struct {
	Dir      string // 目的目录， 前面不需要加斜杠。示例： “test/ddd”
	Data     []byte
	FileName string // 文件名
}

type CollectFastDfsData struct {
	Cells []CellFastDfs
}

type OutputFastDfs struct {
	DfsAddr  string
	TaskName string // 任务名
	// IfDisTask bool   // 是否区分任务保存
}

type UpLoadResult struct {
	Url    string `json:"url"`
	Md5    string `json:"md5"`
	Path   string `json:"path"`
	Domain string `json:"domain"`
	Scene  string `json:"scene"`
	//Just for Compatibility
	Scenes  string `json:"scenes"`
	Retmsg  string `json:"retmsg"`
	Retcode int    `json:"retcode"`
	Src     string `json:"src"`
}

func NewCollectFastDfsData() *CollectFastDfsData {
	return &CollectFastDfsData{
		Cells: make([]CellFastDfs, 0, 2),
	}
}

func (self *CollectFastDfsData) Add(dir string, fileName string, data []byte) {
	cell := CellFastDfs{
		Dir:      dir,
		FileName: fileName,
		Data:     data,
	}
	self.Cells = append(self.Cells, cell)
}

func (self *CollectFastDfsData) ToCollectData() *pipeline.CollectData {
	cd := pipeline.NewCollectData(1)
	cd.Cell["dfs_data"] = self
	return cd
}

func (self *OutputFastDfs) Reset(dfsAddr string, taskName string) {
	self.DfsAddr = dfsAddr
	self.TaskName = taskName
}

func (self *OutputFastDfs) Init() error {
	// 测试dfs的是否可以使用

	return nil
}

func postFile(filename string, target_url string, data []byte, params map[string]string) (*http.Response, error) {
	code_info := "OutputFastDfs.postFile"
	body_buf := bytes.NewBufferString("")
	body_writer := multipart.NewWriter(body_buf)
	for k, v := range params {
		body_writer.WriteField(k, v)
	}
	// use the body_writer to write the Part headers to the buffer
	_, err := body_writer.CreateFormFile("file", filename)
	if err != nil {
		log.Errorf("[%s]error writing to buffer", code_info)
		return nil, err
	}

	// the file data will be the second part of the body
	fh := bytes.NewBuffer(data)
	// need to know the boundary to properly close the part myself.
	boundary := body_writer.Boundary()
	//close_string := fmt.Sprintf("\r\n--%s--\r\n", boundary)
	close_buf := bytes.NewBufferString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// use multi-reader to defer the reading of the file data until
	// writing to the socket buffer.
	request_reader := io.MultiReader(body_buf, fh, close_buf)
	req, err := http.NewRequest("POST", target_url, request_reader)
	if err != nil {
		log.Errorf("[%s]new request fail:%s", code_info, err.Error())
		return nil, err
	}

	// Set headers for multipart, and Content Length
	req.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	req.ContentLength = int64(len(data)) + int64(body_buf.Len()) + int64(close_buf.Len())
	log.Tracef("req:%+v", req)
	return http.DefaultClient.Do(req)

}
func (self *OutputFastDfs) uploadData(data *CellFastDfs) {
	code_info := "OuputFastDfs.uploadData"
	url := self.DfsAddr + "/upload"
	param := make(map[string]string)
	param["output"] = "json"
	param["scene"] = "default"
	param["path"] = data.Dir
	log.Debugf("ready upload file:%s/%s", data.Dir, data.FileName)
	resp, err := postFile(data.FileName, url, data.Data, param)
	if err != nil {
		log.Errorf("[%s]upload file fail:%s", code_info, err.Error())
		return
	}
	// log.Debugf("[%s]upload ret:%+v.", code_info, resp)
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("[%s]read response  fail:%s", code_info, err.Error())
		return
	}
	ret := &UpLoadResult{}
	if err = json.Unmarshal(body, ret); err != nil {
		log.Errorf("[%s]Unmarshal upload ret fail:%s", code_info, string(body))
		return
	}
	// 打印保存结果
	if ret.Retcode == 0 {
		monitorHandler.AddOne(status.MONI_SYS_DFS_UPLOAD_SUCC)
		log.Infof("[%s]collect data[%s] succ.", code_info, data.FileName)
	} else {
		monitorHandler.AddOne(status.MONI_SYS_DFS_UPLOAD_FAIL)
		log.Infof("[%s]collect data[%s] fail.ret:%s", code_info, data.FileName, string(body))
	}
}

func (self *OutputFastDfs) CollectData(collectData *pipeline.CollectData) {
	// 类型转换
	code_info := "OutputFastDfs.CollectData"
	data, ok := collectData.Cell["dfs_data"].(*CollectFastDfsData)
	if !ok {
		log.Errorf("[%s] CollectData not found dfs data.type is :%T", code_info, collectData.Cell["dfs_data"])
		return
	}
	for _, cell := range data.Cells {
		self.uploadData(&cell)
	}
}

func (self *OutputFastDfs) Stop() {

}

func (self *OutputFastDfs) GetName() string {
	return "OutputFile"
}
