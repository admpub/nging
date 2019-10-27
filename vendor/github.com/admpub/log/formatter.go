package log

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"time"
)

// Formatter formats a log message into an appropriate string.
type Formatter func(*Logger, *Entry) string

// DefaultFormatter is the default formatter used to format every log message.
func DefaultFormatter(l *Logger, e *Entry) string {
	return e.Time.Format(time.RFC3339) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
}

func NormalFormatter(l *Logger, e *Entry) string {
	return e.Time.Format(`2006-01-02 15:04:05`) + "|" + e.Level.String() + "|" + e.Category + "|" + e.Message + e.CallStack
}

func ShortFileFormatter(l *Logger, e *Entry) string {
	callDepth := DefaultCallDepth
	if l.CallDepth > 0 {
		callDepth = l.CallDepth
	}
	_, file, line, ok := runtime.Caller(callDepth)
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
