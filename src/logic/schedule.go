package logic

import (
	"time"
)

type Schedule struct {
	Start time.Time
	End   time.Time
}

func (s Schedule) UnderSchedule(t time.Time) bool {
	return (t.After(s.Start) && t.Before(s.End)) || t.Equal(s.Start) || t.Equal(s.End)
}
