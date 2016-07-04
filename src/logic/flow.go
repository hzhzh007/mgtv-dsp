package logic

import (
	clog "github.com/hzhzh007/context_log"
	"github.com/jinzhu/now"
	"time"
)

const (
	FlowNumInfiniteType       = 0 //无限量投放
	FlowNumTotalType          = 1 //总共投
	FlowNumPerDayType         = 2 //每天投
	SpeedAvgType              = 0 //投放速度
	SpeedAsSoonAsPossibleType = 1
	DaySeconds                = 86400
)

type FlowNumType int
type FlowType []FlowStratage

type FlowStratage struct {
	Type  FlowNumType
	Speed int
	Num   int
}

func (f *FlowStratage) UnderFlow(today, total int, schedules Schedules, nowTime time.Time) bool {
	clog.Log.Debug("today:%d, total:%d, schdule:%+v, flowStartage:%+v", today, total, schedules, *f)
	switch f.Type {
	case FlowNumInfiniteType:
		return true
	case FlowNumTotalType:
		return (total < f.Num) && f.UnderSpeedTotal(today, total, schedules, nowTime)
	case FlowNumPerDayType:
		return (today < f.Num) && f.UnderSpeedPerDay(today, total, schedules, nowTime)
	}
	return false
}

func (f *FlowStratage) UnderSpeedPerDay(today, total int, schedules Schedules, nowTime time.Time) bool {
	switch f.Speed {
	case SpeedAsSoonAsPossibleType:
		return true
	case SpeedAvgType:
		todayStart := now.New(nowTime).BeginningOfDay()
		return today <= ((int(nowTime.Sub(todayStart)) * f.Num) / int(time.Hour*24))
	}
	return false
}

func (f *FlowStratage) UnderSpeedTotal(today, total int, schedules Schedules, nowTime time.Time) bool {
	switch f.Speed {
	case SpeedAsSoonAsPossibleType:
		return true
	case SpeedAvgType:
		return total <= (int(schedules.PassedDuration(nowTime)) * f.Num / int(schedules.TotalDuration()))
	}
	return false
}
