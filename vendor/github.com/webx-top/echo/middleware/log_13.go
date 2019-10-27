// +build go1.13

package middleware

import (
	"log"
)

func init() {
	DefaultLogWriter = log.Writer()
}
