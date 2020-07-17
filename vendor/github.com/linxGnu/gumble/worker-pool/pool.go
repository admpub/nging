package workerpool

import (
	"context"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

var numCPU int

func init() {
	numCPU = runtime.NumCPU()
}

// TaskResult represent result of task.
type TaskResult struct {
	Result interface{}
	Err    error
}

// Task represent a task.
type Task struct {
	ctx      context.Context
	executor func(context.Context) (interface{}, error)
	future   chan *TaskResult
}

// NewTask create new task.
func NewTask(ctx context.Context, executor func(context.Context) (interface{}, error)) *Task {
	return &Task{
		ctx:      ctx,
		executor: executor,
		future:   make(chan *TaskResult, 1),
	}
}

// Execute task.
func (t *Task) Execute() {
	var result interface{}
	var err error

	if t.executor != nil {
		result, err = t.executor(t.ctx)
	}

	t.future <- &TaskResult{Result: result, Err: err}
}

// Result pushed via channel
func (t *Task) Result() <-chan *TaskResult {
	return t.future
}

// Option represents pool option.
type Option struct {
	// NumberWorker number of workers.
	// Default: runtime.NumCPU()
	NumberWorker int `yaml:"number_worker" json:"number_worker"`
	// ExpandableLimit limits number of workers to be expanded on demand.
	// Default: 0 (no expandable)
	ExpandableLimit int32 `yaml:"expandable_limit" json:"expandable_limit"`
	// ExpandedLifetime represents lifetime of expanded worker (in nanoseconds)/
	// Default: 1 minute.
	ExpandedLifetime time.Duration `yaml:"expanded_lifetime" json:"expanded_lifetime"`
}

func (o *Option) normalize() {
	if o.NumberWorker <= 0 {
		o.NumberWorker = numCPU
	}

	if o.ExpandableLimit < 0 {
		o.ExpandableLimit = 0
	}

	if o.ExpandedLifetime <= 0 {
		o.ExpandedLifetime = time.Minute
	}
}

// Pool is a lightweight worker pool with capable of auto-expand on demand.
type Pool struct {
	ctx    context.Context
	cancel context.CancelFunc

	opt Option

	wg        sync.WaitGroup
	taskQueue chan *Task
	expanded  int32

	state uint32 // 0: not start, 1: started, 2: stopped
}

// NewPool create new worker pool
func NewPool(ctx context.Context, opt Option) (p *Pool) {
	if ctx == nil {
		ctx = context.Background()
	}

	// normalize option
	opt.normalize()

	// set up pool
	p = &Pool{
		opt:       opt,
		taskQueue: make(chan *Task, opt.NumberWorker),
	}
	p.ctx, p.cancel = context.WithCancel(ctx)
	return
}

// Start workers.
func (p *Pool) Start() {
	if atomic.CompareAndSwapUint32(&p.state, 0, 1) {
		numWorker := p.opt.NumberWorker

		p.wg.Add(numWorker)
		for i := 0; i < numWorker; i++ {
			go p.worker()
		}
	}
}

// Stop worker. Wait all task done.
func (p *Pool) Stop() {
	if atomic.CompareAndSwapUint32(&p.state, 1, 2) || atomic.CompareAndSwapUint32(&p.state, 0, 2) {
		// cancel context
		p.cancel()

		// wait child workers
		close(p.taskQueue)
		p.wg.Wait()
	}
}

// Execute a task.
func (p *Pool) Execute(exec func(context.Context) (interface{}, error)) (t *Task) {
	return p.ExecuteWithCtx(p.ctx, exec)
}

// ExecuteWithCtx a task with custom context.
func (p *Pool) ExecuteWithCtx(ctx context.Context, exec func(context.Context) (interface{}, error)) (t *Task) {
	if ctx == nil {
		ctx = p.ctx
	}
	t = NewTask(ctx, exec)
	p.Do(t)
	return
}

// TryExecute try to execute a task. If task queue is full, returns immediately and
// addedToQueue is false.
func (p *Pool) TryExecute(exec func(context.Context) (interface{}, error)) (t *Task, addedToQueue bool) {
	return p.TryExecuteWithCtx(p.ctx, exec)
}

// TryExecuteWithCtx try to execute a task with custom context. If task queue is full, returns immediately and
// addedToQueue is false.
func (p *Pool) TryExecuteWithCtx(ctx context.Context, exec func(context.Context) (interface{}, error)) (t *Task, addedToQueue bool) {
	if ctx == nil {
		ctx = p.ctx
	}
	t = NewTask(ctx, exec)
	addedToQueue = p.TryDo(t)
	return
}

// Do a task.
func (p *Pool) Do(t *Task) {
	if t != nil {
		if p.opt.ExpandableLimit == 0 {
			p.push(t)
		} else {
			select {
			case <-p.ctx.Done():
			case p.taskQueue <- t:
			default:
				if atomic.AddInt32(&p.expanded, 1) <= p.opt.ExpandableLimit {
					p.wg.Add(1)
					go p.expandedWorker()
				} else {
					atomic.AddInt32(&p.expanded, -1)
				}

				// push again
				p.push(t)
			}
		}
	}
}

func (p *Pool) push(t *Task) {
	select {
	case <-p.ctx.Done():
	case p.taskQueue <- t:
	}
}

// TryDo try to execute a task. If task queue is full, returns immediately and
// addedToQueue is false.
func (p *Pool) TryDo(t *Task) (addedToQueue bool) {
	if t != nil {
		select {
		case p.taskQueue <- t:
			addedToQueue = true

		case <-p.ctx.Done():

		default:
		}
	}
	return
}

func (p *Pool) worker() {
	for task := range p.taskQueue {
		task.Execute()
	}
	p.wg.Done()
}

func (p *Pool) expandedWorker() {
	lifetime := p.opt.ExpandedLifetime
	timer := time.NewTimer(lifetime)
	for {
		select {
		case task, ok := <-p.taskQueue:
			if !ok {
				p.stopExpendedWorker(timer)
				return
			}

			task.Execute()

			// received task, expand the lifetime
			resetTimer(timer, lifetime)

		case <-timer.C:
			p.stopExpendedWorker(timer)
			return

		}
	}
}

func (p *Pool) stopExpendedWorker(t *time.Timer) {
	stopTimer(t)
	p.wg.Done()
	atomic.AddInt32(&p.expanded, -1)
}

func stopTimer(t *time.Timer) {
	if !t.Stop() {
		<-t.C
	}
}

func resetTimer(t *time.Timer, d time.Duration) {
	stopTimer(t)
	t.Reset(d)
}
