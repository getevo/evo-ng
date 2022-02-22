package distributed

import (
	"github.com/getevo/evo-ng/apps/redis"
	"time"
)

type MachineLock struct {
	Key    string
	Holder string
	This   string
}

func NewMachineLock(key string) *MachineLock {
	return &MachineLock{
		Key:  key,
		This: GetAppName(),
	}
}

func (l *MachineLock) TryAcquire(duration time.Duration) bool {
	redis.Get(l.Key, &l.Holder)
	var canAcquire = l.Holder == "" || l.Holder == l.Key
	if canAcquire {
		redis.Set(l.Key, l.Holder, duration)
	}
	return canAcquire
}

func (l *MachineLock) TryUntilAcquire(duration time.Duration, retry time.Duration) bool {
	var canAcquire = false
	for {
		redis.Get(l.Key, &l.Holder)
		canAcquire = l.Holder == "" || l.Holder == l.Key
		if canAcquire {
			redis.Set(l.Key, l.Holder, duration)
		}
		if canAcquire {
			return
		}
	}
	return canAcquire
}
