package utils

import (
	"time"
)

// 按给定的刷新时间获取当天刷新时间点，用于数据库查询
func GetDailyByRefreshTime(hour int64) string {

	now := time.Now()
	daily := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)

	if hour < 0 || hour > 23 { // 超出范围，返回当天0点
		return daily.Format(time.DateTime)
	}

	daily = daily.Add(time.Duration(hour) * time.Hour)

	return daily.Format(time.DateTime)
}
