package monitor

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/admpub/godownloader/model"
)

var States = map[State]string{
	Stopped:   `Stopped`,
	Running:   `Running`,
	Failed:    `Failed`,
	Completed: `Completed`,
}

type State int

func (s State) String() string {
	return States[s]
}

func (s State) Int() int {
	return int(s)
}

const (
	Stopped State = iota
	Running
	Failed
	Completed
)

type DiscretWork interface {
	DoWork(context.Context) (bool, error)
	GetProgress() model.DownloadProgress
	BeforeRun(context.Context) error
	AfterStop() error
	IsPartialDownload() bool
	ResetProgress()
}

func genUid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type MonitoredWorker struct {
	lc         sync.Mutex
	Itw        DiscretWork
	wgrun      sync.WaitGroup
	guid       string
	state      State
	stateLock  sync.RWMutex
	ondone     func(context.Context) error
	ctx        context.Context
	cancelFunc context.CancelFunc
	id         atomic.Value
}

func (mw *MonitoredWorker) setState(state State) {
	mw.stateLock.Lock()
	mw.state = state
	mw.stateLock.Unlock()
}

func (mw *MonitoredWorker) doWorkExec(ctx context.Context, id interface{}) (bool, error) {
	isdone, err := mw.Itw.DoWork(ctx)
	if err != nil {
		log.Println("error: guid", mw.guid, "work failed", err)
		mw.setState(Failed)
		return isdone, err
	}
	if isdone {
		if err = mw.onDoneExec(ctx); err != nil {
			return isdone, err
		}
		mw.setState(Completed)
		log.Println("info: work done")
		return isdone, err
	}
	if mw.id.Load() != id {
		return false, nil
	}
	return mw.doWorkExec(ctx, id)
}

func (mw *MonitoredWorker) wgoroute() {
	log.Println("info: work start", mw.GetId())
	id := mw.id.Load()
	defer func() {
		mw.wgrun.Done()
	}()

	done := make(chan struct{})
	go func() {
		mw.doWorkExec(mw.ctx, id)
		done <- struct{}{}
		close(done)
	}()

	for {
		select {
		case <-mw.ctx.Done():
			log.Println("info: stop work guid", mw.GetId())
			return
		case <-done:
			mw.cancelFunc()
			log.Println("info: release work guid", mw.GetId())
			return
		}
	}
}

func (mw *MonitoredWorker) onDoneExec(ctx context.Context) error {
	if mw.ondone == nil {
		return nil
	}
	err := mw.ondone(ctx)
	if err != nil {
		log.Println("ondone:", err)
		mw.setState(Failed)
	}
	return err
}

func (mw *MonitoredWorker) GetState() State {
	mw.stateLock.RLock()
	defer mw.stateLock.RUnlock()
	return mw.state
}

func (mw *MonitoredWorker) GetId() string {
	if len(mw.guid) == 0 {
		mw.guid = genUid()
	}
	return mw.guid

}

func (mw *MonitoredWorker) Start(ctx context.Context) error {
	mw.lc.Lock()
	defer mw.lc.Unlock()
	switch mw.GetState() {
	case Running:
		return ErrRunRunningJob
	case Completed:
		if mw.GetProgress().IsCompleted() {
			return ErrRunCompletedJob
		}
	}
	mw.ctx, mw.cancelFunc = context.WithCancel(ctx)
	if err := mw.Itw.BeforeRun(mw.ctx); err != nil {
		mw.setState(Failed)
		mw.cancelFunc()
		return err
	}
	mw.setState(Running)
	mw.wgrun.Add(1)
	mw.id.Store(time.Now().Format(`20060102150405.000000`))
	go mw.wgoroute()

	return nil
}

func (mw *MonitoredWorker) Stop(ctx context.Context) error {
	mw.lc.Lock()
	defer mw.lc.Unlock()
	if mw.GetState() != Running {
		return ErrStopNonRunningJob
	}
	mw.cancelFunc()
	mw.id.Store(``)
	mw.wgrun.Wait()
	mw.setState(Stopped)
	log.Println("info: work stopped")
	if err := mw.Itw.AfterStop(); err != nil {
		return err
	}
	return nil
}

func (mw *MonitoredWorker) GetProgress() model.DownloadProgress {
	return mw.Itw.GetProgress()
}

func (mw *MonitoredWorker) ResetProgress() {
	mw.Itw.ResetProgress()
}
