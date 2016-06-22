package tag

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_User2Array(t *testing.T) {
	assert := assert.New(t)
	userTag := UserTag{
		Tags: []*Tag{
			&Tag{Id: 1}, &Tag{Id: 5}, &Tag{Id: 3},
		},
	}
	intArray := UserTag2Array(userTag)
	assert.Equal(3, len(intArray))
	assert.EqualValues(1, intArray[0])
	assert.EqualValues(5, intArray[1])
	assert.EqualValues(3, intArray[2])
}
