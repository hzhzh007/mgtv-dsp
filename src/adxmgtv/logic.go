package adxmgtv

import (
	"common/bridge"
	. "common/consts"
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
	//	cl.Debug("row activities:%+v", *candidateAds)
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

	//tag filter
	uids, _ := request.GetUids()
	clog.StartTimer()
	userTags, err := tag.RequestTag(ctx, bridge.LogicUids2RpcUser(uids))
	clog.StopTimer("tag_rpc")
	candidateAds.UserTagFilter(ctx, bridge.RpcUserTag2LogicTag(userTags))

	candidateAds.Sort()

	//duplicate removal
	candidateAds.CommonFilter(ctx, request.DuplicateRemovalFn(ctx), "mgtv reduplicate")

	candidateAds.Sort()

	return candidateAds, nil
}
