//go:build windows

package goforever

import (
	"fmt"
	"syscall"
)

func (p *Process) setSysProcAttr(attr *syscall.SysProcAttr) error {
	// token, err := getToken(0)
	// if err != nil {
	// 	return err
	// }
	// attr.Token = token
	return nil
}

func getToken(pid int) (syscall.Token, error) {
	var token syscall.Token

	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return token, fmt.Errorf("Token Process Error: %w", err)
	}
	defer syscall.CloseHandle(handle)

	// Find process token via win32
	err = syscall.OpenProcessToken(handle, syscall.TOKEN_ALL_ACCESS, &token)
	if err != nil {
		return token, fmt.Errorf("Open Token Process Error: %w", err)
	}
	return token, err
}
