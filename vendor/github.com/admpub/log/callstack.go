package log

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"
)

type CallStack struct {
	Depth  int
	Filter string
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
