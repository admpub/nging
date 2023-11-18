package notice

import (
	"context"
	"time"
)

func SetUser(user string) func(*HTTPNoticerConfig) {
	return func(cfg *HTTPNoticerConfig) {
		cfg.User = user
	}
}

func SetClientID(clientID string) func(*HTTPNoticerConfig) {
	return func(cfg *HTTPNoticerConfig) {
		cfg.ClientID = clientID
	}
}

func SetType(typ string) func(*HTTPNoticerConfig) {
	return func(cfg *HTTPNoticerConfig) {
		cfg.Type = typ
	}
}

func SetID(id interface{}) func(*HTTPNoticerConfig) {
	return func(cfg *HTTPNoticerConfig) {
		cfg.ID = id
	}
}

func SetMode(mode string) func(*HTTPNoticerConfig) {
	return func(cfg *HTTPNoticerConfig) {
		cfg.Mode = mode
	}
}

func SetTimeout(timeout time.Duration) func(*HTTPNoticerConfig) {
	return func(cfg *HTTPNoticerConfig) {
		cfg.Timeout = timeout
	}
}

func SetIsExited(isExited IsExited) func(*HTTPNoticerConfig) {
	return func(cfg *HTTPNoticerConfig) {
		cfg.IsExited = isExited
	}
}

type HTTPNoticerConfig struct {
	User     string
	Type     string // Topic
	ClientID string
	ID       interface{}
	IsExited IsExited
	Timeout  time.Duration
	Mode     string // element / notify
}

func NewHTTPNoticerConfig() *HTTPNoticerConfig {
	return &HTTPNoticerConfig{}
}

func (c *HTTPNoticerConfig) SetUser(user string) *HTTPNoticerConfig {
	c.User = user
	return c
}

func (c *HTTPNoticerConfig) SetType(typ string) *HTTPNoticerConfig {
	c.Type = typ
	return c
}

func (c *HTTPNoticerConfig) SetClientID(clientID string) *HTTPNoticerConfig {
	c.ClientID = clientID
	return c
}

func (c *HTTPNoticerConfig) SetID(id interface{}) *HTTPNoticerConfig {
	c.ID = id
	return c
}

func (c *HTTPNoticerConfig) SetTimeout(t time.Duration) *HTTPNoticerConfig {
	c.Timeout = t
	return c
}

func (c *HTTPNoticerConfig) SetIsExited(isExited IsExited) *HTTPNoticerConfig {
	c.IsExited = isExited
	return c
}

func (c *HTTPNoticerConfig) SetMode(mode string) *HTTPNoticerConfig {
	c.Mode = mode
	return c
}

func (c *HTTPNoticerConfig) Noticer(ctx context.Context) Noticer {
	return NewNoticer(ctx, c)
}
