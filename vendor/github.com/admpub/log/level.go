package log

import (
	"fmt"
	"strings"

	"github.com/admpub/color"
)

// RFC5424 log message levels.
const (
	LevelFatal Level = iota
	LevelError
	LevelWarn
	LevelOkay
	LevelInfo
	LevelProgress
	LevelDebug
)

type (
	// Level 日志等级编号
	Level int

	// Action 日志触发行为编号
	Action int

	// Leveler 日志等级接口
	Leveler interface {
		fmt.Stringer
		IsEnabled(level Level) bool
		Int() int
		Tag() string
		Color() *color.Color
		Level() Level
	}
)

var (
	// LevelNames maps log levels to names
	LevelNames = map[Leveler]string{
		LevelDebug:    "Debug",
		LevelProgress: "Progress",
		LevelInfo:     "Info",
		LevelOkay:     "Okay",
		LevelWarn:     "Warn",
		LevelError:    "Error",
		LevelFatal:    "Fatal",
	}

	// LevelUppers 日志大写名称前缀
	LevelUppers = map[string]string{
		`Debug`:    "DEBUG",
		"Progress": "PROGR",
		`Info`:     " INFO",
		"Okay":     " OKAY",
		`Warn`:     " WARN",
		`Error`:    "ERROR",
		`Fatal`:    "FATAL",
	}

	// Levels 所有日志等级
	Levels = map[string]Leveler{
		"Debug":    LevelDebug,
		"Progress": LevelProgress,
		"Info":     LevelInfo,
		"Okay":     LevelOkay,
		"Warn":     LevelWarn,
		"Error":    LevelError,
		"Fatal":    LevelFatal,
	}
)

// HTTPStatusLevelName HTTP状态码相应级别名称
func HTTPStatusLevelName(httpCode int) string {
	s := `Info`
	switch {
	case httpCode >= 500:
		s = `Error`
	case httpCode >= 400:
		s = `Warn`
	case httpCode >= 300:
		s = `Debug`
	}
	return s
}

// GetLevel 获取指定名称的等级信息
func GetLevel(level string) (Leveler, bool) {
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

// IsEnabled 是否启用了某个等级
func (l Level) IsEnabled(level Level) bool {
	return l >= level
}

// Int 等级数值
func (l Level) Int() int {
	return int(l)
}

func (l Level) Level() Level {
	return l
}

// Tag 标签
func (l Level) Tag() string {
	return `[` + LevelUppers[l.String()] + `]`
}

// Color 颜色
func (l Level) Color() *color.Color {
	return colorBrushes[l]
}
