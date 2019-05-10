// Copyright 2015 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements logging with severity levels and message categories.
package log

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// RFC5424 log message levels.
const (
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
)

const (
	ActionNothing Action = iota
	ActionPanic
	ActionExit
)

type (
	// Level describes the level of a log message.
	Level  int
	Action int
)

var (
	// LevelNames maps log levels to names
	LevelNames = map[Level]string{
		LevelDebug: "Debug",
		LevelInfo:  "Info",
		LevelWarn:  "Warn",
		LevelError: "Error",
		LevelFatal: "Fatal",
	}

	Levels = map[string]Level{
		"Debug": LevelDebug,
		"Info":  LevelInfo,
		"Warn":  LevelWarn,
		"Error": LevelError,
		"Fatal": LevelFatal,
	}

	CallDepth = 5
)

func GetLevel(level string) (Level, bool) {
	level = strings.Title(level)
	l, y := Levels[level]
	return l, y
}

// String returns the string representation of the log level
func (l Level) String() string {
	if name, ok := LevelNames[l]; ok {
		return name
	}
	return "Unknown"
}

type LoggerWriter struct {
	Level Level
	*Logger
}

func (l *LoggerWriter) Write(p []byte) (n int, err error) {
	var s string
	n = len(p)
	if p[n-1] == '\n' {
		s = string(p[0 : n-1])
	} else {
		s = string(p)
	}
	l.Logger.Log(l.Level, s)
	return
}

func (l *LoggerWriter) Printf(format string, v ...interface{}) {
	l.Logger.Logf(l.Level, format, v...)
}

// Entry represents a log entry.
type Entry struct {
	Level     Level
	Category  string
	Message   string
	Time      time.Time
	CallStack string

	FormattedMessage string
}

// String returns the string representation of the log entry
func (e *Entry) String() string {
	return e.FormattedMessage
}

// Target represents a target where the logger can send log messages to for further processing.
type Target interface {
	// Open prepares the target for processing log messages.
	// Open will be invoked when Logger.Open() is called.
	// If an error is returned, the target will be removed from the logger.
	// errWriter should be used to write errors found while processing log messages.
	Open(errWriter io.Writer) error
	// Process processes an incoming log message.
	Process(*Entry)
	// Close closes a target.
	// Close is called when Logger.Close() is called, which gives each target
	// a chance to flush the logged messages to their destination storage.
	Close()
	SetLevel(interface{})
	SetLevels(...Level)
}

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
	CallStack     map[Level]*CallStack
	MaxLevel      Level    // the maximum level of messages to be logged
	Targets       []Target // targets for sending log messages to
	MaxGoroutines int32    // Max Goroutine
	AddSpace      bool     // Add a space between two arguments.
}

type CallStack struct {
	Depth  int
	Filter string
}

// Formatter formats a log message into an appropriate string.
type Formatter func(*Logger, *Entry) string

// Logger records log messages and dispatches them to various targets for further processing.
type Logger struct {
	*coreLogger
	Category   string    // the category associated with this logger
	Formatter  Formatter // message formatter
	categories map[string]*Logger
}

// NewLogger creates a root logger.
// The new logger takes these default options:
// ErrorWriter: os.Stderr, BufferSize: 1024, MaxLevel: LevelDebug,
// Category: app, Formatter: DefaultFormatter
func NewLogger(args ...string) *Logger {
	logger := &coreLogger{
		ErrorWriter: os.Stderr,
		BufferSize:  1024,
		MaxLevel:    LevelDebug,
		CallStack:   make(map[Level]*CallStack),
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
		logger = &Logger{
			coreLogger: l.coreLogger,
			Category:   category,
			categories: make(map[string]*Logger),
		}
		if len(formatter) > 0 {
			logger.Formatter = formatter[0]
		} else {
			logger.Formatter = l.Formatter
		}
		l.categories[category] = logger
	} else {
		if len(formatter) > 0 {
			logger.Formatter = formatter[0]
		}
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
func (l *Logger) Logf(level Level, format string, a ...interface{}) {
	if level > l.MaxLevel || !l.open {
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
func (l *Logger) Log(level Level, a ...interface{}) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if level > l.MaxLevel || !l.open {
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
func (l *Logger) newEntry(level Level, message string) {
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
		return errors.New("Logger.ErrorWriter must be set.")
	}
	var size int
	if !l.syncMode {
		if l.BufferSize < 0 {
			return errors.New("Logger.BufferSize must be no less than 0.")
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
	l.lock.RLock()
	defer l.lock.RUnlock()
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

// DefaultFormatter is the default formatter used to format every log message.
func DefaultFormatter(l *Logger, e *Entry) string {
	return e.Time.Format(time.RFC3339) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
}

func NormalFormatter(l *Logger, e *Entry) string {
	return e.Time.Format(`2006-01-02 15:04:05`) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
}

func ShortFileFormatter(l *Logger, e *Entry) string {
	_, file, line, ok := runtime.Caller(CallDepth)
	if !ok {
		return e.Time.Format(`2006-01-02 15:04:05`) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
	}

	return e.Time.Format(`2006-01-02 15:04:05`) + "|" + filepath.Base(file) + ":" + strconv.Itoa(line) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
}

type JSONL struct {
	Time      string          `bson:"time" json:"time"`
	Level     string          `bson:"level" json:"level"`
	Category  string          `bson:"category" json:"category"`
	Message   json.RawMessage `bson:"message" json:"message"`
	CallStack string          `bson:"callStack" json:"callStack"`
}

func JSONFormatter(l *Logger, e *Entry) string {
	jsonl := &JSONL{
		Time:      e.Time.Format(`2006-01-02 15:04:05`),
		Level:     e.Level.String(),
		Category:  e.Category,
		Message:   []byte(`"` + e.Message + `"`),
		CallStack: e.CallStack,
	}
	if len(e.Message) > 0 {
		switch e.Message[0] {
		case '{', '[', '"':
			jsonl.Message = []byte(e.Message)
		}
	}
	b, err := json.Marshal(jsonl)
	if err != nil {
		fmt.Println(err.Error())
	}
	return string(b)
}

// GetCallStack returns the current call stack information as a string.
// The skip parameter specifies how many top frames should be skipped, while
// the frames parameter specifies at most how many frames should be returned.
func GetCallStack(skip int, frames int, filter string) string {
	buf := new(bytes.Buffer)
	hasFilter := len(filter) > 0
	for i, count := skip, 0; count < frames; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if !hasFilter || strings.Contains(file, filter) {
			fmt.Fprintf(buf, "\n%s:%d", file, line)
			count++
		}
	}
	return buf.String()
}
