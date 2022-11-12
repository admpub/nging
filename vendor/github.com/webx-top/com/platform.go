package com

import (
	"log"
	"os"
	"runtime"
	"time"
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

func Getenv(key string, defaults ...string) string {
	v := os.Getenv(key)
	if len(v) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return v
}

func GetenvInt(key string, defaults ...int) int {
	v := os.Getenv(key)
	if len(v) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return Int(v)
}

func GetenvUint(key string, defaults ...uint) uint {
	v := os.Getenv(key)
	if len(v) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return Uint(v)
}

func GetenvInt64(key string, defaults ...int64) int64 {
	v := os.Getenv(key)
	if len(v) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return Int64(v)
}

func GetenvUint64(key string, defaults ...uint64) uint64 {
	v := os.Getenv(key)
	if len(v) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return Uint64(v)
}

func GetenvInt32(key string, defaults ...int32) int32 {
	v := os.Getenv(key)
	if len(v) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return Int32(v)
}

func GetenvUint32(key string, defaults ...uint32) uint32 {
	v := os.Getenv(key)
	if len(v) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return Uint32(v)
}

func GetenvFloat32(key string, defaults ...float32) float32 {
	v := os.Getenv(key)
	if len(v) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return Float32(v)
}

func GetenvFloat64(key string, defaults ...float64) float64 {
	v := os.Getenv(key)
	if len(v) == 0 && len(defaults) > 0 {
		return defaults[0]
	}
	return Float64(v)
}

func GetenvDuration(key string, defaults ...time.Duration) time.Duration {
	v := os.Getenv(key)
	if len(v) > 0 {
		t, err := time.ParseDuration(v)
		if err == nil {
			return t
		}
		log.Printf("GetenvDuration: %v: %v\n", v, err)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return 0
}
