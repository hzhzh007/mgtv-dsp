//flow and frequency data store module
//TODO: here we just one by one insert, but should use batch way
package feedback

import (
	. "common/consts"
	"dynamic"
	"errors"
	"github.com/garyburd/redigo/redis"
	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"logic"
	"resource"
	"time"
)

const (
	RecordClickType      = 0
	RecordImpressionType = 1
)

type Record struct {
	Uid string
	//	UidType  string
	Type     int
	Activity int
	//	Creative int
}

type Feedback struct {
	ctx        context.Context
	Input      chan Record
	activities *logic.Activities
}

func NewFeedback(ctx context.Context, chanSize, workerNum int) (*Feedback, error) {
	clog := clog.NewContext("module=update")
	ctx = context.WithValue(ctx, ContextLogKey, clog)

	feedback := &Feedback{
		ctx:   ctx,
		Input: make(chan Record, chanSize),
	}
	go feedback.UpdateActivities()
	feedback.DeriveWorker(workerNum)
	return feedback, nil
}

func (f *Feedback) DeriveWorker(workerNum int) {
	for i := 0; i < workerNum; i++ {
		go f.Worker()
	}
}

func (f *Feedback) UpdateActivities() {
	resource := f.ctx.Value(ResourceKey).(*resource.Resource)
	for {
		activities, err := resource.GetActivitiesCopy(nil, false)
		if err != nil {
			clog.Log.Error("copy activities err:%s", err)
		} else {
			f.activities = activities
		}
		time.Sleep(time.Second * 60)
	}
}

func (f *Feedback) GetActivityById(id int) (*logic.Activity, error) {
	if f.activities == nil {
		return nil, errors.New("self activities is nil")
	}
	return f.activities.GetActivityById(id)
}

func (f *Feedback) incFrequency(redisConn redis.Conn, record Record, activity *logic.Activity) error {
	if len(activity.Frequency) == 0 {
		return nil
	}
	freq, err := dynamic.GetFreqByUid(redisConn, record.Uid)
	if err != nil {
		return err
	}
	if len(activity.Frequency) == 0 {
		clog.Log.Debug("freq is nil ,just pass, activity:%d", activity.Id)
		return nil
	}
	for _, freqSetting := range activity.Frequency {
		clog.Log.Debug("freqSetting:%+v", freqSetting)
		freq.IncActivityFreqByOneSetting(freqSetting, activity)
	}
	return freq.Persister(redisConn)
}

func (f *Feedback) Worker() {
	redisPool := f.ctx.Value(RedisKey).(*redis.Pool)
	for record := range f.Input {
		clog.Log.Debug("received a record:%v", record)
		redisConn := redisPool.Get()
		activity, err := f.GetActivityById(record.Activity)
		if err != nil {
			clog.Log.Error("get activity by id err:%s", err)
			continue
		}

		if len(activity.Frequency) > 0 {
			err := f.incFrequency(redisConn, record, activity)
			if err != nil {
				clog.Log.Error("inc freq error:%s", err)
			}
		}
		dynamic.IncrActivityFlow(redisConn, activity.Id, dynamic.FlowImpressionType)
	}
}
