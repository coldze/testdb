package sources

import (
	"context"

	"github.com/coldze/testdb/logic/structs"
	"github.com/coldze/testdb/utils"
)

type cached struct {
	original Source
	cache    Cache
}

func (c *cached) getFromCache(ctx context.Context, key structs.Request) (res structs.Data, absent *structs.Request, err error) {
	absent = &key
	res = structs.Data{}

	if key.IgnoreCache {
		return
	}
	fromCache, err := c.cache.Get(ctx, key)
	if err != nil {
		return
	}
	if fromCache == nil {
		return
	}
	res = fromCache
	absentKeys := []string{}
	for _, k := range key.IDs {
		_, ok := res[k]
		if !ok {
			absentKeys = append(absentKeys, k)
		}
	}
	if len(absentKeys) <= 0 {
		absent = nil
	} else {
		absent = &structs.Request{
			IDs:         absentKeys,
			Date:        key.Date,
			IgnoreCache: key.IgnoreCache,
		}
	}
	return
}

func (c *cached) Get(ctx context.Context, key structs.Request) (structs.Data, error) {
	logger := utils.GetLogger(ctx)

	res, absentKey, err := c.getFromCache(ctx, key)
	if err != nil {
		logger.Warningf("Failed to get data from cache. Error: %v", err)
	}
	logger.Infof("In cache: %+v", res)
	if absentKey == nil {
		logger.Infof("Full cache hit. Res: %+v", res)
		return res, nil
	}
	absent, err := c.original.Get(ctx, *absentKey)
	if err != nil {
		return nil, err
	}
	err = c.cache.Put(ctx, key.Date, absent)
	if err != nil {
		logger.Warningf("Failed to put data to cache. Error: %v", err)
	}
	logger.Infof("Not in cache: %+v", absent)
	for k, v := range absent {
		res[k] = v
	}
	logger.Infof("Mixed response. Res: %+v", res)
	return res, nil
}

func NewSourceWithCache(source Source, cache Cache) Source {
	return &cached{
		cache:    cache,
		original: source,
	}
}
