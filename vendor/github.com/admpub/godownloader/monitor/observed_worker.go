package monitor

import (
	"crypto/rand"
	"errors"
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
}

func genUid() string {
	b := make([]byte, 16)
	rand.Read(b)
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

type MonitoredWorker struct {
	lc     sync.Mutex
	Itw    DiscretWork
	wgrun  sync.WaitGroup
	guid   string
	state  State
	chsig  chan State
	stwg   sync.WaitGroup
	ondone func()
}

func (mw *MonitoredWorker) wgoroute() {
	log.Println("info: work start", mw.GetId())
	defer func() {
		log.Print("info: release work guid ", mw.GetId())
		mw.wgrun.Done()
		if mw.ondone != nil && mw.state == Completed {
			mw.ondone()
		}
		close(mw.chsig)
		mw.chsig = nil
	}()

	for {
		select {
		case newState := <-mw.chsig:
			if newState == Stopped {
				mw.state = newState
				log.Println("info: work stopped")
				return
			}
		default:
			isdone, err := mw.Itw.DoWork()
			if err != nil {
				log.Println("error: guid", mw.guid, " work failed", err)
				mw.state = Failed
				return
			}
			if isdone {
				mw.state = Completed
				log.Println("info: work done")
				return
			}
		}
	}
}

func (mw MonitoredWorker) GetState() State {
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
	if mw.state == Completed {
		return errors.New("error: try run completed job")
	}
	if mw.state == Running {
		return errors.New("error: try run runing job")
	}
	if err := mw.Itw.BeforeRun(); err != nil {
		mw.state = Failed
		return err
	}
	mw.chsig = make(chan State, 1)
	mw.state = Running
	mw.wgrun.Add(1)
	go mw.wgoroute()

	return nil
}

func (mw *MonitoredWorker) Stop() error {
	mw.lc.Lock()
	defer mw.lc.Unlock()
	if mw.state != Running {
		return errors.New("error: imposible stop non runing job")
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

func (mw MonitoredWorker) GetProgress() interface{} {
	return mw.Itw.GetProgress()
}
