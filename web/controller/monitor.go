/*
 * @Author: rayou
 * @Date: 2019-04-27 16:01:21
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-04 17:35:34
 */

package controller

import (
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	default_real_monitor_interval int64 = 10 * 60          // 10 分钟内的
	default_max_query_start2end   int64 = 60 * 60 * 24 * 2 // 最多查询2天的数据
)

type MonitorInfo struct {
	XAxsis []string `json:"xAxsis"`
	Series []uint32 `json:"series"`
}

func MonitorInfoHandler(response http.ResponseWriter, request *http.Request) {
	code_info := "MonitorInfoHandler"
	log.Info("Get MonitorInfo")
	form_monitor := &FormQueryMonitor{}
	err := parseForm(request, form_monitor)
	if err != nil {
		log.Errorf("[MonitorInfohandler] parse info fail:%+v", request.Form)
		responseJson(response, 1, "表单参数不能为空", "")
		return
	}
	log.Infof("Get Monitor Form:%+v", form_monitor)
	err = checkQueryMonitorInfo(form_monitor)
	if err != nil {
		log.Errorf("[%s]err:%s, Form:%+v", code_info, err.Error(), form_monitor)
		responseJson(response, 2, err.Error(), "")
		return
	}
	// 查询
	monitor_data, qerr := DbService.GetMonitorInfoByForm(form_monitor)
	if qerr != nil {
		responseJson(response, 3, qerr.Error(), "")
	}

	monitor_info := MonitorInfo{
		XAxsis: make([]string, 0, len(monitor_data)),
		Series: make([]uint32, 0, len(monitor_data)),
	}
	timeLayout := "2006-01-02 15:04:05"
	for i, _ := range monitor_data {
		monitor_info.XAxsis = append(monitor_info.XAxsis, time.Unix(monitor_data[i].Time, 0).Format(timeLayout))
		monitor_info.Series = append(monitor_info.Series, uint32(monitor_data[i].Value))
	}
	log.Infof("[%s]get monitor info succ.len%d", code_info, len(monitor_info.Series))
	err = responseJson(response, 0, "succ", monitor_info)
	if err != nil {
		log.Errorf("[monitorInfoHandler] response err:%s", err.Error())
	}
}

func checkQueryMonitorInfo(form *FormQueryMonitor) error {
	// toBeCharge := "2019-04-30 00:00:00"                             //待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	time_layout := "2006-01-02 15:04:05"                                            //转化所需模板
	loc, _ := time.LoadLocation("Local")                                            //重要：获取时区
	temp_start_time, serr := time.ParseInLocation(time_layout, form.StartTime, loc) //使用模板在对应时区转化为time.time类型
	temp_end_time, eerr := time.ParseInLocation(time_layout, form.EndTime, loc)     //使用模板在对应时区转化为time.time类型
	if serr != nil || eerr != nil {
		return errors.New("时间参数格式不对")
	}
	start_time := temp_start_time.Unix()
	end_time := temp_end_time.Unix()
	if start_time >= end_time {
		return errors.New("开始时间不能大于截止时间")
	}
	if (end_time - start_time) > default_max_query_start2end {
		return errors.New("查询时间程度不能大于2天。")
	}
	return nil
}
