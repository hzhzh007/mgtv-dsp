package logic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DurationZeroUnder(t *testing.T) {
	assert := assert.New(t)
	zeroDuration := Duration{0, 0}
	assert.True(zeroDuration.UnderDuration(0))
	assert.True(zeroDuration.UnderDuration(1))
	assert.True(zeroDuration.UnderDuration(10))
}
func Test_DurationLeftZeroUnder(t *testing.T) {
	assert := assert.New(t)
	zeroDuration := Duration{0, 360}
	assert.True(zeroDuration.UnderDuration(0))
	assert.True(zeroDuration.UnderDuration(1))
	assert.True(zeroDuration.UnderDuration(10))
	assert.True(zeroDuration.UnderDuration(300))
	assert.False(zeroDuration.UnderDuration(360))
	assert.False(zeroDuration.UnderDuration(370))
}

func Test_DurationRightZeroUnder(t *testing.T) {
	assert := assert.New(t)
	zeroDuration := Duration{360, 0}
	assert.False(zeroDuration.UnderDuration(0))
	assert.False(zeroDuration.UnderDuration(1))
	assert.False(zeroDuration.UnderDuration(10))
	assert.False(zeroDuration.UnderDuration(300))
	assert.True(zeroDuration.UnderDuration(360))
	assert.True(zeroDuration.UnderDuration(370))
}
