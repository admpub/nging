package factory

import "context"

type Cacher interface {
	Put(ctx context.Context, key string, value interface{}, ttlSeconds int64) error
	Del(ctx context.Context, key string) error
	Get(ctx context.Context, key string, recv interface{}) error
	Do(ctx context.Context, key string, recv interface{}, fn func() error, ttlSeconds int64) error
}
