//adapter of request and dsp logic filter
package adxmgtv

import (
	. "common/consts"
	"common/lock"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	clog "github.com/hzhzh007/context_log"
	proto "github.com/hzhzh007/mgtvPmpProto"
	"golang.org/x/net/context"
	"io/ioutil"
	"iplib"
	"logic"
	"math/rand"
	"net/http"
	"time"
)

const (
	ProtocolVersion          = 3
	DuplicateRemovalInterval = 2 // 2s
	WinNoticeUrl             = "http://wy.mgtv.com/winnotice?c=%%SETTLE_PRICE%%"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type BidRequest struct {
	proto.Request
}

//TODO: here iplib do not handle err
func (b *BidRequest) GetLocation(ctx context.Context) (logic.Location, error) {
	iplib := ctx.Value(IPLibKey).(iplib.IpCollections)
	ipStr := b.Device.Ip
	city := iplib.Ip2CityCode(ipStr)
	return logic.Location(city), nil
}

//here we use the mgtv platform type, so we do not need to translate it
func (b *BidRequest) GetPlatform(ctx context.Context) (logic.Platform, error) {
	return logic.Platform(b.Device.Type), nil
}

func (b *BidRequest) GetUids() (logic.Uids, error) {
	uids := make(logic.Uids, 0, 3)
	//android
	uids.AddId(logic.UidImeiSha1Type, b.Device.Imei)
	uids.AddId(logic.UidAnidType, b.Device.Anid)

	//ios
	uids.AddId(logic.UidOpenudidType, b.Device.Openudid)
	uids.AddId(logic.UidIdfaType, b.Device.Idfa)

	uids.AddId(logic.UidMacSha1Type, b.Device.Mac)
	uids.AddId(logic.UidCookieType, b.Device.Duid)
	return uids, nil
}

//for the mgtv, the same user one visitor will cause multiple rtb request ,so we need to removal duplicate return
//return func(a *Activity, index int) (filterd , isContinue bool)
func (b *BidRequest) DuplicateRemovalFn(ctx context.Context) func(a *logic.Activity, index int) (filterd bool, isContinue bool) {
	var keyPrefix string
	hasFirstAd := false
	redisPool := ctx.Value(RedisKey).(*redis.Pool)
	redisConn := redisPool.Get()
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	uids, _ := b.GetUids()
	if len(uids) == 0 {
		keyPrefix = ""
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
		clog.StopTimer("redis_lock")
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

func parseInput(ctx context.Context, r *http.Request) (*BidRequest, error) {
	var pmpRequest proto.Request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &pmpRequest)
	if err != nil {
	}
	return &BidRequest{
		pmpRequest,
	}, nil
}

//TODO: add self win notice and macro replacement
func (b *BidRequest) ImpressionUrl(ctx context.Context, candidateAd *logic.Activity) []proto.IURL {
	impressions := candidateAd.ImpressionUrl()
	mgtvIUrl := make([]proto.IURL, 0, len(impressions)+1)

	//dsp win_notice url
	mgtvIUrl = append(mgtvIUrl, proto.IURL{Event: proto.ImpressionTypeZero, Url: WinNoticeUrl})

	//other third party url
	for i := 0; i < len(impressions); i++ {
		mgtvIUrl = append(mgtvIUrl, proto.IURL{Event: int(impressions[i].Event), Url: impressions[i].Url})
	}
	return mgtvIUrl
}

func (b *BidRequest) ClickUrl(ctx context.Context, candidateAd *logic.Activity) []string {
	return candidateAd.ClickUrl()
}
