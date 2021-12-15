package log

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

var (
	DefaultStackDepth  = 5
	DefaultSkipStack   = 3
	DefaultStackFilter = `github.com/admpub/log`
)

// CallStack 调用栈信息
type CallStack struct {
	Depth   int
	Skip    int
	Filters []string
}

// GetCallStack returns the current call stack information as a string.
// The skip parameter specifies how many top frames should be skipped, while
// the frames parameter specifies at most how many frames should be returned.
func GetCallStack(skip int, frames int, filters ...string) string {
	buf := new(bytes.Buffer)
	hasFilter := len(filters) > 0
	for i, count := skip, 0; count < frames; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if !hasFilter || !matchFilters(file, filters) {
			fmt.Fprintf(buf, "\n%s:%d", file, line)
			count++
		}
	}
	return buf.String()
}

func matchFilters(source string, filters []string) bool {
	for _, filter := range filters {
		if len(filter) == 0 {
			continue
		}
		if strings.Contains(source, filter) {
			return true
		}
	}
	return false
}

// GetCallSingleStack 获取单个记录
func GetCallSingleStack(skip int, filters ...string) (fileName string, lineNo int, found bool) {
	hasFilter := len(filters) > 0
	for i := skip; ; i++ {
		_, file, line, ok := runtime.Caller(i)
		if !ok {
			return
		}
		if !hasFilter || !matchFilters(file, filters) {
			fileName = file
			lineNo = line
			found = ok
			return
		}
	}
}
