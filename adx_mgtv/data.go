package adx_mgtv

import (
	"encoding/json"
	proto "github.com/hzhzh007/mgtvPmpProto"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
)

const (
	ProtocolVersion = 1
)

type BidRequest struct {
	proto.Request
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
