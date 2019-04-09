package record

import (
	"sync/atomic"
	"time"

	"github.com/AzuresYang/arx7/app/spider/downloader/request"
)

type CountType uint32

const (
	COUNT_DOWNLOAD_SUCC CountType = iota
	COUNT_DOWNLOAD_FAIL
)

var (
	AppStartTime  time.Time               // app开始运行的时间点
	TaskStartTime time.Time               // app开始运行任务的时间点
	recordCount   [5]uint64 = [5]uint64{} // 统计用数组， 注意和CountType定义的多少一致
)

func ResetRecordCount() {
	recordCount = [5]uint64{}
}

func GetCount(countType CountType) uint64 {
	return recordCount[countType]
}

func CountAddOne(countType CountType) {
	atomic.AddUint64(&recordCount[countType], 1)
}

func DownloadSuccReq(req *request.ArxRequest, msg string) {

}

func DownloadFailReq(req *request.ArxRequest, msg string) {

}
