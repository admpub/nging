//go:build windows
// +build windows

package conpty

import (
	"golang.org/x/sys/windows"
)

func (cpty *ConPty) ProcessInformation() *windows.ProcessInformation {
	return cpty.pi
}
