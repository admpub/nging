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
	LevelInfo
	LevelDebug
)

type (
	// Level describes the level of a log message.
	Level   int
	Action  int
	Leveler interface {
		fmt.Stringer
		Int() int
		Tag() string
		Color() *color.Color
	}
)

var (
	// LevelNames maps log levels to names
	LevelNames = map[Leveler]string{
		LevelDebug: "Debug",
		LevelInfo:  "Info",
		LevelWarn:  "Warn",
		LevelError: "Error",
		LevelFatal: "Fatal",
	}

	LevelUppers = map[string]string{
		`Debug`: "DEBUG",
		`Info`:  " INFO",
		`Warn`:  " WARN",
		`Error`: "ERROR",
		`Fatal`: "FATAL",
	}

	Levels = map[string]Leveler{
		"Debug": LevelDebug,
		"Info":  LevelInfo,
		"Warn":  LevelWarn,
		"Error": LevelError,
		"Fatal": LevelFatal,
	}

	DefaultCallDepth = 4
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

// Int 等级数值
func (l Level) Int() int {
	return int(l)
}

// Tag 标签
func (l Level) Tag() string {
	return `[` + LevelUppers[l.String()] + `]`
}

// Color 颜色
func (l Level) Color() *color.Color {
	return colorBrushes[l]
}
