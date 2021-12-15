package engine

import (
	"strconv"
	"strings"
	"unsafe"
)

var (
	HeaderSetCookie = `Set-Cookie`
)

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func AddressPort(address string) int {
	delim := `:`
	if len(address) > 0 && address[0] == '[' {
		delim = `]:`
	}
	parts := strings.SplitN(address, delim, 2)
	if len(parts) > 1 {
		port, _ := strconv.Atoi(parts[1])
		return port
	}
	return 80
}
