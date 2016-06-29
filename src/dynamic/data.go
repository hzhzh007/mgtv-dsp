package dynamic

import (
	"github.com/jinzhu/now"
	"logic"
	"time"
)

const (
	WeekCycleSeconds  = 60 * 60 * 24 * 7
	MonthCycleSeconds = 60 * 60 * 24 * 30
)

func (r *Record) SetExpire(a *logic.Activity, nowTime time.Time) {
	switch r.Type {
	case FreqType_FreqPerDay:
		r.Expire = int32(now.New(nowTime).EndOfDay().Unix())

	case FreqType_FreqPerWeek:
		r.Expire = r.ExpireEnd(int32(a.CreativeStart().Unix()), int32(nowTime.Unix()), WeekCycleSeconds)

	case FreqType_FreqPerMonth:
		r.Expire = r.ExpireEnd(int32(a.CreativeStart().Unix()), int32(nowTime.Unix()), MonthCycleSeconds)

	case FreqType_FreqCustom:
		//should not run to here
	}
}

//TODO: the algrithm to be reviewed
func (r *Record) ExpireEnd(start, now, cycle int32) int32 {
	return (((now-start)/cycle)+1)*cycle + start
}
