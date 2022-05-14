package monitor

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"sync"
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
	DoWork() (bool, error)
	GetProgress() interface{}
	BeforeRun() error
	AfterStop() error
	IsPartialDownload() bool
}

func genUid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type MonitoredWorker struct {
	lc    sync.Mutex
	Itw   DiscretWork
	wgrun sync.WaitGroup
	guid  string
	state State
	chsig chan State
	//stwg   sync.WaitGroup
	ondone     func(context.Context) error
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func (mw *MonitoredWorker) wgoroute() {
	log.Println("info: work start", mw.GetId())
	defer func() {
		mw.cancelFunc()
		log.Println("info: release work guid", mw.GetId())
		mw.wgrun.Done()
		close(mw.chsig)
		mw.chsig = nil
	}()
	done := make(chan error)
	go func() {
		isdone, err := mw.Itw.DoWork()
		defer func() {
			done <- err
			close(done)
		}()
		if err != nil {
			log.Println("error: guid", mw.guid, " work failed", err)
			mw.state = Failed
			return
		}
		if isdone {
			if mw.ondone != nil {
				err = mw.ondone(mw.ctx)
				if err != nil {
					log.Println("ondone:", err)
					mw.state = Failed
					return
				}
			}
			mw.state = Completed
			log.Println("info: work done")
		}
	}()

	for {
		select {
		case newState := <-mw.chsig:
			if newState == Stopped {
				mw.state = newState
				log.Println("info: work stopped")
				return
			}
		case <-done:
			return
		}
	}
}

func (mw *MonitoredWorker) GetState() State {
	return mw.state
}

func (mw *MonitoredWorker) GetId() string {
	if len(mw.guid) == 0 {
		mw.guid = genUid()
	}
	return mw.guid

}

func (mw *MonitoredWorker) Start() error {
	mw.lc.Lock()
	defer mw.lc.Unlock()
	if mw.state == Running {
		return ErrRunRunningJob
	}
	if mw.state == Completed {
		return ErrRunCompletedJob
	}
	if err := mw.Itw.BeforeRun(); err != nil {
		mw.state = Failed
		return err
	}
	mw.ctx, mw.cancelFunc = context.WithCancel(context.Background())
	mw.chsig = make(chan State)
	mw.state = Running
	mw.wgrun.Add(1)
	go mw.wgoroute()

	return nil
}

func (mw *MonitoredWorker) Stop() error {
	mw.lc.Lock()
	defer mw.lc.Unlock()
	if mw.state != Running {
		return ErrStopNonRunningJob
	}
	if mw.chsig != nil {
		mw.chsig <- Stopped
	}
	mw.wgrun.Wait()
	if err := mw.Itw.AfterStop(); err != nil {
		return err
	}
	return nil
}

func (mw *MonitoredWorker) Wait() {
	mw.wgrun.Wait()
}

func (mw *MonitoredWorker) GetProgress() interface{} {
	return mw.Itw.GetProgress()
}
