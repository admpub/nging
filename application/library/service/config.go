package service

import (
	"io"

	"github.com/admpub/service"
)

// Config is the runner app config structure.
type Config struct {
	service.Config

	Dir  string
	Exec string
	Args []string
	Env  []string

	Stderr, Stdout io.Writer
}
