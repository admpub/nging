//go:build linux

package reactor

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/godzie44/go-uring/uring"
)

type Callback func(event uring.CQEvent)

type configuration struct {
	tickDuration time.Duration
	logger       Logger
}

type Option func(cfg *configuration)

//WithTickTimeout set tick duration for event loop.
func WithTickTimeout(duration time.Duration) Option {
	return func(cfg *configuration) {
		cfg.tickDuration = duration
	}
}

//WithLogger set logger for event loop.
func WithLogger(l Logger) Option {
	return func(cfg *configuration) {
		cfg.logger = l
	}
}

//Reactor is event loop's manager with main responsibility - handling client requests and return responses asynchronously.
type Reactor struct {
	loops []*ringEventLoop

	currentNonce uint64

	config *configuration
}

//New create new reactor instance.
//rings - io_uring instances. The reactor will create one event loop for each instance.
//opts - reactor options.
func New(rings []*uring.Ring, opts ...Option) (*Reactor, error) {
	for _, ring := range rings {
		if err := checkRingReq(ring, false); err != nil {
			return nil, err
		}
	}

	r := &Reactor{
		config: &configuration{
			tickDuration: time.Millisecond * 1,
			logger:       &nopLogger{},
		},
	}

	for _, opt := range opts {
		opt(r.config)
	}

	for _, ring := range rings {
		loop := newRingEventLoop(ring, r.config.logger)
		r.loops = append(r.loops, loop)
	}

	return r, nil
}

//Run start reactor.
func (r *Reactor) Run(ctx context.Context) {
	for _, loop := range r.loops {
		go loop.runConsumer(r.config.tickDuration)
		go loop.runPublisher()
	}

	<-ctx.Done()

	for _, loop := range r.loops {
		loop.stopConsumer()
		loop.stopPublisher()
	}
}

func (r *Reactor) queue(op uring.Operation, cb Callback, timeout time.Duration) (uint64, error) {
	nonce := r.nextNonce()

	loop := r.loopForNonce(nonce)
	err := loop.Queue(subSqeRequest{op, 0, nonce, timeout}, cb)

	return nonce, err
}

func (r *Reactor) loopForNonce(nonce uint64) *ringEventLoop {
	n := len(r.loops)
	return r.loops[nonce%uint64(n)]
}

//Queue io_uring operation. Callback function `cb` calling when receive cqe.
//Return uint64 which can be used as the SQE identifier.
func (r *Reactor) Queue(op uring.Operation, cb Callback) (uint64, error) {
	return r.queue(op, cb, time.Duration(0))
}

//QueueWithDeadline io_uring operation. Callback function `cb` calling when receive cqe.
//After a deadline time, a CQE with the error ECANCELED will be placed in the channel retChan.
func (r *Reactor) QueueWithDeadline(op uring.Operation, cb Callback, deadline time.Time) (uint64, error) {
	if deadline.IsZero() {
		return r.Queue(op, cb)
	}

	return r.queue(op, cb, time.Until(deadline))
}

//Cancel queued operation.
//nonce - SQE id returned by Queue method.
func (r *Reactor) Cancel(nonce uint64) error {
	loop := r.loopForNonce(nonce)
	return loop.cancel(nonce)
}

type ringEventLoop struct {
	ring *uring.Ring

	callbacks     map[uint64]Callback
	callbacksLock sync.Mutex

	queueSQELock sync.Mutex

	submitSignal chan struct{}

	logger Logger

	stopConsumerCh  chan struct{}
	stopPublisherCh chan struct{}

	needSubmit uint32
}

func newRingEventLoop(ring *uring.Ring, logger Logger) *ringEventLoop {
	return &ringEventLoop{
		ring:            ring,
		submitSignal:    make(chan struct{}),
		stopConsumerCh:  make(chan struct{}),
		stopPublisherCh: make(chan struct{}),
		callbacks:       map[uint64]Callback{},
		logger:          logger,
	}
}

func (loop *ringEventLoop) runConsumer(tickDuration time.Duration) {
	runtime.LockOSThread()

	cqeBuff := make([]*uring.CQEvent, cqeBuffSize)
	for {
		loop.submitSignal <- struct{}{}

		_, err := loop.ring.WaitCQEventsWithTimeout(1, tickDuration)

		if errors.Is(err, syscall.EAGAIN) || errors.Is(err, syscall.EINTR) || errors.Is(err, syscall.ETIME) {
			runtime.Gosched()
			goto CheckCtxAndContinue
		}

		if err != nil {
			loop.logger.Log("io_uring wait", err)
			goto CheckCtxAndContinue
		}

		for n := loop.ring.PeekCQEventBatch(cqeBuff); n > 0; n = loop.ring.PeekCQEventBatch(cqeBuff) {
			for i := 0; i < n; i++ {
				cqe := cqeBuff[i]

				nonce := cqe.UserData
				if nonce == timeoutNonce || nonce == cancelNonce {
					continue
				}

				loop.callbacksLock.Lock()
				loop.callbacks[nonce](uring.CQEvent{
					UserData: cqe.UserData,
					Res:      cqe.Res,
					Flags:    cqe.Flags,
				})
				delete(loop.callbacks, nonce)
				loop.callbacksLock.Unlock()
			}

			loop.ring.AdvanceCQ(uint32(n))
		}

	CheckCtxAndContinue:
		select {
		case <-loop.stopConsumerCh:
			close(loop.stopConsumerCh)
			return
		default:
			continue
		}
	}
}

func (loop *ringEventLoop) stopConsumer() {
	loop.stopConsumerCh <- struct{}{}
	<-loop.stopConsumerCh
}

func (loop *ringEventLoop) stopPublisher() {
	loop.stopPublisherCh <- struct{}{}
	<-loop.stopPublisherCh
}

func (loop *ringEventLoop) cancel(nonce uint64) error {
	op := uring.Cancel(nonce, 0)

	return loop.Queue(subSqeRequest{
		op:       op,
		userData: cancelNonce,
	}, nil)
}

func (loop *ringEventLoop) Queue(req subSqeRequest, cb Callback) (err error) {
	loop.queueSQELock.Lock()
	defer loop.queueSQELock.Unlock()

	atomic.StoreUint32(&loop.needSubmit, 1)

	if req.timeout == 0 {
		err = loop.ring.QueueSQE(req.op, req.flags, req.userData)
	} else {
		err = loop.ring.QueueSQE(req.op, req.flags|uring.SqeIOLinkFlag, req.userData)
		if err == nil {
			_ = loop.ring.QueueSQE(uring.LinkTimeout(req.timeout), 0, timeoutNonce)
		}
	}

	if err == nil {
		loop.callbacksLock.Lock()
		loop.callbacks[req.userData] = cb
		loop.callbacksLock.Unlock()
	}

	return err
}

func (loop *ringEventLoop) runPublisher() {
	runtime.LockOSThread()

	defer close(loop.submitSignal)

	var err error
	for {
		select {
		case <-loop.submitSignal:
			if atomic.CompareAndSwapUint32(&loop.needSubmit, 1, 0) {
				loop.queueSQELock.Lock()
				_, err = loop.ring.Submit()
				loop.queueSQELock.Unlock()

				if err != nil {
					loop.logger.Log("io_uring submit", err)
				}
			}
		case <-loop.stopPublisherCh:
			close(loop.stopPublisherCh)
			return
		}
	}
}

func (r *Reactor) nextNonce() uint64 {
	local := atomic.AddUint64(&r.currentNonce, 1)

	for local >= cancelNonce {
		local = atomic.AddUint64(&r.currentNonce, 1)
	}

	return local
}
