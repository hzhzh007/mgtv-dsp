package adxmgtv

import (
	. "common/consts"
	"common/encrypt"
	"config"
	"feedback"
	"fmt"
	"github.com/guregu/kami"
	clog "github.com/hzhzh007/context_log"
	"github.com/mattrobenolt/emptygif"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
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
	Feedback         *feedback.Feedback
)

//TODO: here can do some other things like flow and frequency controll
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
		}
		forms.Add(DecryptedPriceKey, decrypted)
	}
	clog.AddNotes("", map2String(forms))
	clog.Debug("redirect: %s", fmt.Sprintf("http://%s/%s?%s", LocaleRecordHost, r.URL.Path, forms.Encode()))
	if LocaleRecordHost != "" {
		http.Get(fmt.Sprintf("http://%s/%s?%s", LocaleRecordHost, r.URL.Path, forms.Encode()))
	}
	if strings.HasSuffix(r.URL.Path, "win") {
		select {
		case Feedback.Input <- feedback.Record{
			Uid:      formsString(forms, "d"),
			Type:     feedback.RecordImpressionType,
			Activity: formInt(forms, "cd"),
		}:
		default:
			clog.Error("feedback queue full:%s", forms.Encode())
		}
	}
	w.Write(emptygif.PIXEL)
}

func formsString(form url.Values, key string) string {
	v := form[key]
	if len(v) > 0 {
		return v[0]
	}
	return ""
}
func formInt(form url.Values, key string) int {
	v := form[key]
	if len(v) > 0 {
		vi, _ := strconv.Atoi(v[0])
		return vi
	}
	return 0
}

func InitHandler(ctx context.Context) {
	conf := ctx.Value(ConfKey).(*config.DspConfig)
	Feedback = ctx.Value(FeedbackKey).(*feedback.Feedback)

	//some local variables init
	aesKey = conf.Mgtv.Key
	WinNoticeUrl = conf.Mgtv.WinNoticeUrl
	ClickNoticeUrl = conf.Mgtv.ClickNonticeUrl
	LocaleRecordHost = conf.Mgtv.RedirectHost

	//path set
	kami.Post(conf.Mgtv.BidPath, BidHandler)
	kami.Get(conf.Mgtv.NoticePath, WinNoticeHandler)
}
