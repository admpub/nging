package log

import (
	"io"
	"strings"

	"github.com/admpub/color"
)

// DefaultLog 默认日志实例
var DefaultLog = &defaultLogger{Logger: New()}

type defaultLogger struct {
	*Logger
}

func GetLogger(category string, formatter ...Formatter) *Logger {
	return DefaultLog.GetLogger(category, formatter...)
}

func Categories() []string {
	return DefaultLog.Categories()
}

func HasCategory(category string) bool {
	return DefaultLog.HasCategory(category)
}

func SetEmoji(on bool) *Logger {
	return DefaultLog.SetEmoji(on)
}

func EmojiOfLevel(level Level) string {
	return DefaultLog.EmojiOfLevel(level)
}

func Sync(args ...bool) *Logger {
	return DefaultLog.Sync(args...)
}

func Async(args ...bool) *Logger {
	return DefaultLog.Async(args...)
}

func SetTarget(targets ...Target) *Logger {
	return DefaultLog.SetTarget(targets...)
}

func SetFormatter(formatter Formatter) *Logger {
	return DefaultLog.SetFormatter(formatter)
}

func SetFatalAction(action Action) *Logger {
	return DefaultLog.SetFatalAction(action)
}

func AddTarget(targets ...Target) *Logger {
	return DefaultLog.AddTarget(targets...)
}

func SetLevel(level string) *Logger {
	return DefaultLog.SetLevel(level)
}

func SetCallStack(level Level, callStack *CallStack) *Logger {
	return DefaultLog.SetCallStack(level, callStack)
}

// IsEnabled 是否启用了某个等级的日志输出
func IsEnabled(level Level) bool {
	return DefaultLog.IsEnabled(level)
}

func Fatalf(format string, a ...interface{}) {
	DefaultLog.Fatalf(format, a...)
}

func Errorf(format string, a ...interface{}) {
	DefaultLog.Errorf(format, a...)
}

func Warnf(format string, a ...interface{}) {
	DefaultLog.Warnf(format, a...)
}

func Okayf(format string, a ...interface{}) {
	DefaultLog.Okayf(format, a...)
}

func Infof(format string, a ...interface{}) {
	DefaultLog.Infof(format, a...)
}

func Progressf(format string, a ...interface{}) {
	DefaultLog.Progressf(format, a...)
}

func Debugf(format string, a ...interface{}) {
	DefaultLog.Debugf(format, a...)
}

func Fatal(a ...interface{}) {
	DefaultLog.Fatal(a...)
}

func Error(a ...interface{}) {
	DefaultLog.Error(a...)
}

func Warn(a ...interface{}) {
	DefaultLog.Warn(a...)
}

func Okay(a ...interface{}) {
	DefaultLog.Okay(a...)
}

func Info(a ...interface{}) {
	DefaultLog.Info(a...)
}

func Progress(a ...interface{}) {
	DefaultLog.Progress(a...)
}

func Debug(a ...interface{}) {
	DefaultLog.Debug(a...)
}

func Writer(level Level) io.Writer {
	return DefaultLog.Writer(level)
}

func Close() {
	DefaultLog.Close()
}

var (
	// target console
	DefaultConsoleColorize = !color.NoColor

	// target file
	DefaultFileMaxBytes    int64 = 100 * 1024 * 1024 // 100M
	DefaultFileBackupCount       = 30                // 30

	// target network
	DefaultNetworkType    = `tcp`
	DefaultNetworkAddress = ``

	// target mail
	DefaultMailHost       = ``
	DefaultMailUsername   = ``
	DefaultMailPassword   = ``
	DefaultMailSubject    = ``
	DefaultMailSender     = ``
	DefaultMailRecipients = []string{}
)

func UseCommonTargets(levelName string, targetNames ...string) *Logger {
	DefaultLog.SetLevel(levelName)
	targets := []Target{}

	for _, targetName := range targetNames {
		ti := strings.SplitN(targetName, `:`, 2)
		var categories []string
		if len(ti) == 2 {
			targetName = ti[0]
			if len(ti[1]) > 0 {
				categories = strings.Split(ti[1], `,`)
			}
		}
		switch targetName {
		case "console":
			//输出到命令行
			consoleTarget := NewConsoleTarget()
			consoleTarget.ColorMode = DefaultConsoleColorize
			consoleTarget.Categories = categories
			targets = append(targets, consoleTarget)

		case "file":
			//输出到文件
			if DefaultLog.MaxLevel.Int() >= LevelInfo.Int() {
				fileTarget := NewFileTarget()
				fileTarget.FileName = `logs/{date:20060102}_info.log`
				fileTarget.Levels = map[Leveler]bool{LevelInfo: true}
				fileTarget.Categories = categories
				fileTarget.MaxBytes = DefaultFileMaxBytes
				targets = append(targets, fileTarget)
			}
			if DefaultLog.MaxLevel.Int() >= LevelWarn.Int() {
				fileTarget := NewFileTarget()
				fileTarget.FileName = `logs/{date:20060102}_warn.log` //按天分割日志
				fileTarget.Levels = map[Leveler]bool{LevelWarn: true}
				fileTarget.Categories = categories
				fileTarget.MaxBytes = DefaultFileMaxBytes
				targets = append(targets, fileTarget)
			}
			if DefaultLog.MaxLevel.Int() >= LevelError.Int() {
				fileTarget := NewFileTarget()
				fileTarget.FileName = `logs/{date:20060102}_error.log` //按天分割日志
				fileTarget.MaxLevel = LevelError
				fileTarget.Categories = categories
				fileTarget.MaxBytes = DefaultFileMaxBytes
				targets = append(targets, fileTarget)
			}
			if DefaultLog.MaxLevel == LevelDebug {
				fileTarget := NewFileTarget()
				fileTarget.FileName = `logs/{date:20060102}_debug.log`
				fileTarget.Levels = map[Leveler]bool{LevelDebug: true}
				fileTarget.Categories = categories
				fileTarget.MaxBytes = DefaultFileMaxBytes
				targets = append(targets, fileTarget)
			}

		case "mail":
			if DefaultLog.MaxLevel.Int() == LevelFatal.Int() || DefaultLog.MaxLevel.Int() >= LevelError.Int() {
				mailTarget := NewMailTarget()
				mailTarget.MaxLevel = LevelError
				mailTarget.Categories = categories
				targets = append(targets, mailTarget)
			}

		case "network":
			netTarget := NewNetworkTarget()
			netTarget.Categories = categories
			targets = append(targets, netTarget)
		}
	}
	SetTarget(targets...)
	SetFatalAction(ActionExit)
	return DefaultLog.Logger
}
