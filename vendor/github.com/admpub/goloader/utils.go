package goloader

import (
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"unsafe"
)

//go:linkname add runtime.add
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer

//go:linkname adduintptr runtime.add
func adduintptr(p uintptr, x int) unsafe.Pointer

func putUint24(b []byte, v uint32) {
	_ = b[2] // early bounds check to guarantee safety of writes below
	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
}

func alignof(i int, align int) int {
	if i%align != 0 {
		i = i + (align - i%align)
	}
	return i
}

func bytearrayAlign(b *[]byte, align int) {
	length := len(*b)
	if length%align != 0 {
		*b = append(*b, make([]byte, align-length%align)...)
	}
}

func putAddressAddOffset(b []byte, offset *int, addr uint64) {
	if PtrSize == Uint32Size {
		binary.LittleEndian.PutUint32(b[*offset:], uint32(addr))
	} else {
		binary.LittleEndian.PutUint64(b[*offset:], uint64(addr))
	}
	*offset = *offset + PtrSize
}

func putAddress(b []byte, addr uint64) {
	if PtrSize == Uint32Size {
		binary.LittleEndian.PutUint32(b, uint32(addr))
	} else {
		binary.LittleEndian.PutUint64(b, uint64(addr))
	}
}

// sign extend a 24-bit integer
func signext24(x int64) int32 {
	return (int32(x) << 8) >> 8
}

func copy2Slice(dst []byte, src uintptr, size int) {
	s := sliceHeader{
		Data: src,
		Len:  size,
		Cap:  size,
	}
	copy(dst, *(*[]byte)(unsafe.Pointer(&s)))
}

func append2Slice(dst *[]byte, src uintptr, size int) {
	s := sliceHeader{
		Data: src,
		Len:  size,
		Cap:  size,
	}
	*dst = append(*dst, *(*[]byte)(unsafe.Pointer(&s))...)
}

//go:nosplit
//go:noinline
//see runtime.internal.atomic.Loadp
func loadp(ptr unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(ptr)
}

func grow(bytes *[]byte, size int) {
	if len(*bytes) < size {
		*bytes = append(*bytes, make([]byte, size-len(*bytes))...)
	}
}

//see $GOROOT/src/cmd/internal/loader/loader.go:preprocess
func ispreprocesssymbol(name string) bool {
	if len(name) > 5 {
		switch name[:5] {
		case "$f32.", "$f64.", "$i64.":
			return true
		default:
		}
	}
	return false
}

func preprocesssymbol(name string, bytes []byte) error {
	val, err := strconv.ParseUint(name[5:], 16, 64)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to parse $-symbol %s: %v", name, err))
	}
	switch name[:5] {
	case "$f32.":
		if uint64(uint32(val)) != val {
			return errors.New(fmt.Sprintf("$-symbol %s too large: %d", name, val))
		}
		binary.LittleEndian.PutUint32(bytes, uint32(val))
		bytes = bytes[:4]
	case "$f64.", "$i64.":
		binary.LittleEndian.PutUint64(bytes, val)
	default:
		return errors.New(fmt.Sprintf("unrecognized $-symbol: %s", name))
	}
	return nil
}
