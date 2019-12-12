package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"

	"github.com/coldze/testdb/logic/sources/wraps"
)

type redisCache struct {
	client *redis.Client
}

func (r *redisCache) Set(key string, data interface{}, ttl time.Duration) error {
	return r.client.Set(key, data, ttl).Err()
}

func (r *redisCache) Del(key string) error {
	return r.client.Del(key).Err()
}

func (r *redisCache) Get(key string) (interface{}, error) {
	return r.client.Get(key).Result()
}

func (r *redisCache) Flush() error {
	return r.client.FlushDB().Err()
}

func (r *redisCache) Close() error {
	return r.client.Close()
}

func NewRedis(cfg *redis.Options) (wraps.CacheWrap, error) {
	client := redis.NewClient(cfg)
	if client == nil {
		return nil, errors.New("internal error - redis client is nil")
	}
	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return &redisCache{
		client: client,
	}, nil
}
