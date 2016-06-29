//distributed lock
package lock

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"time"
)

const (
	GetLockOK = "OK"
)

var (
	ErrorEmptyKey = errors.New("empty key and return locked")
)

// only one can get lock at X seconds
// use redis and the cmd `SET KEY VALUE NX  EX time` as the distribute lock
// here set the key to Now() to make debug more easy
func OnceInXSecond(redisConn redis.Conn, key string, second int) (bool, error) {
	if key == "" {
		return false, ErrorEmptyKey
	}
	OK, err := redis.String(redisConn.Do("SET", key, time.Now().Unix(), "NX", "EX", second))
	if err == redis.ErrNil {
		err = nil
	}
	return OK == GetLockOK, err

}
