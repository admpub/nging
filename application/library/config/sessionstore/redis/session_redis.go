package redis

import (
	"reflect"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/echo/middleware/session/engine"
	"github.com/webx-top/echo/middleware/session/engine/cookie"
	"github.com/webx-top/echo/middleware/session/engine/redis"
	"github.com/webx-top/echo/param"
)

func init() {
	config.RegisterSessionStore(`redis`, `Redis存储`, initSessionStoreRedis)
}

var sessionStoreRedisOptions *redis.RedisOptions

func initSessionStoreRedis(_ *config.Config, cookieOptions *cookie.CookieOptions, sessionConfig param.Store) (changed bool, err error) {
	redisOptions := &redis.RedisOptions{
		Size:         sessionConfig.Int(`maxIdle`),
		Network:      sessionConfig.String(`network`),
		Address:      sessionConfig.String(`address`),
		Password:     sessionConfig.String(`password`),
		DB:           sessionConfig.Uint(`db`),
		KeyPairs:     cookieOptions.KeyPairs,
		MaxAge:       sessionConfig.Int(`maxAge`),
		MaxReconnect: sessionConfig.Int(`maxReconnect`),
	}
	if redisOptions.Size <= 0 {
		redisOptions.Size = 10
	}
	if len(redisOptions.Network) == 0 {
		redisOptions.Network = `tcp`
	}
	if len(redisOptions.Address) == 0 {
		redisOptions.Address = `127.0.0.1:6379`
	}
	if redisOptions.MaxReconnect <= 0 {
		redisOptions.MaxReconnect = 30
	}
	if sessionStoreRedisOptions == nil || !engine.Exists(`redis`) || !reflect.DeepEqual(redisOptions, sessionStoreRedisOptions) {
		redis.RegWithOptions(redisOptions)
		sessionStoreRedisOptions = redisOptions
		changed = true
	}
	return
}
