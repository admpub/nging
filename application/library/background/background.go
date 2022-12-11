package background

import (
	"context"
	"time"

	"github.com/webx-top/echo"
)

// New 新建后台执行信息
func New(c context.Context, opt echo.H) *Background {
	if c == nil {
		c = context.Background()
	}
	if opt == nil {
		opt = echo.H{}
	}
	ctx, cancel := context.WithCancel(c)
	return &Background{
		ctx:     ctx,
		cancel:  cancel,
		Options: opt,
		Started: time.Now(),
	}
}

// Background 后台执行信息
type Background struct {
	ctx     context.Context
	cancel  context.CancelFunc
	Options echo.H
	Started time.Time
}

// Context 暂存上下文信息
func (b *Background) Context() context.Context {
	return b.ctx
}

// Cancel 取消执行
func (b *Background) Cancel() {
	b.cancel()
}
