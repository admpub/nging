package once

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

// Do will execute the function and will make the Once as done.
// All subsequent calls to Do will not execute the function unless Reset() is called on Once.
func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) == 0 {
		o.doSlow(f)
	}
}

// DoForce will execute the function even if Do has already been called on Once.
// It will not call f directly because that will not be a thread safe call.
func (o *Once) DoForce(f func()) {
	atomic.StoreUint32(&o.done, 0)
	o.doSlow(f)
}

// Reset will reset the Once.
// The next time Do is called the function will execute
func (o *Once) Reset() {
	if atomic.LoadUint32(&o.done) == 1 {
		atomic.StoreUint32(&o.done, 0)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}
