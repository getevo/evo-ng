package standard

import "time"

type Limiter interface {
	Reset()
	TryAcquireDuration(duration time.Duration) bool
	TryAcquire() bool
	TryUntilAcquire(duration time.Duration, retry time.Duration, timeout time.Duration) bool
}

type Cache interface {
	Get(key string) (interface{}, error)
	Peek(key string) (interface{}, error)
	Set(key string, value interface{}) error
	Delete(key string) error
	Keys() ([]string, error)
	Purge() error
}
