package adxmgtv

import (
	. "common/consts"
	"common/encrypt"
	"config"
	"fmt"
	"github.com/guregu/kami"
	clog "github.com/hzhzh007/context_log"
	"github.com/mattrobenolt/emptygif"
	"golang.org/x/net/context"
	"net/http"
)

const (
	WinNoticePriceKey = "c"
	DecryptedPriceKey = "price"
)

var (
	aesKey           string
	WinNoticeUrl     string
	ClickNoticeUrl   string
	LocaleRecordHost string
)

func BidHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	clog := clog.NewContext("adx=mgtv_pmp")
	defer clog.Flush()

	ctx = context.WithValue(ctx, ContextLogKey, clog)
	clog.StartTimer()
	request, err := parseInput(ctx, r)
	clog.StopTimer("parse_cost")
	if err != nil {
		clog.Error("parse input error:%s", err)
		ResponseEmpty(ctx, request, w)
		return
	}
	result, err := filterAd(ctx, request)

	clog.StartTimer()
	if err != nil || result.AvailableLen() == 0 {
		ResponseEmpty(ctx, request, w)
		return
	}
	ResponseAd(ctx, w, request, (*result)[0:1])
	clog.StopTimer("resp_cost")

}

//TODO: here can do some other things like flow and frequency controll
func WinNoticeHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	clog := clog.NewContext("adx=mgtv_win_notice")
	defer clog.Flush()

	w.Header().Set("Connection", "close")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Content-Type", "image/gif")

	r.ParseForm()
	forms := r.Form

	clog.AddNotes("path", r.URL.Path)
	price := forms[WinNoticePriceKey]
	if len(price) > 0 {
		decrypted, err := encrypt.AesCBCDecrypte(price[0], aesKey)
		if err != nil {
			clog.Error("decrypt str:%s error:%s", price[0], err)
		} else if LocaleRecordHost != "" {
			http.Get(fmt.Sprintf("http://%s/%s?%s", LocaleRecordHost, r.URL.Path, forms.Encode()))
		}
		forms.Add(DecryptedPriceKey, decrypted)
		clog.AddNotes("", map2String(forms))
	}
	w.Write(emptygif.PIXEL)
}

func InitHandler(ctx context.Context) {
	conf := ctx.Value(ConfKey).(*config.DspConfig)

	//some local variables init
	aesKey = conf.Mgtv.Key
	WinNoticeUrl = conf.Mgtv.WinNoticeUrl
	ClickNoticeUrl = conf.Mgtv.ClickNonticeUrl
	LocaleRecordHost = conf.Mgtv.RedirectHost

	//path set
	kami.Post(conf.Mgtv.BidPath, BidHandler)
	kami.Get(conf.Mgtv.NoticePath, WinNoticeHandler)
}
