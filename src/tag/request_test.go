package tag

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func init() {
	Init("127.0.0.1:8989", 10, time.Millisecond*100)
}

func Test_Request(t *testing.T) {
	assert := assert.New(t)
	user := NewUser("imei", "824224242342")
	userTag, err := RequestTag(context.Background(), user)
	assert.Nil(err)
	assert.NotNil(userTag)
}

func BenchmarkGrpc(b *testing.B) {
	assert := assert.New(b)
	user := NewUser("imei", "824224242342")
	userTag, err := RequestTag(context.Background(), user)
	assert.Nil(err)
	assert.NotNil(userTag)
	for i := 0; i < b.N; i++ {
		userTag, err := RequestTag(context.Background(), user)
		assert.Nil(err)
		assert.NotNil(userTag)
	}
	b.Log(userTag)
}
