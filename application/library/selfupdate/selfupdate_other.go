//go:build !linux && !windows

package selfupdate

import "syscall"

func sysProcAttr() syscall.SysProcAttr {
	return syscall.SysProcAttr{
		Setpgid: true,
	}
}
