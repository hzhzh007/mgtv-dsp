package logic

import (
	"fmt"
	"strconv"
	"time"
)

const (
	ActivityStatusRun = 1
)

type Activity struct {
	//activity id
	Id int `yaml:"id"`

	ActiveTime Schedules `yaml:"active_time"`

	//for mgtv the platform type
	//for the mgtv:`1 PC`  `21 ANDROID_TABLET_H5`  `22 ANDROID_PHONE_H5`  `23 IPAD_H5`  `24 IPHONE_H5`  `31 ANDROID_TABLET_APP`  `32 ANDROID_PHONE_APP`  `33 IPAD_APP`   `34 IPHONE_APP`   `41 XIAOMI_SDK`   `42 BAIDU_SDK`   `100 OTT`  `101 OTT_NEW`
	//for others, we can have more personalized definitions
	Platform        []Platform  `yaml:"platform"`
	AdType          int         `yaml:"ad_type"`
	IncludeLocation []Location  `yaml:"include_location"`
	IncludeTag      []Tag       `yaml:"include_tag"`
	ExcludeTag      []Tag       `yaml:"exclude_tag"`
	Length          Duration    `yaml:"length"`
	Creative        []Creative  `yaml:"creative"`
	Click           []string    `yaml:"click_url"`
	MonitorUrl      Impressions `yaml:"monitor_url"`
	LandingPage     string      `yaml:"landing_page"`
	Percent         Probability `yaml:"percent"` //0-10000
	MaxPrice        int         `yaml:"max_price"`
	Pmp             []PMP       `yaml:"pmp"`

	Flow      FlowType      `yaml:"flow"`
	Frequency FrequencyType `yaml:"frequency"`

	//the flag wen filter
	filtered          bool
	selectedCreateive *Creative
	selectedDeal      *PMP
	Status            int `yaml:"status"`
	//the below is the extention as the dsp developed
	//ExcludeLocations []Location   `yaml:"exclude_location"`
}

func (a *Activity) LocationOK(requestLocation Location) bool {
	if len(a.IncludeLocation) == 0 {
		return true
	}
	for _, location := range a.IncludeLocation {
		if location.Include(requestLocation) {
			return true
		}
	}
	return false
}

func (a *Activity) ScheduleOK(requestTime time.Time) bool {
	if len(a.ActiveTime) == 0 {
		return true
	}
	for _, schedule := range a.ActiveTime {
		if schedule.UnderSchedule(requestTime) {
			return true
		}
	}
	return false
}

//inlcude priority higher than exclude priority
func (a *Activity) TagOK(userTag []Tag) bool {
	lenIncludeTag := len(a.IncludeTag)
	lenExcludeTag := len(a.ExcludeTag)
	if lenIncludeTag == 0 && lenExcludeTag == 0 {
		return true
	}
	//include stratage priority is higher
	for _, tag := range a.IncludeTag {
		if !tag.In(userTag) {
			return false
		}
	}
	if lenIncludeTag > 0 {
		return true
	}

	for _, tag := range a.ExcludeTag {
		if tag.In(userTag) {
			return false
		}
	}
	return true
}

func (a *Activity) Filtered() bool {
	return a.filtered || a.Status != ActivityStatusRun
}

func (a *Activity) SetFiltered() bool {
	a.filtered = true
	return a.filtered
}

//TODO: implement it
//sometimes it may the heart of the dsp
func (a *Activity) GetECPM() int {
	return a.MaxPrice
}

func (a *Activity) Price() int {
	return a.MaxPrice
}
func (a *Activity) PlatformOK(requestPlatform Platform) bool {
	if len(a.Platform) == 0 {
		return true
	}
	for _, platform := range a.Platform {
		if platform == requestPlatform {
			return true
		}
	}
	return false
}

//TODO: select the creative
func (a *Activity) selectedCreative() *Creative {
	if a.selectedCreateive == nil {
		a.selectedCreateive = &a.Creative[0]
	}
	return a.selectedCreateive
}

func (a *Activity) SetDeal(i int) {
	a.selectedDeal = &a.Pmp[i]
}
func (a *Activity) SelectedDealId() string {
	if a.selectedDeal == nil {
		return ""
	}
	return a.selectedDeal.DealId
}

func (a *Activity) CreativeUrl() string {
	return a.selectedCreative().Url
}

func (a *Activity) LandingPageUrl() string {
	selectedCreative := a.selectedCreative()
	if selectedCreative.LandingPage == "" {
		return a.LandingPage
	}
	return selectedCreative.LandingPage
}

func (a *Activity) ImpressionUrl() Impressions {
	selectedCreative := a.selectedCreative()
	if len(selectedCreative.MonitorUrl) != 0 {
		return selectedCreative.MonitorUrl
	}
	return a.MonitorUrl
}

func (a *Activity) ClickUrl() []string {
	selectedCreative := a.selectedCreative()
	if len(selectedCreative.Click) != 0 {
		return selectedCreative.Click
	}
	return a.Click
}

func (a *Activity) ActiveId() string {
	return strconv.Itoa(a.Id)
}
func (a *Activity) CreativeId() string {
	selectedCreative := a.selectedCreative()
	return strconv.Itoa(selectedCreative.Id)
}

func (a *Activity) Duration() int {
	return a.selectedCreateive.Duration
}

//TODO: implement it
func (a *Activity) IsFlowOK(today, total int) bool {
	for _, flowSetting := range a.Flow {
		if !flowSetting.UnderFlow(today, total, a.ActiveTime, time.Now()) {
			return false
		}
	}
	return true
}

func (a *Activity) CreativeWidth() int {
	return a.selectedCreative().Width
}

func (a *Activity) CreativeHeight() int {
	return a.selectedCreative().Height
}

func (a *Activity) CreativeType() int {
	return a.selectedCreative().Type
}
func (a *Activity) CreativeStart() time.Time {
	return a.ActiveTime[0].Start
}

//just for log
func (a *Activity) String(msg string) string {
	return fmt.Sprintf("{activity:%d,creative:%d,pd:%s,%s}", a.Id, a.selectedCreative().Id, a.SelectedDealId(), msg)
}
