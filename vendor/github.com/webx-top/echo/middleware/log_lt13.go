// +build !go1.13

package middleware

import (
	"io"

	"github.com/admpub/log"
)

func GetDefaultLogWriter() io.Writer {
	return log.Writer(log.LevelInfo)
}
