package packer

import (
	"os"

	oncex "github.com/admpub/once"
)

var (
	empty      Manager
	defaultMgr Manager
	defaultErr error
	once       oncex.Once
)

var (
	Stdout = os.Stdout
	Stderr = os.Stderr
)

func Default() (mgr Manager, err error) {
	once.Do(func() {
		defaultMgr, defaultErr = DetectManager()
	})
	mgr = defaultMgr
	err = defaultErr
	return
}

func Reset() {
	once.Reset()
}
