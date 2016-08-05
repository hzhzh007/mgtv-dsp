//adapter of request and dsp logic filter
package adxmgtv

import (
	"bytes"
	. "common/consts"
	"encoding/json"
	"fmt"
	clog "github.com/hzhzh007/context_log"
	proto "github.com/hzhzh007/mgtvPmpProto"
	"golang.org/x/net/context"
	"io/ioutil"
	"iplib"
	"logic"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	ProtocolVersion          = 3
	DuplicateRemovalInterval = 2 // 2s
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type BidRequest struct {
	proto.Request
	city int
}

//TODO: here iplib do not handle err
func (b *BidRequest) GetLocation(ctx context.Context) (logic.Location, error) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	iplib := ctx.Value(IPLibKey).(iplib.IpCollections)
	ipStr := b.Device.Ip
	city := iplib.Ip2CityCode(ipStr)
	b.city = int(city)
	clog.AddNotes("city", strconv.Itoa(int(city)))
	return logic.Location(city), nil
}

func (b *BidRequest) VideoDuration(ctx context.Context) int {
	return b.Video.Duration
}

func (b *BidRequest) AdType(ctx context.Context) int {
	if len(b.Imp) > 0 {
		return b.Imp[0].Location
	}
	return 0
}

//here we use the mgtv platform type, so we do not need to translate it
func (b *BidRequest) GetPlatform(ctx context.Context) (logic.Platform, error) {
	return logic.Platform(b.Device.Type), nil
}

func (b *BidRequest) LogInfo(ctx context.Context, candidateAds *logic.Activities) {
	clog := ctx.Value(ContextLogKey).(*clog.ServerContext)
	clog.AddNotes("bid", b.Bid)
	if len(b.Imp) > 0 {
		clog.AddNotes("slot", b.Imp[0].SpaceId)
		clog.AddNotes("player", strconv.Itoa(b.Imp[0].PlayerId))
		clog.AddNotes("location", strconv.Itoa(b.Imp[0].Location))
		clog.AddNotes("order", strconv.Itoa(b.Imp[0].Order))
	}
	clog.AddNotes("ip", b.Device.Ip)
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

func parseInput(ctx context.Context, r *http.Request) (*BidRequest, error) {
	var pmpRequest proto.Request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &pmpRequest)
	return &BidRequest{
		Request: pmpRequest,
	}, err
}

//TODO: add self win notice and macro replacement
func (b *BidRequest) ImpressionUrl(ctx context.Context, candidateAd *logic.Activity, index int) []proto.IURL {
	replacer := b.CreateReplacer(candidateAd, index)
	impressions := candidateAd.ImpressionUrl()
	mgtvIUrl := make([]proto.IURL, 0, len(impressions)+1)

	//dsp win_notice url
	mgtvIUrl = append(mgtvIUrl, proto.IURL{Event: proto.ImpressionTypeZero, Url: replacer.Replace(WinNoticeUrl)})

	//other third party url
	for i := 0; i < len(impressions); i++ {
		mgtvIUrl = append(mgtvIUrl, proto.IURL{Event: int(impressions[i].Event), Url: replacer.Replace(impressions[i].Url)})
	}
	return mgtvIUrl
}

//TODO: to be more efficient here click and impression replcer can use the samve and reduce once creation
func (b *BidRequest) ClickUrl(ctx context.Context, candidateAd *logic.Activity, index int) []string {
	clickUrl := candidateAd.ClickUrl()
	replacer := b.CreateReplacer(candidateAd, index)
	result := make([]string, 0, len(clickUrl))

	//self monitor
	result = append(result, replacer.Replace(ClickNoticeUrl))
	for _, url := range clickUrl {
		result = append(result, replacer.Replace(url))
	}
	return result
}

func (b *BidRequest) CreateReplacer(activity *logic.Activity, index int) *strings.Replacer {
	var uid string
	uids, _ := b.GetUids()
	if len(uids) > 0 {
		uid = uids[0].Value
	}
	replacer := strings.NewReplacer("${ADSPACE_ID}", b.Imp[index].SpaceId,
		"${ORDER}", strconv.Itoa(b.Imp[index].Order),
		"${PLAYER_ID}", strconv.Itoa(b.Imp[index].PlayerId),
		"${CREATIVE_ID}", activity.CreativeId(),
		"${ACTIVE_ID}", activity.ActiveId(),
		"${CONTENT_ID}", strconv.Itoa(b.Video.VideoId),
		"${HID}", strconv.Itoa(b.Video.CollectionId),
		"${ADTYPE}", strconv.Itoa(activity.AdType),
		"${UID}", uid,
		"${BID}", b.Bid,
		"${PD}", activity.SelectedDealId(),
		"${CITY}", strconv.Itoa(b.city),
		"__OS__", b.Device.Os,
		"__IMEI__", b.Device.Imei,
		"__MAC__", b.Device.Mac,
		"__IDFA__", b.Device.Idfa,
		"__OPENUDID__", b.Device.Openudid,
		"__ANDROIDID__", b.Device.Anid,
		"__UDID__", b.Device.Udid,
		"__ODIN__", b.Device.Odin,
		"__DUID__", b.Device.Duid,
		"__IP__", b.Device.Ip,
		"__UA__", url.QueryEscape(b.Device.Ua),
		"__TS__", fmt.Sprintf("%d", time.Now().Unix()))

	return replacer
}

func map2String(m url.Values) string {
	var buf bytes.Buffer
	for key, value := range m {

		buf.WriteByte(' ')
		buf.WriteString(key)
		buf.WriteByte('=')
		if len(value) > 0 {
			buf.WriteString(value[0])
		}
	}
	return buf.String()
}
