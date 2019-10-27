// +build !go1.13

package middleware

import (
	"github.com/admpub/log"
)

func init() {
	DefaultLogWriter = log.Writer(log.LevelInfo)
}
