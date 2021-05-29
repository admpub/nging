// +build go1.13

package middleware

import (
	"io"
	"log"
)

func GetDefaultLogWriter() io.Writer {
	return log.Writer()
}
