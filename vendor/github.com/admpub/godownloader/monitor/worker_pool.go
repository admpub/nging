package monitor

import (
	"context"
	"log"
	"sync/atomic"
)

type WorkerPool struct {
	workers    map[string]*MonitoredWorker
	total      int32
	done       int32
	onComplete func(context.Context) error
	state      State
}

func (wp *WorkerPool) AppendWork(iv *MonitoredWorker) {
	if wp.workers == nil {
		wp.workers = make(map[string]*MonitoredWorker)
	}
	iv.ondone = func(ctx context.Context) (err error) {
		atomic.AddInt32(&wp.done, 1)
		log.Printf("info: complete %d/%d", wp.done, wp.total)
		if wp.Completed() {
			if wp.onComplete != nil {
				err = wp.onComplete(ctx)
				if err != nil {
					return
				}
			}
			wp.state = Completed
		}
		return
	}
	wp.workers[iv.GetId()] = iv
	atomic.AddInt32(&wp.total, 1)
}

func (wp *WorkerPool) AfterComplete(fn func(context.Context) error) {
	wp.onComplete = fn
}

func (wp *WorkerPool) Completed() bool {
	return atomic.LoadInt32(&wp.total) == atomic.LoadInt32(&wp.done)
}

func (wp *WorkerPool) StartAll() []error {
	var errs []error
	for _, value := range wp.workers {
		if err := value.Start(); err != nil && err != ErrRunRunningJob {
			errs = append(errs, err)
		}
	}
	wp.state = Running
	return errs
}

func (wp *WorkerPool) StopAll() []error {
	var errs []error
	for _, value := range wp.workers {
		if err := value.Stop(); err != nil && err != ErrStopNonRunningJob {
			errs = append(errs, err)
		}
	}
	wp.state = Stopped
	return errs
}

func (wp *WorkerPool) GetAllProgress() interface{} {
	var pr []interface{}
	for _, value := range wp.workers {
		pr = append(pr, value.GetProgress())
	}
	return pr
}

func (wp *WorkerPool) State() State {
	return wp.state
}
