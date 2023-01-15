package com

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

var ErrExitedByContext = context.Canceled

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

func NewDelayOnce(delay time.Duration, timeout time.Duration, debugMode ...bool) *DelayOnce {
	if timeout <= delay {
		panic(`timeout must be greater than delay`)
	}
	var debug bool
	if len(debugMode) > 0 {
		debug = debugMode[0]
	}
	return &DelayOnce{
		mp:      sync.Map{},
		delay:   delay,
		timeout: timeout,
		debug:   debug,
	}
}

// DelayOnce 触发之后延迟一定的时间后再执行。如果在延迟处理的时间段内再次触发，则延迟时间基于此处触发时间顺延
// d := NewDelayOnce(time.Second*5, time.Hour)
// ctx := context.TODO()
//
//	for i:=0; i<10; i++ {
//		d.Do(ctx, `key`,func() error { return nil  })
//	}
type DelayOnce struct {
	mp      sync.Map
	delay   time.Duration
	timeout time.Duration
	debug   bool
}

type eventSession struct {
	cancel context.CancelFunc
	time   time.Time
	mutex  sync.RWMutex
	stop   chan struct{}
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

func (e *eventSession) Cancel() <-chan struct{} {
	e.cancel()
	return e.stop
}

func (d *DelayOnce) checkAndStore(parentCtx context.Context, key string) (*eventSession, context.Context) {
	v, loaded := d.mp.Load(key)
	if loaded {
		session := v.(*eventSession)
		if time.Since(session.Time()) < d.timeout { // 超过 d.timeout 后重新处理，d.timeout 内记录当前时间
			session.Renew(time.Now())
			d.mp.Store(key, session)
			return nil, nil
		}

		if d.debug {
			log.Println(`[DelayOnce] cancel -------------> ` + key)
		}

		<-session.Cancel()

		if d.debug {
			log.Println(`[DelayOnce] canceled -------------> ` + key)
		}
	}
	ctx, cancel := context.WithCancel(parentCtx)
	session := &eventSession{
		cancel: cancel,
		time:   time.Now(),
		stop:   make(chan struct{}, 1),
	}
	d.mp.Store(key, session)
	return session, ctx
}

func (d *DelayOnce) Do(parentCtx context.Context, key string, f func() error) (isNew bool) {
	session, ctx := d.checkAndStore(parentCtx, key)
	if session == nil {
		return false
	}
	go func(key string) {
		for {
			t := time.NewTicker(time.Second)
			defer t.Stop()
			select {
			case <-ctx.Done(): // 如果先进入“<-t.C”分支，会等“<-t.C”分支内的代码执行完毕后才有机会执行本分支
				d.mp.Delete(key)
				session.stop <- struct{}{}
				close(session.stop)
				if d.debug {
					log.Println(`[DelayOnce] close -------------> ` + key)
				}
				return
			case <-t.C:
				if time.Since(session.Time()) > d.delay { // 时间超过d.delay才触发
					err := f()
					session.Cancel()
					if err != nil {
						log.Println(key+`:`, err)
					}
				}
			}
		}
	}(key)
	return true
}

func (d *DelayOnce) DoWithState(parentCtx context.Context, key string, f func(func() bool) error) (isNew bool) {
	session, ctx := d.checkAndStore(parentCtx, key)
	if session == nil {
		return false
	}
	go func(key string) {
		var state int32
		isAbort := func() bool {
			return atomic.LoadInt32(&state) > 0
		}
		go func() {
			<-ctx.Done()
			atomic.AddInt32(&state, 1)
		}()
		for {
			t := time.NewTicker(time.Second)
			defer t.Stop()
			select {
			case <-ctx.Done(): // 如果先进入“<-t.C”分支，会等“<-t.C”分支内的代码执行完毕后才有机会执行本分支
				d.mp.Delete(key)
				session.stop <- struct{}{}
				close(session.stop)
				if d.debug {
					log.Println(`[DelayOnce] close -------------> ` + key)
				}
				return
			case <-t.C:
				if time.Since(session.Time()) > d.delay { // 时间超过d.delay才触发
					err := f(isAbort)
					session.Cancel()
					if err != nil {
						log.Println(key+`:`, err)
					}
				}
			}
		}
	}(key)
	return true
}
