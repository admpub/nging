package notice

import (
	"context"
	"time"
)

func NewControlWithContext(ctx context.Context, timeout time.Duration) IsExited {
	defaultCtrl := &Control{}
	defaultCtrl.ListenContextAndTimeout(ctx, timeout)
	return defaultCtrl
}

type Control struct {
	exited bool
}

func (c *Control) IsExited() bool {
	return c.exited
}

func (c *Control) Exited() *Control {
	c.exited = true
	return c
}

func (c *Control) ListenContextAndTimeout(ctx context.Context, timeouts ...time.Duration) *Control {
	timeout := 24 * time.Hour
	if len(timeouts) > 0 && timeouts[0] != 0 {
		timeout = timeouts[0]
	}
	t := time.NewTicker(timeout)
	defer t.Stop()
	go func() {
		for {
			select {
			case <-ctx.Done():
				c.Exited()
				return
			case <-t.C:
				c.Exited()
				return
			}
		}
	}()
	return c
}

type IsExited interface {
	IsExited() bool
}
