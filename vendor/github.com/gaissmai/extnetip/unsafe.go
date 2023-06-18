package extnetip

import (
	"net/netip"
	"unsafe"
)

// addr is a struct for unsafe peeks into netip.Addr for uint128 math calculations.
type addr struct {
	ip uint128
	z  uintptr
}

// peek the singleton pointer for z4 from netip
var z4 uintptr = peek(netip.AddrFrom4([4]byte{})).z

// peek into the private internals of netip.Addr with unsafe.Pointer
func peek(a netip.Addr) addr {
	return *(*addr)(unsafe.Pointer(&a))
}

// back conversion to netip.Addr
func back(a addr) netip.Addr {
	return *(*netip.Addr)(unsafe.Pointer(&a))
}
