//some input/output stream operation
package adxmgtv

import (
	. "common/consts"
	"encoding/json"
	"fmt"
	clog "github.com/hzhzh007/context_log"
	proto "github.com/hzhzh007/mgtvPmpProto"
	"golang.org/x/net/context"
	"logic"
	"net/http"
)

//TODO: here only response one ad, update 2016160621
func ResponseAd(ctx context.Context, w http.ResponseWriter, bidRequest *BidRequest, candidateAds logic.Activities) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	candidateAd := &candidateAds[0]
	responseAd := proto.ADS{
		Price:           candidateAd.Price(),
		AdUrl:           candidateAd.CreativeUrl(),
		ClickThroughUrl: candidateAd.LandingPageUrl(),

		IUrl: bidRequest.ImpressionUrl(ctx, candidateAd, 0),

		CUrl:     bidRequest.ClickUrl(ctx, candidateAd, 0),
		AdId:     candidateAd.ActiveId(),
		Title:    "test", //TODO
		Duration: candidateAd.Duration(),
		Ctype:    proto.AdTypeVideo, //TODO: CreativeType()
		Width:    candidateAd.CreativeWidth(),
		Height:   candidateAd.CreativeHeight(),
		DealId:   candidateAd.SelectedDealId(),
	}
	reqsponse := proto.Response{
		ErrCode: proto.ResponseBidOK,
		Version: ProtocolVersion,
		Bid:     bidRequest.Bid,
		Ads:     []proto.ADS{responseAd},
	}
	WriteResponse(ctx, w, reqsponse, 200)
	clog.AddNotes("adinfo", candidateAd.String(fmt.Sprintf("price=%d", candidateAd.Price())))
}

func ResponseEmpty(ctx context.Context, bidRequest *BidRequest, w http.ResponseWriter) error {
	reqsponse := proto.Response{
		ErrCode: proto.ResponseNoAd,
		Version: ProtocolVersion,
	}

	if bidRequest != nil {
		reqsponse.Bid = bidRequest.Bid
	}
	WriteResponse(ctx, w, reqsponse, 200)
	return nil
}

func WriteResponse(ctx context.Context, w http.ResponseWriter, response interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
