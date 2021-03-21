package log

import "io"

// DefaultLog 默认日志实例
var DefaultLog = &defaultLogger{Logger: New()}

type defaultLogger struct {
	*Logger
}

func GetLogger(category string, formatter ...Formatter) *Logger {
	return DefaultLog.GetLogger(category, formatter...)
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

func Infof(format string, a ...interface{}) {
	DefaultLog.Infof(format, a...)
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

func Info(a ...interface{}) {
	DefaultLog.Info(a...)
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

func UseCommonTargets(levelName string, targetNames ...string) *Logger {
	DefaultLog.SetLevel(levelName)
	targets := []Target{}

	for _, targetName := range targetNames {
		switch targetName {
		case "console":
			//输出到命令行
			consoleTarget := NewConsoleTarget()
			consoleTarget.ColorMode = false
			targets = append(targets, consoleTarget)

		case "file":
			//输出到文件
			if DefaultLog.MaxLevel.Int() >= LevelInfo.Int() {
				fileTarget := NewFileTarget()
				fileTarget.FileName = `logs/{date:20060102}_info.log`
				fileTarget.Filter.Levels = map[Leveler]bool{LevelInfo: true}
				fileTarget.MaxBytes = 10 * 1024 * 1024
				targets = append(targets, fileTarget)
			}
			if DefaultLog.MaxLevel.Int() >= LevelWarn.Int() {
				fileTarget := NewFileTarget()
				fileTarget.FileName = `logs/{date:20060102}_warn.log` //按天分割日志
				fileTarget.Filter.Levels = map[Leveler]bool{LevelWarn: true}
				fileTarget.MaxBytes = 10 * 1024 * 1024
				targets = append(targets, fileTarget)
			}
			if DefaultLog.MaxLevel.Int() >= LevelError.Int() {
				fileTarget := NewFileTarget()
				fileTarget.FileName = `logs/{date:20060102}_error.log` //按天分割日志
				fileTarget.Filter.MaxLevel = LevelError
				fileTarget.MaxBytes = 10 * 1024 * 1024
				targets = append(targets, fileTarget)
			}
			if DefaultLog.MaxLevel == LevelDebug {
				fileTarget := NewFileTarget()
				fileTarget.FileName = `logs/{date:20060102}_debug.log`
				fileTarget.Filter.Levels = map[Leveler]bool{LevelDebug: true}
				fileTarget.MaxBytes = 10 * 1024 * 1024
				targets = append(targets, fileTarget)
			}
		}
	}
	SetTarget(targets...)
	SetFatalAction(ActionExit)
	return DefaultLog.Logger
}
