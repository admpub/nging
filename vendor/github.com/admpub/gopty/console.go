package gopty

import (
	"github.com/admpub/gopty/interfaces"
)

// Console communication interface
type Console interfaces.Console

// New creates a new console with initial size
func New(w int, h int) (Console, error) {
	return newNative(w, h)
}
