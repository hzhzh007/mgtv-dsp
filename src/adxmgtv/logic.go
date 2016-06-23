package adxmgtv

import (
	"common/bridge"
	. "common/consts"
	"fmt"
	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"logic"
	"tag"
	"time"
)

//TODO: implement it
func filterAd(ctx context.Context, request *BidRequest) (*logic.Activities, error) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	candidateAds, err := bridge.GetActivitiesCopy(request)
	if err != nil {
		clog.Error("get candidateAds error:%s", err)
		return candidateAds, err
	}

	//log some info
	request.LogInfo(ctx, candidateAds)

	clog.StartTimer()
	/*
	 *some filte logic
	 */
	//location filter
	location, _ := request.GetLocation(ctx)
	candidateAds.LocationFilter(ctx, location)

	// schedule time  filter
	candidateAds.ScheduleFilter(ctx, time.Now()) //here should cache

	//platform filter
	platform, _ := request.GetPlatform(ctx)
	candidateAds.PlatformFilter(ctx, platform)
	clog.StopTimer("filter1_cost")

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

	//first sort
	clog.StartTimer()
	candidateAds.Sort()
	clog.StopTimer("sort1_cost")

	//duplicate removal
	candidateAds.CommonFilter(ctx, request.DuplicateRemovalFn(ctx), "mgtv reduplicate")

	candidateAds.Sort()

	return candidateAds, nil
}
