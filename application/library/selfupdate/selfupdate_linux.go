//go:build linux

package selfupdate

import "syscall"

func sysProcAttr() syscall.SysProcAttr {
	return syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.Signal(0),
	}
}
