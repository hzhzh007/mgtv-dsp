package logic

import (
	"time"
)

type UserId interface {
	GetIdType() string
	GEtIdValue() string
}

type AdOpptunities interface {
	GetAdOpptunities() []AdOpptunity
}
type CandidateAds interface {
	GetCandidateAd() CandidateAd
	CandidateAds()
}

//version0.1 we only need to filter
// 1.time
// 2.location
// 3.Platform
// 4.Tag
type AdOpptunity interface {
	// where the ip from
	Ip() string

	//
	AdType() int

	//
	Platform() int

	//
	Uids() []UserId
}

type CandidateAd interface {

	//是否符合地域定向
	LocationOK(Location) bool

	//是否符合排期定向
	ScheduleOK(time.Time) bool

	//是否符合平台定向
	PlatformOK(Platform) bool

	//是否符合tag定向
	TagOK([]Tag) bool

	//
	CreativeUrl() string

	//LoadingUrl
	LandingUrl() string

	//
	ImpressionUrl() []string

	//
	ClickUrl() []string

	//
	ExpectECMP() int

	//Price
	Price() int

	//Id()
	Id() string

	//Duration
	Duration() int

	//other common
	SetFilterStatus(bool)

	FilterStatus() bool
}
