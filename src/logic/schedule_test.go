package logic

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_TimeUnder(t *testing.T) {
	assert := assert.New(t)
	schedule := Schedule{}
	schedule.Start, _ = time.Parse("2006-01-02", "2016-06-15")
	schedule.End, _ = time.Parse("2006-01-02", "2016-06-20")
	time0615, _ := time.Parse("2006-01-02", "2016-06-15")
	time0620, _ := time.Parse("2006-01-02", "2016-06-20")
	time0618, _ := time.Parse("2006-01-02", "2016-06-18")
	time0701, _ := time.Parse("2006-01-02", "2016-07-01")
	time0501, _ := time.Parse("2006-01-02", "2016-05-01")

	assert.True(schedule.UnderSchedule(time0615))
	assert.True(schedule.UnderSchedule(time0620))
	assert.True(schedule.UnderSchedule(time0618))
	t.Log(schedule, time0701)
	assert.False(schedule.UnderSchedule(time0701))
	assert.False(schedule.UnderSchedule(time0501))

}

func Test_TotalDuration(t *testing.T) {
	assert := assert.New(t)
	schedule := Schedule{}
	schedule.Start, _ = time.Parse("2006-01-02", "2016-06-15")
	schedule.End, _ = time.Parse("2006-01-02", "2016-06-20")
	schedules := Schedules{schedule}
	totalDuration := schedules.TotalDuration()
	assert.Equal(24*time.Hour*5, totalDuration)
}
