package redis

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/admpub/redistore"
	"github.com/admpub/sessions"
	"github.com/webx-top/echo/middleware/session/engine"
	ss "github.com/webx-top/echo/middleware/session/engine"
)

var DefaultMaxReconnect = 5

func New(opts *RedisOptions) sessions.Store {
	store, err := NewRedisStore(opts)
	if err != nil {
		if !strings.Contains(err.Error(), `connect:`) {
			panic(err.Error())
		}
		retries := opts.MaxReconnect
		if retries <= 0 {
			retries = DefaultMaxReconnect
		}
		for i := 1; i < retries; i++ {
			log.Println(`[sessions]`, err.Error())
			wait := time.Second
			log.Printf(`[sessions] (%d/%d) reconnect redis after %v`, i, retries, wait)
			time.Sleep(wait)
			store, err = NewRedisStore(opts)
			if err == nil {
				log.Println(`[sessions] reconnect redis successfully`)
				return store
			}
		}
		panic(err.Error())
	}
	return store
}

func Reg(store sessions.Store, args ...string) {
	name := `redis`
	if len(args) > 0 {
		name = args[0]
	}
	ss.Reg(name, store)
}

func RegWithOptions(opts *RedisOptions, args ...string) sessions.Store {
	store := New(opts)
	Reg(store, args...)
	return store
}

type RedisOptions struct {
	Size         int      `json:"size"`
	Network      string   `json:"network"`
	Address      string   `json:"address"`
	Password     string   `json:"password"`
	DB           uint     `json:"db"`
	KeyPairs     [][]byte `json:"keyPairs"`
	MaxAge       int      `json:"maxAge"`
	MaxReconnect int      `json:"maxReconnect"`
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
func NewRedisStore(opts *RedisOptions) (sessions.Store, error) {
	store, err := redistore.NewRediStoreWithDB(opts.Size, opts.Network, opts.Address, opts.Password, fmt.Sprintf("%d", opts.DB), opts.KeyPairs...)
	if err != nil {
		return nil, err
	}
	if opts.MaxAge > 0 {
		store.DefaultMaxAge = opts.MaxAge
	} else {
		store.DefaultMaxAge = engine.DefaultMaxAge
	}
	return &redisStore{store}, nil
}

type redisStore struct {
	*redistore.RediStore
}
