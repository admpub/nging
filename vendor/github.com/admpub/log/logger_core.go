package log

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// coreLogger maintains the log messages in a channel and sends them to various targets.
type coreLogger struct {
	lock        sync.RWMutex
	open        bool        // whether the logger is open
	entries     chan *Entry // log entries
	sendN       uint32
	procsN      uint32
	waiting     *sync.Once
	fatalAction Action
	syncMode    bool // Whether the use of non-asynchronous mode （是否使用非异步模式）

	ErrorWriter   io.Writer // the writer used to write errors caused by log targets
	BufferSize    int       // the size of the channel storing log entries
	CallStack     map[Leveler]*CallStack
	MaxLevel      Leveler  // the maximum level of messages to be logged
	Targets       []Target // targets for sending log messages to
	MaxGoroutines int32    // Max Goroutine
	AddSpace      bool     // Add a space between two arguments.
}

func (l *coreLogger) wait() {
	l.waiting.Do(func() {
		for {
			sendN := atomic.LoadUint32(&l.sendN)
			//fmt.Println(`waiting ...`, len(l.entries), sendN)
			if sendN <= atomic.LoadUint32(&l.procsN) {
				l.sendN = 0
				l.procsN = 0
				l.waiting = &sync.Once{}
				return
			}
			delay := sendN
			if delay < 500 {
				delay = 500
			}
			time.Sleep(time.Duration(delay) * time.Microsecond)
		}
	})
}

// Open prepares the logger and the targets for logging purpose.
// Open must be called before any message can be logged.
func (l *coreLogger) Open() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.open {
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

	l.open = true

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
		atomic.AddUint32(&l.procsN, 1)
	}
}

// Close closes the logger and the targets.
// Existing messages will be processed before the targets are closed.
// New incoming messages will be discarded after calling this method.
func (l *coreLogger) Close() {
	l.lock.Lock()
	defer l.lock.Unlock()

	if !l.open {
		return
	}
	l.open = false
	l.wait()
	// use a nil entry to signal the close of logger
	l.entries <- nil
	for _, target := range l.Targets {
		target.Close()
	}
}
