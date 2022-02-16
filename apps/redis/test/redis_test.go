package test

import (
	"github.com/getevo/evo-ng/apps/redis"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var connected = false

func Init() {
	if !connected {
		redis.Connect(&redis.Config{
			AppID:  "RedisTest",
			Server: []string{"10.15.10.193:6379"},
		})
	}
}

type Test struct {
	String  string
	Integer int
}

func TestSimple(t *testing.T) {
	Init()
	assert.NoError(t, redis.Set("x", "y", 10*time.Second))
	var outString string
	assert.NoError(t, redis.Get("x", &outString))
	assert.Equal(t, outString, "y")
	assert.NoError(t, redis.Del("x"))
	assert.Error(t, redis.Get("x", &outString))

	v := Test{
		"A", 10,
	}
	assert.NoError(t, redis.Set("y", v, 10*time.Second))
	var outTest Test
	assert.NoError(t, redis.Get("y", &outTest))
	assert.Equal(t, outTest, v)
	assert.NoError(t, redis.Del("y"))
	assert.Error(t, redis.Get("y", &outTest))

}
