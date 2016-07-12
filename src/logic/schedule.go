package logic

import (
	"time"
)

type Schedules []Schedule

type Schedule struct {
	Start time.Time
	End   time.Time
}

func (s Schedule) UnderSchedule(t time.Time) bool {
	return (t.After(s.Start) && t.Before(s.End)) || t.Equal(s.Start) || t.Equal(s.End)
}

func (ses Schedules) TotalDuration() time.Duration {
	var d time.Duration = 0
	for i := 0; i < len(ses); i++ {
		d += ses[i].End.Sub(ses[i].Start)
	}
	return d
}

func (ses Schedules) PassedDuration(now time.Time) time.Duration {
	var d time.Duration = 0
	for i := 0; i < len(ses); i++ {
		if now.After(ses[i].End) {
			d += ses[i].End.Sub(ses[i].Start)
		} else if now.After(ses[i].Start) {
			d += now.Sub(ses[i].Start)
		}
	}
	return d

}
