package logic

import (
	"fmt"
	"strconv"
	"time"
)

type Activity struct {
	Id              int         `yaml:"id"`
	ActiveTime      []Schedule  `yaml:"active_time"`
	Platform        []Platform  `yaml:"platform"`
	IncludeLocation []Location  `yaml:"include_location"`
	IncludeTag      []Tag       `yaml:"include_tag"`
	Creative        []Creative  `yaml:"creative"`
	Click           []string    `yaml:"click_url"`
	MonitorUrl      Impressions `yaml:"monitor_url"`
	LandingPage     string
	MaxPrice        int `yaml:"max_price"`

	//the flag wen filter
	filtered          bool
	selectedCreateive *Creative
	//the below is the extention as the dsp developed
	//ExcludeLocations []Location   `yaml:"exclude_location"`
	//ExcludeTag  []Tag        `yaml:"exclude_tag"`
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
func (a *Activity) TagOK(userTag []Tag) bool {
	if len(a.IncludeTag) == 0 {
		return true
	}
	for _, tag := range a.IncludeTag {
		if tag.In(userTag) {
			return true
		}
	}
	return false
}

func (a *Activity) Filtered() bool {
	return a.filtered
}

func (a *Activity) SetFiltered() bool {
	a.filtered = true
	return a.filtered
}

//TODO: implement it
func (a *Activity) GetECPM() int {
	return a.Id
}

func (a *Activity) Price() int {
	return a.MaxPrice
}
func (a *Activity) PlatformOK(requestPlatform Platform) bool {
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
func (a *Activity) Duration() int {
	return a.selectedCreateive.Duration
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

func (a *Activity) String(msg string) string {
	return fmt.Sprintf("{activity:%d, creative:%d,%s}", a.Id, a.selectedCreative().Id, msg)
}
