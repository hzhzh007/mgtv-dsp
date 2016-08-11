package logic

import (
	"math/rand"
	"time"
)

const (
	ProbabilityDefault = 0
	ProbabilityMax     = 10000
)

type Probability int

func (p Probability) UnderProbability() bool {
	if p == ProbabilityDefault || p > ProbabilityMax {
		return true
	}
	return rand.Intn(ProbabilityMax) < int(p)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
