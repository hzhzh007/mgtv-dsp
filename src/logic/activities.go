package logic

import (
	. "common/consts"
	"errors"
	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"sort"
	"time"
)

//type Activities []Activity

//TODO here Activity may use pointer to have more performance
type Activities []Activity

/*
type Activities struct {
	Row Activities
}*/

//common user extend function filter
//receive the func(a *Activity, index int) (filterd , isContinue bool)
func (activities Activities) CommonFilter(ctx context.Context, fn func(a *Activity, index int) (bool, bool), filterType string) {
	var (
		filtered   bool
		isContinue bool
	)
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	for i := 0; i < len(activities); i++ {
		if activities[i].Filtered() {
			continue
		}
		if filtered, isContinue = fn(&activities[i], i); filtered == true {
			activities[i].SetFiltered()
			clog.Debug("activity: %s filtered by %s", activities[i].ActiveId(), filterType)
		}
		if isContinue == false {
			break
		}
	}
}

func (activities Activities) LocationFilter(ctx context.Context, location Location) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	for i := 0; i < len(activities); i++ {
		if activities[i].Filtered() {
			continue
		}
		if !activities[i].LocationOK(location) {
			activities[i].SetFiltered()
			clog.Debug("%s filtered by location", activities[i].ActiveId())
		}
	}
}

func (activities Activities) ScheduleFilter(ctx context.Context, now time.Time) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	for i := 0; i < len(activities); i++ {
		if activities[i].Filtered() {
			continue
		}
		if !activities[i].ScheduleOK(now) {
			activities[i].SetFiltered()
			clog.Debug("%s filtered by schedule", activities[i].ActiveId())
		}
	}
}

func (activities Activities) PlatformFilter(ctx context.Context, platform Platform) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	for i := 0; i < len(activities); i++ {
		if activities[i].Filtered() {
			continue
		}
		if !activities[i].PlatformOK(platform) {
			activities[i].SetFiltered()
			clog.Debug("%+v filtered by platform", activities[i])
		}
	}
}

func (activities Activities) UserTagFilter(ctx context.Context, userTags Tags) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	for i := 0; i < len(activities); i++ {
		if activities[i].Filtered() {
			continue
		}
		if !activities[i].TagOK(userTags) {
			activities[i].SetFiltered()
			clog.Debug("%s filtered by tag", activities[i].ActiveId())
		}
	}
}

func (activities Activities) DurationFilter(ctx context.Context, duration int) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	for i := 0; i < len(activities); i++ {
		if activities[i].Filtered() {
			continue
		}
		if !activities[i].Length.UnderDuration(duration) {
			activities[i].SetFiltered()
			clog.Debug("%s filtered by duration", activities[i].ActiveId())
		}
	}
}

func (activities Activities) AdTypeFilter(ctx context.Context, adType int) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	for i := 0; i < len(activities); i++ {
		if activities[i].Filtered() {
			continue
		}
		if !(activities[i].AdType == adType) {
			activities[i].SetFiltered()
			clog.Debug("%s filtered by adtype a:%d, r:%d", activities[i].ActiveId(), activities[i].AdType, adType)
		}
	}
}

//TODO: use binary search
func (activities Activities) GetActivityById(id int) (*Activity, error) {
	for i := 0; i < len(activities); i++ {
		if activities[i].Id == id {
			return &activities[i], nil
		}
	}
	return nil, errors.New("not found activity, id")
}

func (activities Activities) Sort() {
	sort.Sort(activities)
}

func (r Activities) Less(i, j int) bool {
	if r[i].Filtered() {
		return false
	}
	if r[j].Filtered() {
		return true
	}
	return r[i].GetECPM() > r[j].GetECPM()
}
func (r Activities) Swap(i, j int) {
	temp := r[j]
	r[j] = r[i]
	r[i] = temp
}

//TODO
func (r Activities) Len() int {
	return len(r)
}

//ATTENTION: must call after Sort()
func (r Activities) AvailableLen() int {
	for i := 0; i < len(r); i++ {
		if r[i].Filtered() {
			return i
		}
	}
	return len(r)
}
