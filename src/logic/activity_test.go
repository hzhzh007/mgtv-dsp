package logic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TagIncluePriorierThanExcludeOk(t *testing.T) {
	assert := assert.New(t)
	a := Activity{
		IncludeTag: []Tag{1000101, 2000101},
		ExcludeTag: []Tag{1000101, 2000101},
	}
	tag := []Tag{1000101, 2000101}
	assert.True(a.TagOK(tag))
}
func Test_TagNormalInlucdeOk(t *testing.T) {
	assert := assert.New(t)
	a := Activity{
		IncludeTag: []Tag{1000101, 2000101},
	}
	tag := []Tag{1000101, 2000101}
	assert.True(a.TagOK(tag))

	tag = []Tag{1000101}
	assert.False(a.TagOK(tag))

	tag = []Tag{}
	assert.False(a.TagOK(tag))

	tag = []Tag{2000101}
	assert.False(a.TagOK(tag))

	tag = []Tag{3000101}
	assert.False(a.TagOK(tag))

	tag = []Tag{1001000101, 20012000101}
	assert.True(a.TagOK(tag))
}

func Test_TagNormalExlucdeOk(t *testing.T) {
	assert := assert.New(t)
	a := Activity{
		ExcludeTag: []Tag{3000101, 4000101},
	}
	tag := []Tag{3000101, 4000101}
	assert.False(a.TagOK(tag))

	tag = []Tag{3000101}
	assert.False(a.TagOK(tag))

	tag = []Tag{4000101}
	assert.False(a.TagOK(tag))

	tag = []Tag{1004000101}
	assert.False(a.TagOK(tag))

	tag = []Tag{2000101}
	assert.True(a.TagOK(tag))
}

func Test_TagNormalExlucdeAndIncludeOk(t *testing.T) {
	assert := assert.New(t)
	a := Activity{
		IncludeTag: []Tag{1000101, 2000101},
		ExcludeTag: []Tag{3000101, 4000101},
	}
	tag := []Tag{3000101, 4000101}
	assert.False(a.TagOK(tag))

	tag = []Tag{1000101, 2000101, 3000101}
	assert.True(a.TagOK(tag))

	tag = []Tag{1000101, 2000101, 4000101}
	assert.True(a.TagOK(tag))

	tag = []Tag{10001000101, 2000101, 4000101}
	assert.True(a.TagOK(tag))

	tag = []Tag{4000101}
	assert.False(a.TagOK(tag))
}

func Test_Status(t *testing.T) {
	assert := assert.New(t)
	a := Activity{
		filtered: false,
	}
	assert.True(a.Filtered())

	a.Status = 1
	assert.False(a.Filtered())

	a.Status = 2
	assert.True(a.Filtered())

	a.filtered = true
	assert.True(a.Filtered())

	a.filtered = true
	a.Status = 2
	assert.True(a.Filtered())
}
