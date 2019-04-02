package record

import (
	"sync/atomic"
	"time"
)

type CountType uint32

const (
	COUNT_DOWNLOAD_SUCC CountType = iota
	COUNT_DOWNLOAD_FAIL
)

var (
	AppStartTime  time.Time // app开始运行的时间点
	TaskStartTime time.Time // app开始运行任务的时间点
	recordCount   [2]uint64 // 统计用数组， 注意和CountType定义的多少一致
)

func ResetRecordCount() {
	recordCount = [2]uint64{}
}

func GetCount(countType CountType) uint64 {
	return recordCount(countType)
}

func CountAddOne(countType CountType) {
	atomic.AddUint64(&recordCount[countType], 1)
}
