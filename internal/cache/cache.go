package cache

import "time"

type Cache interface {
	Set(key string, value interface{}, expire time.Duration) error
	Extend(key string, duration time.Duration) error
	Replace(key string, value interface{}, duration time.Duration)
	Delete(key string) error
	Get(key string, out interface{}) error
	Flush()
	Keys() []string
}
