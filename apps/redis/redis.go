package redis

import (
	"context"
	"github.com/getevo/evo-ng"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redis_rate/v9"
	"github.com/kelindar/binary"
	"time"
)

var client redis.UniversalClient

var config ConfigWrapper

type ConfigWrapper struct {
	Redis Config `json:"redis" yaml:"redis"`
}
type Config struct {
	AppID    string   `json:"app_id" yaml:"app_id"`
	NodeID   string   `json:"node_id" yaml:"node_id"`
	Username string   `json:"username" yaml:"username"`
	Password string   `json:"password" yaml:"password"`
	Server   []string `json:"server" yaml:"server"`
}

func Register() error {
	evo.ParseConfig(&config)
	return Connect(&config.Redis)
}

func GetConfig() Config {
	return config.Redis
}

func Connect(config *Config) error {
	client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:       config.Server,
		Username:    config.Username,
		Password:    config.Password,
		MaxRetries:  3,
		DialTimeout: time.Duration(2 * time.Second),
	})

	Locker = newLock(client)
	limiter = redis_rate.NewLimiter(client)
	return nil
}

func SetRaw(key string, value interface{}, expiration time.Duration) error {
	key = config.Redis.AppID + "#" + key
	if err := client.Set(context.Background(), key, value, expiration).Err(); err != nil {
		return err
	}
	return nil
}

func Set(key string, value interface{}, expiration time.Duration) error {
	key = config.Redis.AppID + "#" + key
	b, err := binary.Marshal(value)
	if err != nil {
		return err
	}
	if err := client.Set(context.Background(), key, b, expiration).Err(); err != nil {
		return err
	}

	return nil
}

func Search(key string) []string {
	key = config.Redis.AppID + "#" + key
	var v []string
	result := client.Scan(context.Background(), 0, key, 0).Iterator()
	for result.Next(context.Background()) {
		v = append(v, result.Val())
	}
	return v
}

func GetString(key string) (string, error) {
	key = config.Redis.AppID + "#" + key
	b, err := GetRaw(key)
	return string(b), err
}

func GetRaw(key string) ([]byte, error) {
	key = config.Redis.AppID + "#" + key
	b, err := client.Get(context.Background(), key).Bytes()
	if err != nil {
		return b, err
	}
	return b, nil

}

func Exists(key string) bool {
	key = config.Redis.AppID + "#" + key
	return client.Exists(context.Background(), key).Err() != nil
}

func Extend(key string, duration time.Duration) error {
	key = config.Redis.AppID + "#" + key
	return client.Expire(context.Background(), key, duration).Err()
}

func Del(key string) error {
	key = config.Redis.AppID + "#" + key
	return client.Del(context.Background(), key).Err()
}

func Decr(key string, step int) error {
	key = config.Redis.AppID + "#" + key
	return client.DecrBy(context.Background(), key, int64(step)).Err()
}

func Incr(key string, step int) error {
	key = config.Redis.AppID + "#" + key
	return client.IncrBy(context.Background(), key, int64(step)).Err()
}

func Get(key string, v interface{}) error {
	key = config.Redis.AppID + "#" + key
	b, err := client.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	return binary.Unmarshal(b, v)
}

func GetDel(key string, v interface{}) error {
	key = config.Redis.AppID + "#" + key
	b, err := client.GetDel(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	return binary.Unmarshal(b, v)
}

func GetEx(key string, v interface{}, duration time.Duration) error {
	key = config.Redis.AppID + "#" + key
	b, err := client.GetEx(context.Background(), key, duration).Bytes()
	if err != nil {
		return err
	}
	return binary.Unmarshal(b, v)
}

func Delete(keys ...string) error {
	for i, key := range keys {
		keys[i] = config.Redis.AppID + "#" + key
	}

	return client.Del(context.Background(), keys...).Err()
}

func Publish(channel string, payload interface{}) {
	channel = config.Redis.AppID + "#" + channel
	client.Publish(context.Background(), channel, payload)
}

func Client() redis.UniversalClient {
	return client
}
