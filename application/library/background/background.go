package background

import (
	"context"
	"sync"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// Backgrounds 后台任务集合
var Backgrounds = sync.Map{}

// Exec 执行信息
type Exec struct {
	m  map[string]*Background
	mu sync.RWMutex
}

// Cancel 取消某个任务
func (e *Exec) Cancel(cacheKey string) {
	e.mu.Lock()
	e.cancel(cacheKey)
	e.mu.Unlock()
}

func (e *Exec) cancel(cacheKey string) {
	if bgExec, ok := (*e).m[cacheKey]; ok {
		bgExec.Cancel()()
		delete((*e).m, cacheKey)
	}
}

// Exists 任务是否存在
func (e *Exec) Exists(cacheKey string) bool {
	e.mu.RLock()
	_, ok := (*e).m[cacheKey]
	e.mu.RUnlock()
	return ok
}

// Map 任务列表
func (e *Exec) Map() map[string]*Background {
	e.mu.RLock()
	r := (*e).m
	e.mu.RUnlock()
	return r
}

// Add 新增任务
func (e *Exec) Add(op string, cacheKey string, bgExec *Background) {
	e.mu.Lock()

	e.cancel(cacheKey) // 避免被覆盖后旧任务失去控制，先取消已存在的任务

	(*e).m[cacheKey] = bgExec

	e.mu.Unlock()

	Backgrounds.Store(op, e)
}

// All 所有任务
func All() map[string]map[string]*Background {
	r := map[string]map[string]*Background{}
	Backgrounds.Range(func(key, val interface{}) bool {
		r[param.AsString(key)] = val.(*Exec).Map()
		return true
	})
	return r
}

// ListBy 获取某个操作的所有任务
func ListBy(op string) *Exec {
	old, exists := Backgrounds.Load(op)
	if !exists {
		return nil
	}
	exec := old.(*Exec)
	return exec
}

// Cancel 取消执行
func Cancel(op string, cacheKey string) error {
	exec := ListBy(op)
	if exec == nil {
		return nil
	}
	exec.Cancel(cacheKey)
	Backgrounds.Store(op, exec)
	return nil
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
func (b *Background) Cancel() context.CancelFunc {
	return b.cancel
}

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
