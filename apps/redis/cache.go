package redis

import (
	"errors"
	"time"
)

var (
	ErrKeyNotExists = errors.New("key does not exists")
)

var Cache = RedisCache{}

type RedisCache struct{}

const cacheSuffix = "::cache"

func (c RedisCache) Set(key string, value interface{}, expire time.Duration) error {
	return Set(key+cacheSuffix, value, expire)
}

func (c RedisCache) Extend(key string, expire time.Duration) error {
	return Extend(key+cacheSuffix, expire)
}

func (c RedisCache) Replace(key string, value interface{}, expire time.Duration) error {
	if Exists(key) {
		return Set(key+cacheSuffix, value, expire)
	}
	return ErrKeyNotExists
}

func (c RedisCache) Delete(key string) error {
	return Delete(key + cacheSuffix)
}

func (c RedisCache) Get(key string, out interface{}) error {
	return Get(key+cacheSuffix, out)
}

func (c RedisCache) Flush() error {
	var keys = Search("*" + cacheSuffix)
	if clusterClient != nil {
		return clusterClient.Del(ctx, keys...).Err()
	} else {
		return singleClient.Del(ctx, keys...).Err()
	}
}

func (c RedisCache) Keys() []string {
	var keys = Search("*" + cacheSuffix)
	return keys
}
