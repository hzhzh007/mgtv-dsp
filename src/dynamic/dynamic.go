//user and activity state relate service
//now include flow and frequency stratage
package dynamic

import (
	"github.com/garyburd/redigo/redis"
	"logic"
	//	"time"
)

const (
	FlowImpressionType = "i"
	FlowClickType      = "c"
)

type FrequencyInterface interface {
	//Uid() string
	IsUnderFreq(activityId int32, setting logic.FrequencyStratage) bool
}

type FlowInterface interface {
	//	DataTime() time.Time
	TodayFlow(activityId int) int
	TotalFlow(activityId int) int
}

//func GetFreqByUid(redisConn redis.Conn, uid string) (FrequencyInterface, error) {
func GetFreqByUid(redisConn redis.Conn, uid string) (*FreqValue, error) {
	return NewFreqFromRedis(redisConn, uid)
}

func GetActivityFlowData(redisConn redis.Conn, activityIds []int, flowType string) (FlowInterface, error) {
	return NewFlowResult(redisConn, activityIds, flowType)
}
