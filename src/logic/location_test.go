package logic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_LocationIn(t *testing.T) {
	assert := assert.New(t)
	var tangshan Location = 1156130200 //河北省	唐山市	1156130200
	var baoding Location = 1156130600  //河北省	保定市	1156130600
	var hebei Location = 1156130000
	var shaxi Location = 1156610000

	assert.False(tangshan.Include(baoding))

	assert.False(tangshan.Include(hebei))
	assert.False(shaxi.Include(hebei))

	assert.True(hebei.Include(tangshan))
	assert.True(hebei.Include(baoding))
	assert.True(hebei.Include(hebei))
}
