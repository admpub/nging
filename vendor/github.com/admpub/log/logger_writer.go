package log

import (
	"fmt"
	"strconv"
	"strings"
)

type LoggerWriter struct {
	Level Leveler
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
	level, s := l.detectLevel(s)
	l.Logger.Log(level, s)
	return
}

func (l *LoggerWriter) detectLevel(s string) (Leveler, string) {
	level := l.Level
	if len(s) <= 6 {
		return level, s
	}
	switch s[0] {
	case '>': // stdLog.Println(`>Debug:message`)
		pos := strings.Index(s, `:`)
		if pos >= 0 {
			levelText := s[1:pos]
			if lv, ok := Levels[levelText]; ok {
				level = lv
				s = s[pos+1:]
			}
		}
	case ':': // stdLog.Println(`:200:message`)
		s2 := s[1:]
		pos := strings.Index(s2, `:`)
		if pos >= 0 {
			httpCode := s2[0:pos]
			code, err := strconv.Atoi(httpCode)
			if err != nil {
				return level, s
			}
			level = NewHttpLevel(code, l.Level)
		}
	}
	return level, s
}

func (l *LoggerWriter) Printf(format string, v ...interface{}) {
	level, format := l.detectLevel(format)
	l.Logger.Logf(level, format, v...)
}

func (l *LoggerWriter) Println(v ...interface{}) {
	level := l.Level
	if len(v) > 0 {
		level, v[0] = l.detectLevel(fmt.Sprint(v[0]))
	}
	l.Logger.Log(level, v...)
}
