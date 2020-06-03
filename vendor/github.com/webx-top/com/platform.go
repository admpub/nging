package com

import "runtime"

const (
	IsWindows = runtime.GOOS == "windows"
	IsLinux   = runtime.GOOS == "linux"
	IsMac     = runtime.GOOS == "darwin"
	Is32Bit   = runtime.GOARCH == "386"
	Is64Bit   = runtime.GOARCH == "amd64"
)
