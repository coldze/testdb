package stores

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"github.com/coldze/testdb/logic/sources"
	"github.com/coldze/testdb/logic/sources/wraps"
	"github.com/coldze/testdb/logic/structs"
	"github.com/coldze/testdb/utils"
)

type redisCache struct {
	ttl   time.Duration
	cache wraps.CacheWrap
}

func buildCacheKey(id string, date time.Time) string {
	return fmt.Sprintf("%s_%v", id, date.UnixNano())
}

func (c *redisCache) Get(ctx context.Context, key structs.Request) (structs.Data, error) {
	logger := utils.GetLogger(ctx)
	res := structs.Data{}
	for _, id := range key.IDs {
		cacheKey := buildCacheKey(id, key.Date)
		data, err := c.cache.Get(cacheKey)
		if err != nil {
			if err == redis.Nil {
				continue
			}
			logger.Warningf("Failed to get for key '%s'. Error: %v", cacheKey, err)
			continue
		}
		if data == nil {
			continue
		}
		dataStr, ok := data.(string)
		if !ok {
			logger.Warningf("Failed to get string data for key '%s'. Actual type: %T", cacheKey, data)
			continue
		}
		count, err := strconv.ParseUint(dataStr, 10, 64)
		if err != nil {
			logger.Warningf("Failed to parse data for key '%s'. Error: %v", cacheKey, err)
			continue
		}
		res[id] = count
	}
	return res, nil
}

func (c *redisCache) Put(ctx context.Context, date time.Time, data structs.Data) error {
	logger := utils.GetLogger(ctx)
	for id, count := range data {
		cacheKey := buildCacheKey(id, date)
		err := c.cache.Set(cacheKey, count, c.ttl)
		if err != nil {
			logger.Warningf("Failed to put for key '%s'. Error: %v", cacheKey, err)
		}
	}
	return nil
}

func (c *redisCache) Del(ctx context.Context, key structs.Request) error {
	logger := utils.GetLogger(ctx)
	for _, id := range key.IDs {
		cacheKey := buildCacheKey(id, key.Date)
		err := c.cache.Del(cacheKey)
		if err != nil {
			logger.Warningf("Failed to remove for key '%s'. Error: %v", cacheKey, err)
		}
	}
	return nil
}

func (c *redisCache) Wipe(ctx context.Context) error {
	return c.cache.Flush()
}

func NewCache(wrap wraps.CacheWrap, ttl time.Duration) sources.Cache {
	return &redisCache{
		cache: wrap,
		ttl:   ttl,
	}
}
