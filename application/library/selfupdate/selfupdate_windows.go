//go:build windows

package selfupdate

import "syscall"

func sysProcAttr() syscall.SysProcAttr {
	return syscall.SysProcAttr{}
}
