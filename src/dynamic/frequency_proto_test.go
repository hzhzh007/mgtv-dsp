package dynamic

import (
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_PbStorage(t *testing.T) {
	assert := assert.New(t)
	redisValue := RedisValue{}
	var i int32
	for i = 1; i <= 100; i++ {
		redisValue.Impression = append(redisValue.Impression, &Record{
			Id:      i + 1000,
			Expire:  int32(time.Now().Unix()),
			Type:    FreqType_FreqPerWeek,
			Counter: i,
		})
		data, err := proto.Marshal(&redisValue)
		assert.Nil(err)
		assert.NotZero(len(data))
		//t.Logf("record_num:%d, len:%d\n", i, len(data))
	}
}

//@Result: mpb15` BenchmarkPbSpeed-4   	 5000000	       286 ns/op 1.825s
func BenchmarkPbUmarrshalSpeed(b *testing.B) {
	assert := assert.New(b)
	redisValue := RedisValue{}
	var data []byte
	var err error
	var i int32
	for i = 1; i <= 10; i++ {
		redisValue.Impression = append(redisValue.Impression, &Record{
			Id:      i + 1000,
			Expire:  int32(time.Now().Unix()),
			Type:    FreqType_FreqPerWeek,
			Counter: i,
		})
		data, err := proto.Marshal(&redisValue)
		assert.Nil(err)
		assert.NotZero(len(data))
	}
	newTest := &RedisValue{}
	for j := 0; j <= b.N; j++ {
		err = proto.Unmarshal(data, newTest)
		assert.Nil(err)
	}
	b.Log(newTest)
}

//@RESULT: mpb15` BenchmarkPbMarrshalSpeed-4    	 1000000	      1957 ns/op
func BenchmarkPbMarrshalSpeed(b *testing.B) {
	assert := assert.New(b)
	redisValue := RedisValue{}
	var data []byte
	var err error
	var i int32
	for i = 1; i <= 10; i++ {
		redisValue.Impression = append(redisValue.Impression, &Record{
			Id:      i + 1000,
			Expire:  int32(time.Now().Unix()),
			Type:    FreqType_FreqPerWeek,
			Counter: i,
		})
	}
	for j := 0; j <= b.N; j++ {
		data, err = proto.Marshal(&redisValue)
		assert.Nil(err)
	}
	b.Log(len(data))
}
