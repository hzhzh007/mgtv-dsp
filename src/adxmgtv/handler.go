package adxmgtv

import (
	. "common/consts"
	"common/encrypt"
	"github.com/guregu/kami"
	clog "github.com/hzhzh007/context_log"
	"golang.org/x/net/context"
	"net/http"
)

const (
	WinNoticePriceKey = "c"
	DecryptedPriceKey = "price"
	LocaleRecordUrl   = "http://127.0.0.1/win/record?"
)

var (
	aesKey string = "975dfad4e0b94c38"
)

func BidHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	clog := clog.NewContext("adx=mgtv_pmp")
	defer clog.Flush()

	ctx = context.WithValue(ctx, ContextLogKey, clog)
	clog.StartTimer()
	request, err := parseInput(ctx, r)
	clog.StopTimer("parseInput")
	if err != nil {
		clog.Error("parse input error:%s", err)
		ResponseEmpty(ctx, request, w)
		return
	}
	result, err := filterAd(ctx, request)
	if err != nil || result.AvailableLen() == 0 {
		ResponseEmpty(ctx, request, w)
		return
	}
	ResponseAd(ctx, w, request, (*result)[0:1])

}

func WinNoticeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	clog := clog.NewContext("adx=mgtv_win_notice")
	defer clog.Flush()

	w.Header().Set("Connection", "close")
	w.Header().Set("Cache-Control", "no-store")

	r.ParseForm()
	forms := r.Form

	price := forms[WinNoticePriceKey]
	if len(price) > 0 {
		decrypted, err := encrypt.AesCBCDecrypte(price[0], aesKey)
		if err != nil {
			clog.Error("decrypt str:%s error:%s", price[0], err)
		} else {
			clog.AddNotes("url", forms.Encode())
			forms.Add(DecryptedPriceKey, decrypted)
			http.Get(LocaleRecordUrl + forms.Encode())
		}
	}
}

//TODO: May here we need to config the path from the config file
func InitHandler(ctx context.Context) {
	//	conf := ctx.Value(ConfKey).(*config.DspConfig)
	kami.Post("/mgtv/bid", BidHandler)
	kami.Get("/winnotice", WinNoticeHandler)
}
