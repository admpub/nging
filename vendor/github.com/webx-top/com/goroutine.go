package com

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

var ErrExitedByContext = errors.New(`received an exit notification from the context`)

func Loop(ctx context.Context, exec func() error, duration time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	check := func() <-chan struct{} {
		return ctx.Done()
	}
	for {
		select {
		case <-check():
			log.Println(CalledAtFileLine(2), ErrExitedByContext)
			return ErrExitedByContext
		default:
			if err := exec(); err != nil {
				return err
			}
			time.Sleep(duration)
		}
	}
}

// Notify 等待系统信号
// <-Notify()
func Notify(sig ...os.Signal) chan os.Signal {
	terminate := make(chan os.Signal, 1)
	if len(sig) == 0 {
		sig = []os.Signal{os.Interrupt}
	}
	signal.Notify(terminate, sig...)
	return terminate
}

func NewDelayOnce(delay time.Duration, timeout time.Duration) *DelayOnce {
	if timeout <= delay {
		panic(`timeout must be greater than delay`)
	}
	return &DelayOnce{
		mp:      sync.Map{},
		delay:   delay,
		timeout: timeout,
	}
}

// DelayOnce 触发之后延迟一定的时间后再执行。如果在延迟处理的时间段内再次触发，则延迟时间基于此处触发时间顺延
// d := NewDelayOnce(time.Second*5, time.Hour)
// ctx := context.TODO()
// for i:=0; i<10; i++ {
// 	d.Do(ctx, `key`,func() error { return nil  })
// }
type DelayOnce struct {
	mp      sync.Map
	delay   time.Duration
	timeout time.Duration
}

type eventSession struct {
	cancel context.CancelFunc
	time   time.Time
	mutex  sync.RWMutex
}

func (e *eventSession) Renew(t time.Time) {
	e.mutex.Lock()
	e.time = t
	e.mutex.Unlock()
}

func (e *eventSession) Time() time.Time {
	e.mutex.RLock()
	t := e.time
	e.mutex.RUnlock()
	return t
}

func (d *DelayOnce) checkAndStore(parentCtx context.Context, key string) (exit bool, ctx context.Context) {
	v, loaded := d.mp.Load(key)
	if loaded {
		session := v.(*eventSession)
		if time.Since(session.Time()) < d.timeout { // 超过 d.timeout 后重新处理，d.timeout 内记录当前时间
			session.Renew(time.Now())
			d.mp.Store(key, session)
			return true, nil
		}
		session.cancel()
	}
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(parentCtx)
	d.mp.Store(key, &eventSession{
		cancel: cancel,
		time:   time.Now(),
	})
	return false, ctx
}

func (d *DelayOnce) Do(parentCtx context.Context, key string, f func() error) (isNew bool) {
	exit, ctx := d.checkAndStore(parentCtx, key)
	if exit {
		return false
	}
	go func(key string) {
		for {
			t := time.NewTicker(time.Second)
			defer t.Stop()
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				if err := d.exec(key, f); err != nil {
					log.Println(key+`:`, err)
				}
			}
		}
	}(key)
	return true
}

func (d *DelayOnce) exec(key string, f func() error) (err error) {
	v, ok := d.mp.Load(key)
	if !ok {
		return
	}
	session := v.(*eventSession)
	if time.Since(session.Time()) > d.delay { // 时间超过d.delay才触发
		err = f()
		if err != nil {
			return
		}
		d.mp.Delete(key)
		session.cancel()
	}
	return
}
