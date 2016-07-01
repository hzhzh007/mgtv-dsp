package dynamic

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	proto "github.com/golang/protobuf/proto"
	clog "github.com/hzhzh007/context_log"
	"logic"
	"time"
)

const (
	RedisFrePrifix = "fr_"
)

type FreqValue struct {
	RedisValue
	Uid string
}

var (
	ErrorEmptyKey = errors.New("empty key and return null")
)

func NewFreqFromRedis(redisConn redis.Conn, key string) (*FreqValue, error) {
	if key == "" {
		return nil, ErrorEmptyKey
	}
	redisValue := &RedisValue{}
	data, err := redis.Bytes(redisConn.Do("GET", RedisFrePrifix+key))
	if err == redis.ErrNil {
		return &FreqValue{*redisValue, key}, nil
	} else if err != nil {
		return nil, err
	}
	err = proto.Unmarshal(data, redisValue)
	if err != nil {
		return nil, err
	}
	return &FreqValue{*redisValue, key}, nil
}

func (r *FreqValue) Persister(redisConn redis.Conn) error {
	data, err := proto.Marshal(&r.RedisValue)
	if len(data) == 0 {
		errors.New("zero data")
	}
	_, err = redis.String(redisConn.Do("SET", RedisFrePrifix+r.Uid, data, "EX", r.getMaxExpire()))
	return err
}

func (r *FreqValue) IsUnderFreq(activityId int32, setting logic.FrequencyStratage) bool {
	clog.Log.Debug("freq filter id:%d, setting:%+v", activityId, setting)
	for _, record := range r.Impression {
		clog.Log.Debug("record:%+v", record)
		if record.Id != activityId || int32(record.Type) != int32(setting.GetType()) {
			clog.Log.Debug("not equal")
			continue
		}
		if record.Counter >= int32(setting.GetNum()) {
			clog.Log.Debug("value less")
			return false
		}
	}
	return true
}

func (f *FreqValue) IncActivityFreqByOneSetting(freqSetting logic.FrequencyStratage, a *logic.Activity) {
	now := time.Now().Unix()
	id := int32(a.Id)
	for i := 0; i < len(f.Impression); {
		if f.Impression[i].Expire < int32(now) {
			f.Impression = append(f.Impression[:i], f.Impression[i+1:]...)
		}
		if id == f.Impression[i].Id &&
			int32(f.Impression[i].Type) == int32(freqSetting.Type) {
			f.Impression[i].Counter++
			clog.Log.Debug("++ return :%+v", f.Impression)
			return
		}
		i++
	}
	f.Impression = append(f.Impression, NewRecord(freqSetting, a))
	clog.Log.Debug("last return :%+v", f.Impression)
}

func NewRecord(freqSetting logic.FrequencyStratage, a *logic.Activity) *Record {
	activityId := int32(a.Id)
	r := &Record{Id: activityId,
		Expire:  0,
		Type:    FreqType(freqSetting.Type),
		Counter: 1,
	}
	r.SetExpire(a, time.Now())
	return r
}

func (f *FreqValue) getMaxExpire() (maxExpire int32) {
	maxExpire = 0
	for _, r := range f.Impression {
		if r.Expire > maxExpire {
			maxExpire = r.Expire
		}
	}
	return maxExpire
}
