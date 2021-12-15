package com

import (
	"context"
	"log"
	"sync"
	"time"
)

func NewWaitGroup(wg *sync.WaitGroup) *WaitGroup {
	return &WaitGroup{
		WaitGroup: wg,
		T:         time.Now().Local(),
	}
}

type WaitGroup struct {
	*sync.WaitGroup
	T time.Time
}

func (wg *WaitGroup) Done() {
	defer func() {
		recover()
	}()
	wg.WaitGroup.Done()
}

type Oncer interface {
	CanSet(reqTag interface{}) bool
	Wait(reqTag interface{})
	Release(reqTag interface{})
	StartGC() error
	Close() error
}

var (
	OnceGCInterval = 5 * time.Minute
	OnceGCLifetime = 2 * time.Minute
)

func NewOnce(gcIntervalAndGCLifetime ...time.Duration) Oncer {
	ctx, cancel := context.WithCancel(context.TODO())
	o := &doOnce{
		data:     map[interface{}]*WaitGroup{},
		interval: OnceGCInterval,
		lifetime: OnceGCLifetime,
		context:  ctx,
		cancel:   cancel,
	}
	if len(gcIntervalAndGCLifetime) > 0 {
		o.interval = gcIntervalAndGCLifetime[0]
	}
	if len(gcIntervalAndGCLifetime) > 1 {
		o.lifetime = gcIntervalAndGCLifetime[1]
	}
	go func() {
		if err := o.StartGC(); err != nil {
			log.Println(err)
		}
	}()
	return o
}

type doOnce struct {
	lock     sync.RWMutex
	data     map[interface{}]*WaitGroup
	interval time.Duration
	lifetime time.Duration
	context  context.Context
	cancel   context.CancelFunc
	canceled bool
}

// CanSet 同一时刻只有一个请求能获取执行权限，获得执行权限的线程接下来需要执行具体的业务逻辑，完成后调用release方法通知其他线程，操作完成，获取资源即可，其他请求接下来需要调用wait方法
// reqTag 请求标识 用于标识同一个资源
func (u *doOnce) CanSet(reqTag interface{}) bool {
	u.lock.Lock()

	if u.data == nil {
		u.data = map[interface{}]*WaitGroup{}
	} else {
		_, ok := u.data[reqTag]
		if ok {
			u.lock.Unlock()
			return false
		}
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	u.data[reqTag] = NewWaitGroup(wg)

	u.lock.Unlock()
	return true
}

// Wait 调用wait方法将处于阻塞状态，直到获得执行权限的线程处理完具体的业务逻辑，调用release方法来通知其他线程资源ok了
func (u *doOnce) Wait(reqTag interface{}) {
	u.lock.RLock()
	w, ok := u.data[reqTag]
	u.lock.RUnlock()
	if !ok {
		return
	}

	w.Wait()
}

// Release 获得执行权限的线程需要在执行完业务逻辑后调用该方法通知其他处于阻塞状态的线程
func (u *doOnce) Release(reqTag interface{}) {
	u.lock.Lock()

	if _, ok := u.data[reqTag]; !ok {
		u.lock.Unlock()
		return
	}
	u.data[reqTag].Done()
	delete(u.data, reqTag)
	u.lock.Unlock()
}

func (u *doOnce) StartGC() error {
	if u.interval == 0 {
		return nil
	}
	t := time.NewTicker(u.interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := u.gc(); err != nil {
				return err
			}
		case <-u.context.Done():
			return nil
		}
	}
}

func (u *doOnce) gc() error {
	for key, val := range u.data {
		if val == nil {
			delete(u.data, key)
			continue
		}
		if u.canceled {
			return nil
		}
		if time.Since(val.T) > u.lifetime {
			val.Done()
			delete(u.data, key)
		}
	}
	return nil
}

func (u *doOnce) clear() error {
	for key, val := range u.data {
		if val == nil {
			delete(u.data, key)
			continue
		}
		val.Done()
		delete(u.data, key)
	}
	return nil
}

func (u *doOnce) Close() error {
	u.cancel()
	u.canceled = true
	u.clear()
	return nil
}
