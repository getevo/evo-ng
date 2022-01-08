package redis

import (
	"context"
	"fmt"
	"github.com/getevo/evo-ng"
	"github.com/go-redis/redis/v8"
	"github.com/kelindar/binary"
	"time"
)

var clusterClient *redis.ClusterClient
var singleClient *redis.Client
var ctx = context.Background()
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

	if len(config.Server) > 1 {
		clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:       config.Server,
			Username:    config.Username,
			Password:    config.Password,
			MaxRetries:  3,
			DialTimeout: time.Duration(2 * time.Second),
		})
		if err := clusterClient.Ping(ctx).Err(); err != nil {
			return err
		}
		clusterClient.ReloadState(ctx)
	} else if len(config.Server) == 1 {
		if len(config.Server) == 1 {
			singleClient = redis.NewClient(&redis.Options{
				Addr:        config.Server[0],
				Username:    config.Username,
				Password:    config.Password,
				MaxRetries:  3,
				DialTimeout: time.Duration(2 * time.Second),
			})
			if err := singleClient.Ping(ctx).Err(); err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("invalid server")
	}
	return nil
}

func Set(key string, value interface{}, expiration time.Duration) error {
	key = config.Redis.AppID + "#" + key
	b, err := binary.Marshal(value)
	if err != nil {
		return err
	}
	if clusterClient != nil {
		if err := clusterClient.Set(context.Background(), key, b, expiration).Err(); err != nil {
			return err
		}
	} else {
		if err := singleClient.Set(context.Background(), key, b, expiration).Err(); err != nil {
			return err
		}
	}

	return nil
}

func Search(key string) []string {
	key = config.Redis.AppID + "#" + key
	var v []string
	if clusterClient != nil {
		result := clusterClient.Scan(context.Background(), 0, key, 0).Iterator()
		for result.Next(context.Background()) {
			v = append(v, result.Val())
		}
	} else {
		result := singleClient.Scan(context.Background(), 0, key, 0).Iterator()
		for result.Next(context.Background()) {
			v = append(v, result.Val())
		}
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
	if clusterClient != nil {
		b, err := clusterClient.Get(context.Background(), key).Bytes()
		if err != nil {
			return b, err
		}
		return b, nil
	} else {
		b, err := singleClient.Get(context.Background(), key).Bytes()
		if err != nil {
			return b, err
		}
		return b, nil
	}

}

func Del(key string) error {
	key = config.Redis.AppID + "#" + key
	if clusterClient != nil {
		return clusterClient.Del(context.Background(), key).Err()
	} else {
		return singleClient.Del(context.Background(), key).Err()
	}
}

func Decr(key string, step int) error {
	key = config.Redis.AppID + "#" + key
	if clusterClient != nil {
		return clusterClient.DecrBy(context.Background(), key, int64(step)).Err()
	} else {
		return singleClient.DecrBy(context.Background(), key, int64(step)).Err()
	}
}

func Incr(key string, step int) error {
	key = config.Redis.AppID + "#" + key
	if clusterClient != nil {
		return clusterClient.IncrBy(context.Background(), key, int64(step)).Err()
	} else {
		return singleClient.IncrBy(context.Background(), key, int64(step)).Err()
	}
}

func Get(key string, v interface{}) error {
	key = config.Redis.AppID + "#" + key
	if clusterClient != nil {
		b, err := clusterClient.Get(context.Background(), key).Bytes()
		if err != nil {
			return err
		}
		return binary.Unmarshal(b, v)
	} else {
		b, err := singleClient.Get(context.Background(), key).Bytes()
		if err != nil {
			return err
		}
		return binary.Unmarshal(b, v)
	}
}

func Delete(key string) error {
	key = config.Redis.AppID + "#" + key
	if clusterClient != nil {
		return clusterClient.Del(context.Background(), key).Err()
	} else {
		return singleClient.Del(context.Background(), key).Err()
	}
}

func Publish(channel string, payload interface{}) {
	channel = config.Redis.AppID + "#" + channel
	if clusterClient != nil {
		clusterClient.Publish(context.Background(), channel, payload)
	} else {
		singleClient.Publish(context.Background(), channel, payload)
	}
}
