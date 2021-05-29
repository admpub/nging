package factory

import (
	"time"
)

type Cacher interface {
	Put(key string, value interface{}, lifetime time.Duration) error
	Del(key string) error
	Get(key string) (interface{}, error)
}
