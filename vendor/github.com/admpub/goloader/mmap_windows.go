// +build windows

package goloader

import (
	"os"
	"syscall"
	"unsafe"
)

func Mmap(size int) ([]byte, error) {

	sizelo := uint32(size >> 32)
	sizehi := uint32(size) & 0xFFFFFFFF
	h, errno := syscall.CreateFileMapping(syscall.InvalidHandle, nil,
		syscall.PAGE_EXECUTE_READWRITE, sizelo, sizehi, nil)
	if h == 0 {
		return nil, os.NewSyscallError("CreateFileMapping", errno)
	}

	addr, errno := syscall.MapViewOfFile(h,
		syscall.FILE_MAP_READ|syscall.FILE_MAP_WRITE|syscall.FILE_MAP_EXECUTE,
		0, 0, uintptr(size))
	if addr == 0 {
		return nil, os.NewSyscallError("MapViewOfFile", errno)
	}

	if err := syscall.CloseHandle(syscall.Handle(h)); err != nil {
		return nil, os.NewSyscallError("CloseHandle", err)
	}

	var header sliceHeader
	header.Data = addr
	header.Len = size
	header.Cap = size
	b := *(*[]byte)(unsafe.Pointer(&header))

	return b, nil
}

func Munmap(b []byte) error {

	addr := (uintptr)(unsafe.Pointer(&b[0]))
	if err := syscall.UnmapViewOfFile(addr); err != nil {
		return os.NewSyscallError("UnmapViewOfFile", err)
	}
	return nil
}
