package lock

import (
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_OnceInXSecond(t *testing.T) {
	assert := assert.New(t)
	c, err := redis.Dial("tcp", "10.100.2.90:6379")
	assert.Nil(err)
	//userTag, err := RequestTag(context.Background(), user)
	getLocked, err := OnceInXSecond(c, "test", 2)
	assert.Nil(err)
	assert.True(getLocked)

	getLocked, err = OnceInXSecond(c, "test", 2)
	assert.Nil(err)
	assert.False(getLocked)

	time.Sleep(time.Second)
	getLocked, err = OnceInXSecond(c, "test", 2)
	assert.Nil(err)
	assert.False(getLocked)

	time.Sleep(time.Second)
	getLocked, err = OnceInXSecond(c, "test", 2)
	assert.Nil(err)
	assert.True(getLocked)
}

func Test_OnceInXSecondNull(t *testing.T) {
	assert := assert.New(t)
	c, err := redis.Dial("tcp", "10.100.2.90:6379")
	assert.Nil(err)
	for i := 0; i < 100; i++ {
		getLocked, err := OnceInXSecond(c, "", 2)
		assert.NotNil(err)
		assert.False(getLocked)
	}
}
