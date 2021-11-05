package factory

type Cacher interface {
	Put(key string, value interface{}, ttlSeconds int64) error
	Del(key string) error
	Get(key string) (interface{}, error)
	Do(key string, recv interface{}, fn func() error, ttlSeconds int64) error
}
