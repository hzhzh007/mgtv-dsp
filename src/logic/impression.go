package logic

const (
	ImpressionEventStart        = 0
	ImpressionEventOneQuarter   = 1
	ImpressionEventHalf         = 1
	ImpressionEventThreeQuarter = 1
	ImpressionEventEnd          = 1
)

type ImpressionEvent int

type Impression struct {
	Event ImpressionEvent
	Url   string
}

type Impressions []Impression
