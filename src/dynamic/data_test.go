package dynamic

import (
	"github.com/stretchr/testify/assert"
	"logic"
	"testing"
	"time"
)

func Test_RecordExpirePerDay(t *testing.T) {
	assert := assert.New(t)
	r := Record{
		Type: FreqType_FreqPerDay,
	}
	nowTime, _ := time.Parse("2006-01-02 15:04:05", "2016-06-23 20:14:00")
	t1, _ := time.Parse("2006-01-02 15:04:05", "2016-06-23 23:59:55")
	t2, _ := time.Parse("2006-01-02 15:04:05", "2016-06-24 00:00:02")
	r.SetExpire(nil, nowTime)
	assert.True(time.Unix(int64(r.Expire), 0).After((t1)))
	assert.True(time.Unix(int64(r.Expire), 0).Before((t2)))
}

func Test_RecordExpirePerWeek(t *testing.T) {
	assert := assert.New(t)
	r := Record{
		Type: FreqType_FreqPerWeek,
	}
	t0, _ := time.Parse("2006-01-02 15:04:05", "2016-06-21 12:00:00")
	a := &logic.Activity{
		ActiveTime: []logic.Schedule{logic.Schedule{Start: t0}},
	}
	nowTime, _ := time.Parse("2006-01-02 15:04:05", "2016-06-23 20:14:00")
	t1, _ := time.Parse("2006-01-02 15:04:05", "2016-06-28 11:59:55")
	t2, _ := time.Parse("2006-01-02 15:04:05", "2016-06-28 12:00:02")
	r.SetExpire(a, nowTime)
	assert.True(time.Unix(int64(r.Expire), 0).After((t1)))
	t.Log(time.Unix(int64(r.Expire), 0))
	assert.True(time.Unix(int64(r.Expire), 0).Before((t2)))
}

func Test_RecordExpirePerMonth(t *testing.T) {
	assert := assert.New(t)
	r := Record{
		Type: FreqType_FreqPerMonth,
	}
	t0, _ := time.Parse("2006-01-02 15:04:05", "2016-06-21 12:00:00")
	a := &logic.Activity{
		ActiveTime: []logic.Schedule{logic.Schedule{Start: t0}},
	}
	nowTime, _ := time.Parse("2006-01-02 15:04:05", "2016-06-23 20:14:00")
	t1, _ := time.Parse("2006-01-02 15:04:05", "2016-07-21 11:59:55")
	t2, _ := time.Parse("2006-01-02 15:04:05", "2016-07-22 12:00:02")
	r.SetExpire(a, nowTime)
	assert.True(time.Unix(int64(r.Expire), 0).After((t1)))
	t.Log(time.Unix(int64(r.Expire), 0))
	assert.True(time.Unix(int64(r.Expire), 0).Before((t2)))
}
