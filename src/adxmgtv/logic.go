package adxmgtv

import (
	"common/bridge"
	. "common/consts"
	"fmt"
	"github.com/garyburd/redigo/redis"
	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"logic"
	"resource"
	"strconv"
	"tag"
	"time"
)

//TODO: implement it
func filterAd(ctx context.Context, request *BidRequest) (*logic.Activities, error) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	resource := ctx.Value(ResourceKey).(*resource.Resource)
	redisPool := ctx.Value(RedisKey).(*redis.Pool)
	candidateAds, err := resource.GetActivitiesCopy(request, true)
	if err != nil {
		clog.Error("get candidateAds error:%s", err)
		return candidateAds, err
	}
	//log some info
	clog.AddNotes("vid", strconv.Itoa(request.Video.VideoId))
	clog.AddNotes("hid", strconv.Itoa(request.Video.CollectionId))
	clog.AddNotes("v_len", strconv.Itoa(request.Video.Duration))
	clog.AddNotes("type", strconv.Itoa(request.Device.Type))
	clog.AddNotes("channel", strconv.Itoa(request.Device.Type))

	//log some info
	request.LogInfo(ctx, candidateAds)

	/*
	 *some filte logic
	 */

	//pd, prefer deal
	candidateAds.CommonFilter(ctx, request.PDFilterFn(ctx), "pd")

	//probability
	candidateAds.ProbabilityFilter(ctx)

	//location filter
	location, _ := request.GetLocation(ctx)
	candidateAds.LocationFilter(ctx, location)

	//duration filter
	candidateAds.DurationFilter(ctx, request.VideoDuration(ctx))

	//adtype filter
	candidateAds.AdTypeFilter(ctx, request.AdType(ctx))

	// schedule time  filter
	candidateAds.ScheduleFilter(ctx, time.Now()) //here should cache

	//platform filter
	platform, _ := request.GetPlatform(ctx)
	candidateAds.PlatformFilter(ctx, platform)

	//tag filter
	uids, _ := request.GetUids()

	if len(uids) > 0 { //just for performance
		clog.AddNotes("uid_t", uids[0].Type)
		clog.AddNotes("uid", uids[0].Value)
	}
	clog.StartTimer()
	userTags, err := tag.RequestTag(ctx, bridge.LogicUids2RpcUser(uids))
	clog.StopTimer("tag_rpc")
	localTags := bridge.RpcUserTag2LogicTag(userTags)
	clog.AddNotes("tags", fmt.Sprintf("%v", localTags))
	candidateAds.UserTagFilter(ctx, localTags)

	//get redis connection
	redisConn := redisPool.Get()
	defer redisConn.Close()

	//frequncyFilter
	candidateAds.CommonFilter(ctx, request.FrequencyFilterFn(ctx, redisConn), "freq")

	//first sort
	candidateAds.Sort()

	//duplicate removal
	candidateAds.CommonFilter(ctx, request.DuplicateRemovalFn(ctx, redisConn), "mgtv reduplicate")

	candidateAds.Sort()

	return candidateAds, nil
}
