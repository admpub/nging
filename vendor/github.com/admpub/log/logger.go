// Copyright 2015 Qiang Xue. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package log implements logging with severity levels and message categories.
package log

import (
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

const (
	// ActionNothing 无操作
	ActionNothing Action = iota
	// ActionPanic 触发Panic
	ActionPanic
	// ActionExit 退出程序
	ActionExit
)

// Logger records log messages and dispatches them to various targets for further processing.
type Logger struct {
	*coreLogger
	Category   string    // the category associated with this logger
	Formatter  Formatter // message formatter
	Emoji      bool
	catelock   sync.RWMutex
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
		CallStack:   make(map[Leveler]*CallStack),
		Targets:     make([]Target, 0),
		fatalAction: ActionPanic,
		pid:         os.Getpid(),
		open:        &atomic.Bool{},
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

// New creates a new Logger
func New(args ...string) *Logger {
	return NewLogger(args...)
}

// GetLogger creates a logger with the specified category and log formatter.
// Messages logged through this logger will carry the same category name.
// The formatter, if not specified, will inherit from the calling logger.
// It will be used to format all messages logged through this logger.
func (l *Logger) GetLogger(category string, formatter ...Formatter) *Logger {
	l.catelock.RLock()
	logger, ok := l.categories[category]
	l.catelock.RUnlock()
	if !ok {
		logger = l.clone()
		logger.Category = category
		l.catelock.Lock()
		l.categories[category] = logger
		l.catelock.Unlock()
	}
	if len(formatter) > 0 {
		logger.Formatter = formatter[0]
	}
	return logger
}

func (l *Logger) Categories() []string {
	categories := make([]string, len(l.categories))
	var i int
	for k := range l.categories {
		categories[i] = k
		i++
	}
	sort.Strings(categories)
	return categories
}

func (l *Logger) HasCategory(category string) bool {
	l.catelock.RLock()
	_, ok := l.categories[category]
	l.catelock.RUnlock()
	return ok
}

func (l *Logger) SetEmoji(on bool) *Logger {
	l.Emoji = on
	return l
}

func (l *Logger) EmojiOfLevel(level Level) string {
	if l.Emoji {
		return GetLevelEmoji(level)
	}
	return ``
}

func (l *Logger) clone() *Logger {
	logger := &Logger{
		coreLogger: l.coreLogger,
		Category:   l.Category,
		categories: make(map[string]*Logger),
		Formatter:  l.Formatter,
		Emoji:      l.Emoji,
	}
	return logger
}

// Sync 同步日志
func (l *Logger) Sync(args ...bool) *Logger {
	var on bool
	if len(args) > 0 {
		on = !args[0]
	}
	return l.Async(on)
}

func (l *Logger) sendEntry(entry *Entry) {
	l.entries <- entry
}

// Async 异步日志
func (l *Logger) Async(args ...bool) *Logger {
	var syncMode bool
	if len(args) < 1 {
		syncMode = false
	} else {
		syncMode = !args[0]
	}
	if l.syncMode == syncMode {
		return l
	}
	l.syncMode = syncMode
	if l.open.Load() {
		l.Close()
		l.Open()
	}
	return l
}

// SetTarget 设置日志输出Target
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

// SetFatalAction 设置Fatal类型日志的行为
func (l *Logger) SetFatalAction(action Action) *Logger {
	l.fatalAction = action
	return l
}

// AddTarget 添加日志输出Target
func (l *Logger) AddTarget(targets ...Target) *Logger {
	l.Close()
	l.Targets = append(l.Targets, targets...)
	l.Open()
	return l
}

// SetLevel 添加日志输出等级
func (l *Logger) SetLevel(level string) *Logger {
	if lv, ok := GetLevel(level); ok {
		l.MaxLevel = lv
	}
	return l
}

// SetCallStack 设置各个等级的call stack配置
func (l *Logger) SetCallStack(level Level, callStack *CallStack) *Logger {
	l.coreLogger.setCallStack(level, callStack)
	return l
}

// SetFormatter 设置日志格式化处理函数
func (l *Logger) SetFormatter(formatter Formatter) *Logger {
	l.Formatter = formatter
	return l
}

// IsEnabled 是否启用了某个等级的日志输出
func (l *Logger) IsEnabled(level Level) bool {
	return l.MaxLevel.IsEnabled(level)
}

// Fatalf fatal
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

// Okayf logs a message indicating an okay condition.
func (l *Logger) Okayf(format string, a ...interface{}) {
	l.Logf(LevelOkay, format, a...)
}

// Infof logs a message for informational purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Infof(format string, a ...interface{}) {
	l.Logf(LevelInfo, format, a...)
}

// Progressf logs a message for how things are progressing.
func (l *Logger) Progressf(format string, a ...interface{}) {
	l.Logf(LevelProgress, format, a...)
}

// Debugf logs a message for debugging purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Debugf(format string, a ...interface{}) {
	l.Logf(LevelDebug, format, a...)
}

// Logf logs a message of a specified severity level.
func (l *Logger) Logf(level Leveler, format string, a ...interface{}) {
	if level.Int() > l.MaxLevel.Int() || !l.open.Load() {
		return
	}
	message := format
	if len(a) > 0 {
		message = fmt.Sprintf(format, a...)
	}
	l.newEntry(level, message)
}

// Writer 日志输出Writer
func (l *Logger) Writer(level Level) io.Writer {
	return &LoggerWriter{
		Level:  level,
		Logger: l,
	}
}

// Fatal fatal
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

// Okay logs a message indicating an okay condition.
func (l *Logger) Okay(a ...interface{}) {
	l.Log(LevelOkay, a...)
}

// Info logs a message for informational purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Info(a ...interface{}) {
	l.Log(LevelInfo, a...)
}

// Progress logs a message for how things are progressing.
func (l *Logger) Progress(a ...interface{}) {
	l.Log(LevelProgress, a...)
}

// Debug logs a message for debugging purpose.
// Please refer to Error() for how to use this method.
func (l *Logger) Debug(a ...interface{}) {
	l.Log(LevelDebug, a...)
}

// Log logs a message of a specified severity level.
func (l *Logger) Log(level Leveler, a ...interface{}) {
	l.lock.RLock()
	if level.Int() > l.MaxLevel.Int() || !l.open.Load() {
		l.lock.RUnlock()
		return
	}
	var message string
	if l.AddSpace {
		message = fmt.Sprintln(a...)
		message = message[:len(message)-1]
	} else {
		message = fmt.Sprint(a...)
	}
	l.lock.RUnlock()
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
		l.generateCallStack(entry, level, true)
		entry.FormattedMessage = l.Formatter(l, entry)
		l.sendEntry(entry)
		switch l.fatalAction {
		case ActionPanic:
			panic(entry.FormattedMessage)
		case ActionExit:
			l.Close()
			os.Exit(-1)
		}
		return
	}
	l.generateCallStack(entry, level, false)
	entry.FormattedMessage = l.Formatter(l, entry)
	l.sendEntry(entry)
}

func (l *Logger) generateCallStack(entry *Entry, level Leveler, force bool) *Logger {
	var (
		stackDepth   int
		skipStack    int
		stackFilters = []string{DefaultStackFilter}
	)
	cs, ok := l.CallStack[level]
	if ok && cs != nil {
		if !force && cs.Depth < 1 {
			return l
		}
		stackDepth = cs.Depth
		skipStack = cs.Skip
		if len(cs.Filters) > 0 {
			stackFilters = append(stackFilters, cs.Filters...)
		}
	} else {
		if !force {
			return l
		}
	}
	if stackDepth < 1 {
		skipStack = DefaultSkipStack
		stackDepth = DefaultStackDepth
	} else if skipStack < 0 {
		skipStack = 0
	}
	entry.CallStack = GetCallStack(skipStack, stackDepth, stackFilters...)
	return l
}
