package monitor

import (
	"context"
	"log"
	"sync/atomic"

	"github.com/admpub/godownloader/model"
)

func NewWorkerPool() *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())
	return &WorkerPool{
		ctx:    ctx,
		cancel: cancel,
	}
}

type WorkerPool struct {
	workers    map[string]*MonitoredWorker
	total      int32
	done       int32
	onComplete func(context.Context) error
	state      State
	ctx        context.Context
	cancel     context.CancelFunc
}

func (wp *WorkerPool) AppendWork(iv *MonitoredWorker) {
	if wp.workers == nil {
		wp.workers = make(map[string]*MonitoredWorker)
	}
	iv.ondone = func(ctx context.Context) (err error) {
		doneN := atomic.AddInt32(&wp.done, 1)
		log.Printf("info: complete %d/%d\n", doneN, atomic.LoadInt32(&wp.total))
		if wp.Completed() {
			if wp.onComplete != nil {
				err = wp.onComplete(ctx)
				if err != nil {
					atomic.AddInt32(&wp.done, -1)
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
	if wp.state == Running {
		return nil
	}
	wp.initContext()
	var errs []error
	for _, value := range wp.workers {
		if err := value.Start(wp.ctx); err != nil && err != ErrRunRunningJob {
			errs = append(errs, err)
		}
	}
	wp.state = Running
	return errs
}

func (wp *WorkerPool) StopAll() []error {
	if wp.state == Stopped {
		return nil
	}
	var errs []error
	for _, value := range wp.workers {
		if err := value.Stop(wp.ctx); err != nil && err != ErrStopNonRunningJob {
			errs = append(errs, err)
		}
	}
	wp.state = Stopped
	wp.ctx, wp.cancel = nil, nil
	return errs
}

func (wp *WorkerPool) initContext() {
	if wp.ctx == nil {
		wp.ctx, wp.cancel = context.WithCancel(context.Background())
	}
}

func (wp *WorkerPool) GetAllProgress() []model.DownloadProgress {
	var pr []model.DownloadProgress
	for _, value := range wp.workers {
		pr = append(pr, value.GetProgress())
	}
	return pr
}

func (wp *WorkerPool) ResetAllProgress() {
	for _, value := range wp.workers {
		value.ResetProgress()
	}
	atomic.StoreInt32(&wp.done, 0)
}

func (wp *WorkerPool) State() State {
	return wp.state
}
