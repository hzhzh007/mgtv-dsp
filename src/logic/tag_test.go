package logic

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_TagInclude(t *testing.T) {
	assert := assert.New(t)
	tag1 := Tag(101)
	tag2 := Tag(1000101)
	assert.True(tag1.In([]Tag{tag1}))
	assert.True(tag2.In([]Tag{tag2}))

	//with category 100101 in [101]
	assert.True(tag1.In([]Tag{tag2}))

	assert.False(tag2.In([]Tag{tag1}))
}

func BenchmarkTagIn(b *testing.B) {
	assert := assert.New(b)
	_ = assert
	tag1 := Tag(101)
	tag2 := Tag(1000101)
	tag3 := Tag(1000201)
	sum := 0
	tags := []Tag{tag3, tag3, tag3, tag3, tag2}
	for i := 0; i < b.N; i++ {
		if tag1.In(tags) {
			sum++
		}
	}
	b.Log(sum)
}

func BenchmarkTagIn2(b *testing.B) {
	assert := assert.New(b)
	_ = assert
	tag1 := Tag(101)
	tag2 := Tag(1000101)
	tag3 := Tag(1000201)
	sum := 0
	tags := []Tag{tag3, tag3, tag3, tag3, tag2}
	for i := 0; i < b.N; i++ {
		if tag1.In2(tags) {
			sum++
		}
	}
	b.Log(sum)
}
