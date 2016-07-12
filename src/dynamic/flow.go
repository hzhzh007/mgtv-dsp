package dynamic

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	clog "github.com/hzhzh007/context_log"
	"time"
)

type FlowResult struct {
	totalFlow map[int]int
	todayFlow map[int]int
	cacheTime time.Time
}

func (f *FlowResult) TodayFlow(activityId int) int {
	num, ok := f.todayFlow[activityId]
	if ok {
		return num
	}
	return 0
}

func (f *FlowResult) TotalFlow(id int) int {
	num, ok := f.totalFlow[id]
	if ok {
		return num
	}
	return 0
}

func NewFlowResult(redisConn redis.Conn, ids []int, flowType string) (*FlowResult, error) {
	fr := &FlowResult{
		cacheTime: time.Now(),
		totalFlow: make(map[int]int),
		todayFlow: make(map[int]int),
	}
	err := fr.GetDataOneByOne(redisConn, ids, flowType)
	return fr, err
}

func (f *FlowResult) GetDataOneByOne(redisConn redis.Conn, ids []int, flowType string) error {
	err := f.GetDataTotal(redisConn, ids, flowType)
	if err != nil {
		return err
	}
	err = f.GetDataToday(redisConn, ids, flowType)
	return err

}

//
func (f *FlowResult) GetDataTotal(redisConn redis.Conn, ids []int, flowType string) error {
	for _, id := range ids {
		flow, err := redis.Int(redisConn.Do("GET", KeyTotal(id, flowType)))
		if err == redis.ErrNil {
			continue
		} else if err != nil {
			clog.Log.Error("get flow from redis error:%s", err)
		} else {
			f.totalFlow[id] = flow
		}
	}
	return nil
}

func (f *FlowResult) GetDataToday(redisConn redis.Conn, ids []int, flowType string) error {
	today := time.Now().Format("20060102")
	for _, id := range ids {
		flow, err := redis.Int(redisConn.Do("GET", KeyToday(id, flowType, today)))
		if err == redis.ErrNil {
			continue
		} else if err != nil {
			clog.Log.Error("get flow from redis error:%s", err)
		} else {
			f.todayFlow[id] = flow
		}
	}
	return nil
}

func KeyTotal(id int, flowType string) string {
	return fmt.Sprintf("flow%s_%d_total", flowType, id)
}

func KeyToday(id int, flowType, today string) string {
	return fmt.Sprintf("flow%s_%d_%s", flowType, id, today)
}

func IncrActivityFlow(redisConn redis.Conn, id int, flowType string) error {
	_, err := redis.Int(redisConn.Do("INCR", KeyTotal(id, flowType)))
	if err != nil {
		return err
	}
	_, err = redis.Int(redisConn.Do("INCR", KeyToday(id, flowType, time.Now().Format("20060102"))))
	return err
}
