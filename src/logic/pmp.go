package logic

type PMP struct {
	DealId string
}

func (p *PMP) Match(dealId string) bool {
	return dealId == p.DealId
}
