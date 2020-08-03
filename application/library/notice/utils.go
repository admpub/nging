package notice

import (
	"context"
	"time"
)

type HTTPNoticerConfig struct {
	User     string
	Type     string
	ClientID string
	ID       interface{}
	IsExited IsExited
	Timeout  time.Duration
	Mode     string
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

func NewControlWithContext(ctx context.Context, timeout time.Duration) IsExited {
	defaultCtrl := &Control{}
	defaultCtrl.ListenContextAndTimeout(ctx, timeout)
	return defaultCtrl
}

func NewNoticer(ctx context.Context, config *HTTPNoticerConfig) Noticer {
	var noticeSender Noticer
	if config.IsExited == nil && config.Timeout != 0 {
		config.IsExited = NewControlWithContext(ctx, config.Timeout)
	}
	if len(config.Mode) == 0 {
		if config.ID != nil {
			config.Mode = `element`
		} else {
			config.Mode = `notify`
		}
	}
	progress := NewProgress().SetControl(config.IsExited)
	if len(config.User) > 0 {
		OpenMessage(config.User, config.Type)
		//defer CloseMessage(config.User, config.Type)
		noticeSender = func(message interface{}, statusCode int, progs ...*Progress) error {
			msg := NewMessageWithValue(
				config.Type,
				``,
				message,
				statusCode,
			).SetMode(config.Mode).SetID(config.ID)
			var prog *Progress
			if len(progs) > 0 {
				prog = progs[0]
			}
			if prog == nil {
				prog = progress
			}
			msg.SetProgress(prog).CalcPercent().SetClientID(config.ClientID)
			sendErr := Send(config.User, msg)
			return sendErr
		}
	} else {
		noticeSender = DefaultNoticer
	}
	return noticeSender
}
