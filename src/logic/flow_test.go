package logic

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_FlowNumInfiniteType(t *testing.T) {
	assert := assert.New(t)
	f := FlowStratage{Type: FlowNumInfiniteType}
	assert.True(f.UnderFlow(0, 0, Schedules{}, time.Now()))
	assert.True(f.UnderFlow(999999, 999999, Schedules{}, time.Now()))
}

func Test_FlowNumTotalTypeAndSpeedAvgType(t *testing.T) {
	assert := assert.New(t)
	f := FlowStratage{Type: FlowNumTotalType, Speed: SpeedAvgType, Num: 1000}
	schedule := Schedule{}
	schedule.Start, _ = time.Parse("2006-01-02", "2016-06-15")
	schedule.End, _ = time.Parse("2006-01-02", "2016-06-20")
	now, _ := time.Parse("2006-01-02 15:04:05", "2016-06-16 12:30:00")
	assert.True(f.UnderFlow(0, 0, Schedules{schedule}, now))
	assert.False(f.UnderFlow(1000, 1000, Schedules{schedule}, now))
	assert.True(f.UnderFlow(50, 250, Schedules{schedule}, now))
	assert.True(f.UnderFlow(150, 250, Schedules{schedule}, now))
	assert.False(f.UnderFlow(60, 350, Schedules{schedule}, now))
	assert.False(f.UnderFlow(60, 310, Schedules{schedule}, now))
	assert.True(f.UnderFlow(160, 301, Schedules{schedule}, now))
}
func Test_FlowNumTotalTypeAndSpeedASAP(t *testing.T) {
	assert := assert.New(t)
	f := FlowStratage{Type: FlowNumTotalType, Speed: SpeedAsSoonAsPossibleType, Num: 1000}
	schedule := Schedule{}
	schedule.Start, _ = time.Parse("2006-01-02", "2016-06-15")
	schedule.End, _ = time.Parse("2006-01-02", "2016-06-20")
	now, _ := time.Parse("2006-01-02 15:04:05", "2016-06-16 12:30:00")
	assert.True(f.UnderFlow(0, 0, Schedules{schedule}, now))
	assert.False(f.UnderFlow(1000, 1000, Schedules{schedule}, now))
	assert.True(f.UnderFlow(50, 250, Schedules{schedule}, now))
	assert.True(f.UnderFlow(150, 250, Schedules{schedule}, now))
	assert.True(f.UnderFlow(60, 350, Schedules{schedule}, now))
	assert.True(f.UnderFlow(60, 310, Schedules{schedule}, now))
	assert.False(f.UnderFlow(160, 1000, Schedules{schedule}, now))
}

func Test_FlowNumPerDayTypeAndSpeedAvgType(t *testing.T) {
	assert := assert.New(t)
	f := FlowStratage{Type: FlowNumPerDayType, Speed: SpeedAvgType, Num: 1000}
	schedule := Schedule{}
	schedule.Start, _ = time.Parse("2006-01-02", "2016-06-15")
	schedule.End, _ = time.Parse("2006-01-02", "2016-06-20")
	now, _ := time.Parse("2006-01-02 15:04:05", "2016-06-16 12:30:00")
	assert.True(f.UnderFlow(0, 0, Schedules{schedule}, now))
	assert.False(f.UnderFlow(1000, 1000, Schedules{schedule}, now))
	assert.True(f.UnderFlow(150, 250, Schedules{schedule}, now))
	assert.True(f.UnderFlow(250, 250, Schedules{schedule}, now))
	assert.True(f.UnderFlow(460, 350, Schedules{schedule}, now))
	assert.False(f.UnderFlow(600, 310, Schedules{schedule}, now))
	assert.False(f.UnderFlow(560, 310, Schedules{schedule}, now))
	assert.False(f.UnderFlow(700, 310, Schedules{schedule}, now))
	assert.True(f.UnderFlow(501, 301, Schedules{schedule}, now))
}
