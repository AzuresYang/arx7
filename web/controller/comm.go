/*
 * @Author: rayou
 * @Date: 2019-04-29 19:15:59
 * @Last Modified by: rayou
 * @Last Modified time: 2019-04-30 23:34:04
 */

package controller

type ResponseData struct {
	Status uint32      `json:"Status"`
	Msg    string      `json:"Msg"`
	Data   interface{} `json:"Data"`
}

type FormQueryMonitor struct {
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	SvcId     string `json:"svcid"`
	Metric    string `json:"metric"`
	Classfy   string `json:"classfy"`
}

type MonitorData struct {
	Id      int
	Svcid   int
	Metric  int
	Classfy int
	Value   int
	Ip      string
	Time    int64 // 给前端的是这个数据，就这样吧
}

type MonitorDataWrapper struct {
	Data []MonitorData
}

func (self *MonitorDataWrapper) Len() int { // 重写 Len() 方法
	return len(self.Data)
}
func (self *MonitorDataWrapper) Swap(i, j int) { // 重写 Swap() 方法
	self.Data[i], self.Data[j] = self.Data[j], self.Data[i]
}
func (self *MonitorDataWrapper) Less(i, j int) bool { // 重写 Less() 方法
	return self.Data[i].Time < self.Data[j].Time
}
