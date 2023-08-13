package log

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
)

// coreLogger maintains the log messages in a channel and sends them to various targets.
type coreLogger struct {
	lock        sync.RWMutex
	open        *atomic.Bool // whether the logger is open
	entries     chan *Entry  // log entries
	fatalAction Action
	syncMode    bool

	ErrorWriter io.Writer // the writer used to write errors caused by log targets
	BufferSize  int       // the size of the channel storing log entries
	CallStack   map[Leveler]*CallStack
	MaxLevel    Leveler  // the maximum level of messages to be logged
	Targets     []Target // targets for sending log messages to
	AddSpace    bool     // Add a space between two arguments.
	pid         int
}

// Open prepares the logger and the targets for logging purpose.
// Open must be called before any message can be logged.
func (l *coreLogger) Open() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.open.Load() {
		return nil
	}
	if l.ErrorWriter == nil {
		return errors.New("Logger.ErrorWriter must be set")
	}
	var size int
	if !l.syncMode {
		if l.BufferSize < 0 {
			return errors.New("Logger.BufferSize must be no less than 0")
		}
		size = l.BufferSize
	}
	l.entries = make(chan *Entry, size)
	var targets []Target
	for _, target := range l.Targets {
		if err := target.Open(l.ErrorWriter); err != nil {
			fmt.Fprintf(l.ErrorWriter, "Failed to open target: %v\n", err)
		} else {
			targets = append(targets, target)
		}
	}
	l.Targets = targets

	go l.process()

	l.open.Store(true)

	return nil
}

// process sends the messages to targets for processing.
func (l *coreLogger) process() {
	for {
		entry := <-l.entries
		for _, target := range l.Targets {
			target.Process(entry)
		}
		if entry == nil {
			break
		}
	}
}

// Close closes the logger and the targets.
// Existing messages will be processed before the targets are closed.
// New incoming messages will be discarded after calling this method.
func (l *coreLogger) Close() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if !l.open.Load() {
		return
	}
	l.open.Store(false)
	// use a nil entry to signal the close of logger
	l.entries <- nil
	for _, target := range l.Targets {
		target.Close()
	}
}

func (l *coreLogger) setCallStack(level Level, callStack *CallStack) {
	l.lock.Lock()
	l.CallStack[level] = callStack
	l.lock.Unlock()
}

func (l *coreLogger) Pid() int {
	return l.pid
}
