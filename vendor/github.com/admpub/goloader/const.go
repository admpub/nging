package goloader

import "unsafe"

// size
const (
	PtrSize              = 4 << (^uintptr(0) >> 63)
	Uint32Size           = int(unsafe.Sizeof(uint32(0)))
	IntSize              = int(unsafe.Sizeof(int(0)))
	UInt64Size           = int(unsafe.Sizeof(uint64(0)))
	_FuncSize            = int(unsafe.Sizeof(_func{}))
	ItabSize             = int(unsafe.Sizeof(itab{}))
	FindFuncBucketSize   = int(unsafe.Sizeof(findfuncbucket{}))
	InlinedCallSize      = int(unsafe.Sizeof(inlinedCall{}))
	INVALID_HANDLE_VALUE = ^uintptr(0)
	INVALID_OFFSET       = int(-1)
)

const (
	TLSNAME        = "(TLS)"
	R_CALLIND_NAME = "R_CALLIND"
)

// cpu arch
const (
	ARCH_ARM32 = "arm"
	ARCH_ARM64 = "arm64"
	ARCH_386   = "386"
	ARCH_AMD64 = "amd64"
)

const (
	EMPTY_STRING    = ""
	DEFAULT_PKGPATH = "main"
	EMPTY_PKGPATH   = `""`
	ZERO_BYTE       = byte(0x00)
)

// runtime symbol
const (
	RUNTIME_DEFERRETURN = "runtime.deferreturn"
	RUNTIME_INIT        = "runtime.init"
)

// string match prefix/suffix
const (
	FILE_SYM_PREFIX        = "gofile.."
	TYPE_IMPORTPATH_PREFIX = "type..importpath."
	TYPE_DOUBLE_DOT_PREFIX = "type.."
	TYPE_PREFIX            = "type."
	ITAB_PREFIX            = "go.itab."
	RUNTIME_PREFIX         = "runtime."
	STKOBJ_SUFFIX          = ".stkobj"
	INLINETREE_SUFFIX      = ".inlinetree"
	OS_STDOUT              = "os.Stdout"
)
