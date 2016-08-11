package logic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Default(t *testing.T) {
	assert := assert.New(t)
	var p Probability
	p = 0
	assert.True(p.UnderProbability())
	assert.True(p.UnderProbability())

	p = 10000
	assert.True(p.UnderProbability())
	assert.True(p.UnderProbability())
	p = 10001
	assert.True(p.UnderProbability())
	assert.True(p.UnderProbability())
}

func Test_Probability(t *testing.T) {
	assert := assert.New(t)
	var p Probability
	p = 100
	sum := 0
	for i := 0; i < 100000; i++ {
		if p.UnderProbability() {
			sum++
		}
	}
	t.Log("100 from 100000", sum)
	assert.True(sum > 500)
	assert.True(sum < 2000)
}
