package adx_mgtv

import (
	"github.com/guregu/kami"
	cl "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	. "mgtv-dsp/common/consts"
	"net/http"
)

func bidHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	cl := cl.NewContext("mgtv_pmp")
	defer cl.Flush()

	ctx = context.WithValue(ctx, ContextLogKey, cl)
	request, err := parseInput(ctx, r)
	if err != nil {
		cl.Error("parse input error:%s", err)
		ResponseEmpty(ctx, request, w)
		return
	}
}

//TODO: May here we need to config the path from the config file
func InitHandler(ctx context.Context) {
	//	conf := ctx.Value(ConfKey).(*config.DspConfig)
	kami.Post("/mgtv/bid", bidHandler)
}
