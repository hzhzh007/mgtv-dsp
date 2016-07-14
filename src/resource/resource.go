package resource

import (
	. "common/consts"
	"config"
	"dynamic"
	"errors"
	"github.com/garyburd/redigo/redis"
	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"io/ioutil"
	"logic"
	"net/http"
	"strings"
	"time"
)

type Resource struct {
	ctx context.Context

	//this is a 0-1 buffer type pointer
	activity *logic.Activities

	flow               dynamic.FlowInterface
	activityLocation   string
	activityReloadTime time.Duration
	flowReloaddTime    time.Duration
}

//TODO: do some filter
func (r *Resource) GetActivitiesCopy(request interface{}, onlyValidate bool) (activities *logic.Activities, err error) {
	activities = r.activity
	if activities == nil {
		return nil, errors.New("empty")
	}
	lenA := len(*activities)
	copyActivities := make(logic.Activities, 0, lenA)
	for i := 0; i < lenA; i++ {
		if !(onlyValidate && (*activities)[i].Filtered()) {
			copyActivities = append(copyActivities, (*activities)[i])
		}
	}
	//	copy(copyActivities, *r.activity) how to make it more efficient
	return &copyActivities, nil
}

func (r *Resource) LoadActivities() (*logic.Activities, error) {
	activity := &logic.Activities{}
	if strings.HasPrefix(r.activityLocation, "http") {
		//response, err := http.Get(setUrlQuery(filename, "t", time.Now().Format("0102150405"))) //just because of the bad net cache
		response, err := http.Get(r.activityLocation)
		if err != nil {
			return nil, err
		}
		data, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()
		if err != nil {
			return nil, err
		}
		err = config.LoadConfigFromBytes(data, activity)
		if err != nil {
			return nil, err
		}
	} else {
		err := config.LoadConfig(r.activityLocation, activity)
		if err != nil {
			return nil, err
		}
	}

	//	activity.CommonFilter(r.ctx, r.FlowFilter, "flow")
	//	r.activity = activity
	return activity, nil
}

func (r *Resource) FilterByFlow(activities *logic.Activities) error {
	//	clog.Log.Debug("%v", r.ctx)
	activities.CommonFilter(r.ctx, r.FlowFilter, "flow")
	return nil
}

func (r *Resource) LoadFlow() error {
	redisPool := r.ctx.Value(RedisKey).(*redis.Pool)
	ids := activitiesIds(r.activity)
	redisConn := redisPool.Get()
	defer redisConn.Close()
	flow, err := dynamic.GetActivityFlowData(redisConn, ids, dynamic.FlowImpressionType)
	if err != nil {
		return err
	}
	r.flow = flow
	return nil
}

func activitiesIds(activities *logic.Activities) []int {
	result := make([]int, len(*activities))
	for _, a := range *activities {
		result = append(result, a.Id)
	}
	return result
}

func (r *Resource) Update() {
	flowTimer := time.NewTimer(r.flowReloaddTime)
	activityTimer := time.NewTimer(r.activityReloadTime)
	for {
		select {
		case <-activityTimer.C:
			activityTimer = time.NewTimer(r.activityReloadTime)
			start := time.Now()
			activities, err := r.LoadActivities()
			if err != nil {
				clog.Log.Error("r loadActivities err:%s", err)
				continue
			}
			if err = r.FilterByFlow(activities); err != nil {
				clog.Log.Error("filter Flow error:%s", err)
				continue
			}
			r.activity = activities
			clog.Log.Debug("activity:%+v", r.activity)
			clog.Log.Notice("update actifity  and costs:%s", time.Now().Sub(start))

		case <-flowTimer.C:
			flowTimer = time.NewTimer(r.flowReloaddTime)
			start := time.Now()
			err := r.LoadFlow()
			if err != nil {
				clog.Log.Error("reload flow error:%s", err)
			}
			clog.Log.Notice("update flow and costs:%s", time.Now().Sub(start))
		}
	}
}

func NewResource(ctx context.Context) (res *Resource, err error) {
	conf := ctx.Value(ConfKey).(*config.DspConfig)
	clog := clog.NewContext("module=resource")
	ctx = context.WithValue(ctx, ContextLogKey, clog)
	res = &Resource{
		activityLocation:   conf.Resource.Activity.Location,
		activityReloadTime: conf.Resource.Activity.Reload,
		flowReloaddTime:    conf.Resource.Flow.Reload,
		ctx:                ctx,
	}
	activities, err := res.LoadActivities()
	if err != nil {
		return nil, err
	}
	res.activity = activities //just put it temp
	if err = res.LoadFlow(); err != nil {
		return nil, err
	}
	if err = res.FilterByFlow(activities); err != nil {
		return nil, err
	}
	res.activity = activities
	go res.Update()
	return res, nil

}

func (r *Resource) FlowFilter(a *logic.Activity, index int) (isFilterd, isContinue bool) {
	isFilterd = false
	isContinue = true
	todayFlow := r.flow.TodayFlow(a.Id)
	totalFlow := r.flow.TotalFlow(a.Id)
	if !a.IsFlowOK(todayFlow, totalFlow) {
		isFilterd = true
	}
	return
}
