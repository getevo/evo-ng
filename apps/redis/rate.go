package redis

import (
	"github.com/go-redis/redis_rate/v9"
)

var limiter *redis_rate.Limiter

type RateLimit struct {
}
type Limit redis_rate.Limit
type Result redis_rate.Result

var RateLimiter RateLimit

func (l RateLimit) AllowAtMost(key string, limit Limit, n int) (Result, error) {
	var result, err = limiter.AllowAtMost(ctx, config.Redis.AppID+key, redis_rate.Limit(limit), n)
	return Result(*result), err
}

func (l RateLimit) AllowN(key string, limit Limit, n int) (Result, error) {
	var result, err = limiter.AllowN(ctx, config.Redis.AppID+key, redis_rate.Limit(limit), n)
	return Result(*result), err
}

func (l RateLimit) Allow(key string, limit Limit, n int) (Result, error) {
	var result, err = limiter.Allow(ctx, config.Redis.AppID+key, redis_rate.Limit(limit))
	return Result(*result), err
}

func (l RateLimit) Reset(key string) error {
	return limiter.Reset(ctx, config.Redis.AppID+key)
}
