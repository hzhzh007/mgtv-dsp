package logic

type PMP struct {
	DealId string `yaml:"deal_id"`
}

func (p *PMP) Match(dealId string) bool {
	return dealId == p.DealId
}
