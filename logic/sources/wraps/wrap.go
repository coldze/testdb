package wraps

import "time"

type DbWrap interface {
	Closable
	Query(q string, args ...interface{}) (Scanner, error)
}

type CacheWrap interface {
	Set(key string, data interface{}, ttl time.Duration) error
	Del(key string) error
	Get(key string) (interface{}, error)
	Flush() error
	Close() error
}
