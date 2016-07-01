package adxmgtv

import (
	. "common/consts"
	"common/lock"
	"dynamic"
	"github.com/garyburd/redigo/redis"
	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"logic"
)

//for the mgtv, the same user one visitor will cause multiple rtb request ,so we need to removal duplicate return
//return func(a *Activity, index int) (filterd , isContinue bool)
func (b *BidRequest) DuplicateRemovalFn(ctx context.Context, redisConn redis.Conn) func(a *logic.Activity, index int) (filterd bool, isContinue bool) {
	var keyPrefix string
	hasFirstAd := false
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	uids, _ := b.GetUids()
	if len(uids) == 0 {
		keyPrefix = "lock"
	} else {
		//keyPrefix = fmt.Sprintf("%s_%d", uids[0].Value, rand.Intn(5000))
		keyPrefix = uids[0].Value
	}
	fn := func(a *logic.Activity, index int) (filterd bool, isContinue bool) {
		if hasFirstAd {
			return false, false
		}
		clog.Debug("lock_key:%s", keyPrefix+a.ActiveId())
		clog.StartTimer()
		hasLock, err := lock.OnceInXSecond(redisConn, keyPrefix+a.ActiveId(), DuplicateRemovalInterval)
		clog.StopTimer("lock_cost")
		if err != nil {
			if err == lock.ErrorEmptyKey {
				clog.Notice("get uid key empty")
			} else {
				clog.Error("get redis lock err:%s", err)
			}
		}
		if hasLock {
			hasFirstAd = true
			return false, false
		}
		return true, true
	}
	return fn
}

func (b *BidRequest) FrequencyFilterFn(ctx context.Context, redisConn redis.Conn) func(a *logic.Activity, index int) (filterd bool, isContinue bool) {
	key := ""
	uids, _ := b.GetUids()
	if len(uids) > 0 {
		key = uids[0].Value
	}
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	clog.StartTimer()
	freq, err := dynamic.GetFreqByUid(redisConn, key)
	clog.StopTimer("freq_cost")
	if err != nil {
		if err == dynamic.ErrorEmptyKey {
		} else {
			clog.Error("get freq error", err)
		}
	}
	fn := func(a *logic.Activity, index int) (filtered bool, isContinue bool) {
		filtered = false
		isContinue = true
		for _, s := range a.Frequency {
			if freq != nil && !freq.IsUnderFreq(int32(a.Id), s) {
				filtered = true
				break
			}
		}
		return
	}
	return fn
}
