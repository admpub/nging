package goforever

import (
	"os"
)

type Processer interface {
	Kill() error
	Release() error
	Wait() (*os.ProcessState, error)
	Pid() int
}
