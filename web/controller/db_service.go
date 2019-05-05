/*
 * @Author: rayou
 * @Date: 2019-04-30 12:11:16
 * @Last Modified by: rayou
 * @Last Modified time: 2019-05-01 12:15:21
 */
package controller

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AzuresYang/arx7/app/arxmonitor"
	"github.com/AzuresYang/arx7/config"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type dbService struct {
	db_cfg *config.MysqlConfig
	db     *sql.DB
}

var DbService *dbService = &dbService{}

func (self *dbService) Init(db_cfg *config.MysqlConfig) error {
	self.db_cfg = db_cfg
	path := strings.Join([]string{db_cfg.UserName, ":", db_cfg.Password, "@tcp(", db_cfg.Ip, ":",
		string(db_cfg.Port), ")/", db_cfg.DbName, "?charset=", db_cfg.Charset}, "")
	var err error
	self.db, err = sql.Open("mysql", path)
	if err != nil {
		log.Errorf("[DbService.Init] Open mysql error:%s", err.Error())
		return err
	}
	fmt.Println("start mysql")
	// log.Debugf("[DbService.Init]config :%+v", db_cfg)
	// 设置链接最大空闲时间
	self.db.SetConnMaxLifetime(100 * time.Second)
	// 设置上数据库最大闲置连接数
	self.db.SetMaxIdleConns(db_cfg.MaxIdleConns)
	return nil
}

func (self *dbService) GetMonitorInfo(query *FormQueryMonitor, start_time int64, end_time int64) ([]MonitorData, error) {
	code_info := "DbService.GetDbInfo"
	datas := make([]MonitorData, 0)
	var rows *sql.Rows
	var err error
	fmt.Printf("start_time:%d, end_time:%d\n", start_time, end_time)
	if query.Classfy == "0" {
		rows, err = self.db.Query("select * from monitor_data where svcid=? and metric=? and time >=? and time <=? order by time;",
			query.SvcId, query.Metric, start_time, end_time)
	} else {
		rows, err = self.db.Query("select * from monitor_data where svcid=? and metric=? and classfy=? and time >=? and time <=? order by time;",
			query.SvcId, query.Metric, query.Classfy, start_time, end_time)
	}

	if err != nil {
		log.Errorf("[%s]query error:%s", code_info, err.Error())
		return datas, errors.New("服务内部错误")
	}

	for rows.Next() {
		row := MonitorData{}
		serr := rows.Scan(&row.Id, &row.Svcid, &row.Metric, &row.Classfy, &row.Value, &row.Ip, &row.Time)
		if serr != nil {
			log.Errorf("[%s] scan db data fail:%s", code_info, serr.Error())
		} else {
			datas = append(datas, row)
		}
	}
	return generateMonitorData(start_time, end_time, datas), nil
}

func (self *dbService) GetMonitorInfoByForm(query *FormQueryMonitor) ([]MonitorData, error) {
	time_layout := "2006-01-02 15:04:05" //转化所需模板
	loc, _ := time.LoadLocation("Local")
	temp_start_time, _ := time.ParseInLocation(time_layout, query.StartTime, loc) //使用模板在对应时区转化为time.time类型
	temp_end_time, _ := time.ParseInLocation(time_layout, query.EndTime, loc)     //使用模板在对应时区转化为time.time类型
	start_time := temp_start_time.Unix()
	end_time := temp_end_time.Unix()
	return self.GetMonitorInfo(query, start_time, end_time)
}

// 需要按照精度对时间生成
func generateMonitorData(start_time int64, end_time int64, monitor_data []MonitorData) []MonitorData {
	start_time -= start_time % arxmonitor.DEFAULT_MONITOR_TIME_INTERVAL
	end_time += arxmonitor.DEFAULT_MONITOR_TIME_INTERVAL
	end_time -= end_time % arxmonitor.DEFAULT_MONITOR_TIME_INTERVAL
	capacity := (end_time - start_time) / arxmonitor.DEFAULT_MONITOR_TIME_INTERVAL
	data := make([]MonitorData, 0, capacity)
	// 没有监控数据，只需要填充时间列就好了
	temp := MonitorData{}
	var i int64 = 0
	if len(monitor_data) <= 0 {
		for i = 0; i < capacity; i++ {
			temp.Time = start_time + i*arxmonitor.DEFAULT_MONITOR_TIME_INTERVAL
			data = append(data, temp)
		}
	} else {
		// 有监控数据,对监控数据的格式控制一下。

		monitor_idx := 0

		for i = 0; i < capacity; i++ {
			temp.Time = start_time + i*arxmonitor.DEFAULT_MONITOR_TIME_INTERVAL
			temp.Value = 0

			// 没有数据就直接结束了吧
			if monitor_idx >= len(monitor_data) {
				break
			}
			// 精度不对，填充进去，原先的精度需要补充
			if temp.Time > monitor_data[monitor_idx].Time {
				data = append(data, monitor_data[monitor_idx])
				monitor_idx++
			} else if temp.Time == monitor_data[monitor_idx].Time {
				// 精度相同，值替换一下
				temp.Value = monitor_data[monitor_idx].Value
				monitor_idx++
			}
			data = append(data, temp)
		}
	}
	// 需要对时间做统一化处理就打开这个，不过会耗用一些性能
	data = correctMonitorData(data)
	return data
}

func correctMonitorData(monitor_data []MonitorData) []MonitorData {
	data := make([]MonitorData, 0, len(monitor_data))
	data_idx := -1 // 最后一位
	for i, _ := range monitor_data {
		if i == 0 {
			data = append(data, monitor_data[i])
			data_idx++
			continue
		}
		// 时间统一处理
		monitor_data[i].Time -= monitor_data[i].Time % arxmonitor.DEFAULT_MONITOR_TIME_INTERVAL
		// fmt.Printf("i:%d, monitor:%d, idx:%d, data:%d\n", i, data_idx, len(data_idx))
		if monitor_data[i-1].Time == data[data_idx].Time {
			data[data_idx].Value += monitor_data[i].Value
		} else {
			data = append(data, monitor_data[i])
			data_idx++
		}
	}
	return data
}
