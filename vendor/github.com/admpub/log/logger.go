// Copyright 2015 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements logging with severity levels and message categories.
package log

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

const (
	ActionNothing Action = iota
	ActionPanic
	ActionExit
)

// Logger records log messages and dispatches them to various targets for further processing.
type Logger struct {
	*coreLogger
	Category   string    // the category associated with this logger
	Formatter  Formatter // message formatter
	CallDepth  int
	categories map[string]*Logger
}

// NewLogger creates a root logger.
// The new logger takes these default options:
// ErrorWriter: os.Stderr, BufferSize: 1024, MaxLevel: LevelDebug,
// Category: app, Formatter: DefaultFormatter
func NewLogger(args ...string) *Logger {
	return NewWithCallDepth(DefaultCallDepth, args...)
}

// NewWithCallDepth creates a root logger.
func NewWithCallDepth(callDepth int, args ...string) *Logger {
	logger := &coreLogger{
		ErrorWriter: os.Stderr,
		BufferSize:  1024,
		MaxLevel:    LevelDebug,
		CallStack:   make(map[Leveler]*CallStack),
		Targets:     make([]Target, 0),
		waiting:     &sync.Once{},
	}
	category := `app`
	if len(args) > 0 {
		category = args[0]
	}
	logger.Targets = append(logger.Targets, NewConsoleTarget())
	logger.Open()
	return &Logger{
		coreLogger: logger,
		Category:   category,
		Formatter:  NormalFormatter,
		CallDepth:  callDepth,
		categories: make(map[string]*Logger),
	}
}

func New(args ...string) *Logger {
	return NewLogger(args...)
}

// GetLogger creates a logger with the specified category and log formatter.
// Messages logged through this logger will carry the same category name.
// The formatter, if not specified, will inherit from the calling logger.
// It will be used to format all messages logged through this logger.
func (l *Logger) GetLogger(category string, formatter ...Formatter) *Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	logger, ok := l.categories[category]
	if !ok {
		logger = l.clone()
		logger.Category = category
		l.categories[category] = logger
	}
	if len(formatter) > 0 {
		logger.Formatter = formatter[0]
	}
	return logger
}

func (l *Logger) clone() *Logger {
	logger := &Logger{
		coreLogger: l.coreLogger,
		Category:   l.Category,
		categories: make(map[string]*Logger),
		CallDepth:  l.CallDepth,
		Formatter:  l.Formatter,
	}
	return logger
}

func (l *Logger) Sync(args ...bool) *Logger {
	var on bool
	if len(args) > 0 {
		on = !args[0]
	}
	return l.Async(on)
}

func (l *Logger) sendEntry(entry *Entry) {
	atomic.AddUint32(&l.sendN, 1)
	l.entries <- entry
}

func (l *Logger) Async(args ...bool) *Logger {
	if len(args) < 1 {
		l.syncMode = false
		return l
	}
	l.syncMode = !args[0]
	if l.open {
		l.Close()
		l.Open()
	}
	return l
}

func (l *Logger) SetTarget(targets ...Target) *Logger {
	l.Close()
	if len(targets) > 0 {
		l.Targets = targets
		l.Open()
	} else {
		l.Targets = []Target{}
	}
	return l
}

func (l *Logger) SetFatalAction(action Action) *Logger {
	l.fatalAction = action
	return l
}

func (l *Logger) AddTarget(targets ...Target) *Logger {
	l.Close()
	l.Targets = append(l.Targets, targets...)
	l.Open()
	return l
}

func (l *Logger) SetLevel(level string) *Logger {
	if le, ok := GetLevel(level); ok {
		l.MaxLevel = le
	}
	return l
}

func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.Logf(LevelFatal, format, a...)
}

// Errorf logs a message indicating an error condition.
// This method takes one or multiple parameters. If a single parameter
// is provided, it will be treated as the log message. If multiple parameters
// are provided, they will be passed to fmt.Sprintf() to generate the log message.
func (l *Logger) Errorf(format string, a ...interface{}) {
	l.Logf(LevelError, format, a...)
}

// Warnf logs a message indicating a warning condition.
// Please refer to Error() for how to use this method.
func (l *Logger) Warnf(format string, a ...interface{}) {
	l.Logf(LevelWarn, format, a...)
}

// Infof logs a message for informational purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Infof(format string, a ...interface{}) {
	l.Logf(LevelInfo, format, a...)
}

// Debugf logs a message for debugging purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Logf(LevelDebug, format, a...)
}

// Logf logs a message of a specified severity level.
func (l *Logger) Logf(level Leveler, format string, a ...interface{}) {
	if level.Int() > l.MaxLevel.Int() || !l.open {
		return
	}
	message := format
	if len(a) > 0 {
		message = fmt.Sprintf(format, a...)
	}
	l.newEntry(level, message)
}

func (l *Logger) Writer(level Level) io.Writer {
	return &LoggerWriter{
		Level:  level,
		Logger: l,
	}
}

func (l *Logger) Fatal(a ...interface{}) {
	l.Log(LevelFatal, a...)
}

// Error logs a message indicating an error condition.
// This method takes one or multiple parameters. If a single parameter
// is provided, it will be treated as the log message. If multiple parameters
// are provided, they will be passed to fmt.Sprintf() to generate the log message.
func (l *Logger) Error(a ...interface{}) {
	l.Log(LevelError, a...)
}

// Warn logs a message indicating a warning condition.
// Please refer to Error() for how to use this method.
func (l *Logger) Warn(a ...interface{}) {
	l.Log(LevelWarn, a...)
}

// Info logs a message for informational purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Info(a ...interface{}) {
	l.Log(LevelInfo, a...)
}

// Debug logs a message for debugging purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Debug(a ...interface{}) {
	l.Log(LevelDebug, a...)
}

// Log logs a message of a specified severity level.
func (l *Logger) Log(level Leveler, a ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if level.Int() > l.MaxLevel.Int() || !l.open {
		return
	}
	var message string
	if l.AddSpace {
		message = fmt.Sprintln(a...)
		message = message[:len(message)-1]
	} else {
		message = fmt.Sprint(a...)
	}
	l.newEntry(level, message)
}

// Log logs a message of a specified severity level.
func (l *Logger) newEntry(level Leveler, message string) {
	entry := &Entry{
		Category: l.Category,
		Level:    level,
		Message:  message,
		Time:     time.Now(),
	}
	if level == LevelFatal {
		var (
			stackDepth  int
			stackFilter string
		)
		if cs, ok := l.CallStack[level]; ok && cs != nil {
			stackDepth = cs.Depth
			stackFilter = cs.Filter
		}
		if stackDepth < 20 {
			stackDepth = 20
		}
		entry.CallStack = GetCallStack(3, stackDepth, stackFilter)
		entry.FormattedMessage = l.Formatter(l, entry)
		l.sendEntry(entry)
		l.wait()
		switch l.fatalAction {
		case ActionPanic:
			panic(entry.FormattedMessage)
		case ActionExit:
			os.Exit(-1)
		}
		return
	}
	if cs, ok := l.CallStack[level]; ok && cs != nil && cs.Depth > 0 {
		entry.CallStack = GetCallStack(3, cs.Depth, cs.Filter)
	}
	entry.FormattedMessage = l.Formatter(l, entry)
	l.sendEntry(entry)
}
