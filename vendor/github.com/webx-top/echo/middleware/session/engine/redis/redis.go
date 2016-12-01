package redis

import (
	"github.com/admpub/redistore"
	"github.com/admpub/sessions"
	"github.com/webx-top/echo"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

func New(opts *RedisOptions) RedisStore {
	store, err := NewRedisStore(opts)
	if err != nil {
		panic(err.Error())
	}
	store.Options(*opts.SessionOptions)
	return store
}

func Reg(store RedisStore, args ...string) {
	name := `redis`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *RedisOptions, args ...string) {
	Reg(New(opts), args...)
}

type RedisStore interface {
	ss.Store
}

type RedisOptions struct {
	Size           int                  `json:"size"`
	Network        string               `json:"network"`
	Address        string               `json:"address"`
	Password       string               `json:"password"`
	KeyPairs       [][]byte             `json:"keyPairs"`
	SessionOptions *echo.SessionOptions `json:"session"`
}

// size: maximum number of idle connections.
// network: tcp or udp
// address: host:port
// password: redis-password
// Keys are defined in pairs to allow key rotation, but the common case is to set a single
// authentication key and optionally an encryption key.
//
// The first key in a pair is used for authentication and the second for encryption. The
// encryption key can be set to nil or omitted in the last pair, but the authentication key
// is required in all pairs.
//
// It is recommended to use an authentication key with 32 or 64 bytes. The encryption key,
// if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes.
func NewRedisStore(opts *RedisOptions) (RedisStore, error) {
	store, err := redistore.NewRediStore(opts.Size, opts.Network, opts.Address, opts.Password, opts.KeyPairs...)
	if err != nil {
		return nil, err
	}
	return &redisStore{store}, nil
}

type redisStore struct {
	*redistore.RediStore
}

func (c *redisStore) Options(options echo.SessionOptions) {
	c.RediStore.Options = &sessions.Options{
		Path:     options.Path,
		Domain:   options.Domain,
		MaxAge:   options.MaxAge,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
	}
}
