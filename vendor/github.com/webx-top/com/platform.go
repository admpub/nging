package com

import (
	"os"
	"runtime"
)

const (
	IsWindows = runtime.GOOS == "windows"
	IsLinux   = runtime.GOOS == "linux"
	IsMac     = runtime.GOOS == "darwin"
	Is32Bit   = runtime.GOARCH == "386"
	Is64Bit   = runtime.GOARCH == "amd64"
)

// ExitOnSuccess 成功时退出程序
func ExitOnSuccess(msg string) {
	os.Stdout.WriteString(msg)
	os.Exit(0)
}

// ExitOnFailure 失败时退出程序
func ExitOnFailure(msg string, errCodes ...int) {
	errCode := 1
	if len(errCodes) > 0 {
		errCode = errCodes[0]
	}
	os.Stderr.WriteString(msg)
	os.Exit(errCode)
}
