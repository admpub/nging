package winpty

import (
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func UTF16PtrToString(p *uint16) string {
	var (
		sizeTest uint16
		finalStr []uint16 = make([]uint16, 0)
	)
	for {
		if *p == uint16(0) {
			break
		}

		finalStr = append(finalStr, *p)
		p = (*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + unsafe.Sizeof(sizeTest)))
	}
	return string(utf16.Decode(finalStr[0:]))
}

func UTF16PtrFromStringArray(s []string) (*uint16, error) {
	var r []uint16

	for _, ss := range s {
		a, err := syscall.UTF16FromString(ss)
		if err != nil {
			return nil, err
		}

		r = append(r, a...)
	}

	r = append(r, 0)

	return &r[0], nil
}

func GetErrorMessage(ptr uintptr) string {
	msgPtr, _, _ := winpty_error_msg.Call(ptr)
	if msgPtr == uintptr(0) {
		return "Unknown Error"
	}
	return UTF16PtrToString((*uint16)(unsafe.Pointer(msgPtr)))
}
