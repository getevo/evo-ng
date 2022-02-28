package distributed

import (
	"github.com/getevo/evo-ng/apps/redis"
	"time"
)

type MachineLock struct {
	Key             string        `json:"key"`
	Holder          string        `json:"holder"`
	This            string        `json:"-"`
	DefaultDuration time.Duration `json:"default_timeout"`
}

func NewMachineLock(key string) *MachineLock {
	if !redis.Ready() {
		panic("redis is not connected")
	}
	return &MachineLock{
		Key:  "mlock#" + key,
		This: GetAppName(),
	}
}

func (l *MachineLock) SetDefaultDuration(duration time.Duration) *MachineLock {
	l.DefaultDuration = duration
	return l
}

func (l *MachineLock) TryAcquire() bool {
	redis.Get(l.Key, &l.Holder)
	var canAcquire = l.Holder == "" || l.Holder == l.Key
	if canAcquire {
		redis.Set(l.Key, l.This, l.DefaultDuration)
	}
	return canAcquire
}

func (l *MachineLock) TryAcquireDuration(duration time.Duration) bool {
	redis.Get(l.Key, &l.Holder)
	var canAcquire = l.Holder == "" || l.Holder == l.Key
	if canAcquire {
		redis.Set(l.Key, l.This, duration)
	}
	return canAcquire
}

func (l *MachineLock) TryUntilAcquire(retry time.Duration, timeout time.Duration) bool {
	var canAcquire = false
	for timeout := time.After(time.Second); ; {
		select {
		case <-timeout:
			return false
		default:
			redis.Get(l.Key, &l.Holder)
			canAcquire = l.Holder == "" || l.Holder == l.Key
			if canAcquire {
				redis.Set(l.Key, l.This, l.DefaultDuration)
			}
			if canAcquire {
				return true
			}
			time.Sleep(retry)
		}

	}
}

func (l *MachineLock) Reset() {
	redis.Del(l.Key)
}
