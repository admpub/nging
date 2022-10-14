package gopty

import "errors"

var (
	ErrProcessNotStarted = errors.New("Process has not been started")
	ErrInvalidCmd        = errors.New("Invalid command")
)
