package log

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"time"
)

// Formatter formats a log message into an appropriate string.
type Formatter func(*Logger, *Entry) string

// DefaultFormatter is the default formatter used to format every log message.
func DefaultFormatter(l *Logger, e *Entry) string {
	return l.EmojiOfLevel(e.Level.Level()) + strconv.Itoa(l.Pid()) + "|" + e.Time.Format(time.RFC3339) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
}

func EmptyFormatter(l *Logger, e *Entry) string {
	return e.Message
}

// NormalFormatter 标准格式
func NormalFormatter(l *Logger, e *Entry) string {
	return l.EmojiOfLevel(e.Level.Level()) + strconv.Itoa(l.Pid()) + "|" + e.Time.Format(`2006-01-02 15:04:05`) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
}

// ShortFileFormatter 简介文件名格式
func ShortFileFormatter(skipStack int, filters ...string) Formatter {
	_filters := []string{DefaultStackFilter}
	if len(_filters) > 0 {
		_filters = append(_filters, filters...)
	}
	return func(l *Logger, e *Entry) string {
		file, line, ok := GetCallSingleStack(skipStack, _filters...)
		if !ok {
			return l.EmojiOfLevel(e.Level.Level()) + strconv.Itoa(l.Pid()) + "|" + e.Time.Format(`2006-01-02 15:04:05`) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
		}
		return l.EmojiOfLevel(e.Level.Level()) + strconv.Itoa(l.Pid()) + "|" + e.Time.Format(`2006-01-02 15:04:05`) + "|" + filepath.Base(file) + ":" + strconv.Itoa(line) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
	}
}

// JSONL json格式信息
type JSONL struct {
	Time      string          `bson:"time" json:"time"`
	Level     string          `bson:"level" json:"level"`
	Category  string          `bson:"category" json:"category"`
	Message   json.RawMessage `bson:"message" json:"message"`
	CallStack string          `bson:"callStack" json:"callStack"`
	Pid       int             `bson:"pid" json:"pid"`
}

// JSONFormatter json格式
func JSONFormatter(l *Logger, e *Entry) string {
	jsonl := &JSONL{
		Time:      e.Time.Format(`2006-01-02 15:04:05`),
		Level:     e.Level.String(),
		Category:  e.Category,
		Message:   []byte(`"` + e.Message + `"`),
		CallStack: e.CallStack,
		Pid:       l.Pid(),
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
