//go:build linux

package uring

import (
	"math"
	"os"
	"syscall"
	"unsafe"
)

const (
	sysRingSetup    uintptr = 425
	sysRingEnter    uintptr = 426
	sysRingRegister uintptr = 427

	//copied from signal_unix.numSig
	numSig = 65
)

// sqRing ring flags
const (
	sqNeedWakeup uint32 = 1 << 0 // needs io_uring_enter wakeup
	sqCQOverflow uint32 = 1 << 1 // cq ring is overflown
)

// io_uring_enter flags
const (
	sysRingEnterGetEvents uint32 = 1 << 0
	sysRingEnterSQWakeup  uint32 = 1 << 1
	sysRingEnterSQWait    uint32 = 1 << 2
	sysRingEnterExtArg    uint32 = 1 << 3
)

// SQE flags
const (
	SqeFixedFileFlag    uint8 = 1 << 0
	SqeIODrainFlag      uint8 = 1 << 1
	SqeIOLinkFlag       uint8 = 1 << 2
	SqeIOHardLinkFlag   uint8 = 1 << 3
	SqeAsyncFlag        uint8 = 1 << 4
	SqeBufferSelectFlag uint8 = 1 << 5
)

const (
	libUserDataTimeout = math.MaxUint64
)

func sysEnter(ringFD int, toSubmit uint32, minComplete uint32, flags uint32, sig unsafe.Pointer, raw bool) (uint, error) {
	return sysEnter2(ringFD, toSubmit, minComplete, flags, sig, numSig/8, raw)
}

func sysEnter2(ringFD int, toSubmit uint32, minComplete uint32, flags uint32, sig unsafe.Pointer, sz int, raw bool) (uint, error) {
	var consumed uintptr
	var errno syscall.Errno

	if raw {
		consumed, _, errno = syscall.RawSyscall6(
			sysRingEnter,
			uintptr(ringFD),
			uintptr(toSubmit),
			uintptr(minComplete),
			uintptr(flags),
			uintptr(sig),
			uintptr(sz),
		)
	} else {
		consumed, _, errno = syscall.Syscall6(
			sysRingEnter,
			uintptr(ringFD),
			uintptr(toSubmit),
			uintptr(minComplete),
			uintptr(flags),
			uintptr(sig),
			uintptr(sz),
		)
	}

	if errno != 0 {
		return 0, os.NewSyscallError("io_uring_enter", errno)
	}

	return uint(consumed), nil
}

func sysSetup(entries uint32, params *ringParams) (int, error) {
	fd, _, errno := syscall.Syscall(sysRingSetup, uintptr(entries), uintptr(unsafe.Pointer(params)), 0)
	if errno != 0 {
		return int(fd), os.NewSyscallError("io_uring_setup", errno)
	}

	return int(fd), nil
}

func sysRegister(ringFD int, op int, arg unsafe.Pointer, nrArgs int) error {
	_, _, errno := syscall.Syscall6(
		sysRingRegister,
		uintptr(ringFD),
		uintptr(op),
		uintptr(arg),
		uintptr(nrArgs),
		0,
		0,
	)
	if errno != 0 {
		return os.NewSyscallError("io_uring_register", errno)
	}
	return nil
}

type SQEntry struct {
	OpCode      uint8
	Flags       uint8
	IoPrio      uint16
	Fd          int32
	Off         uint64
	Addr        uint64
	Len         uint32
	OpcodeFlags uint32
	UserData    uint64

	BufIG       uint16
	Personality uint16
	SpliceFdIn  int32
	_pad2       [2]uint64
}

//go:uintptrescapes
func (sqe *SQEntry) fill(op OpCode, fd int32, addr uintptr, len uint32, offset uint64) {
	sqe.OpCode = uint8(op)
	sqe.Flags = 0
	sqe.IoPrio = 0
	sqe.Fd = fd
	sqe.Off = offset
	setAddr(sqe, addr)
	sqe.Len = len
	sqe.OpcodeFlags = 0
	sqe.UserData = 0
	sqe.BufIG = 0
	sqe.Personality = 0
	sqe.SpliceFdIn = 0
	sqe._pad2[0] = 0
	sqe._pad2[1] = 0
}

func (sqe *SQEntry) setUserData(ud uint64) {
	sqe.UserData = ud
}

//go:uintptrescapes
func setAddr(sqe *SQEntry, addr uintptr) {
	sqe.Addr = uint64(addr)
}

type CQEvent struct {
	UserData uint64
	Res      int32
	Flags    uint32
}

func (cqe *CQEvent) Error() error {
	if cqe.Res < 0 {
		return syscall.Errno(uintptr(-cqe.Res))
	}
	return nil
}

type getEventsArg struct {
	sigMask   uintptr
	sigMaskSz uint32
	_pad      uint32
	ts        uintptr
}

//go:uintptrescapes
func newGetEventsArg(sigMask uintptr, sigMaskSz uint32, ts uintptr) *getEventsArg {
	return &getEventsArg{sigMask: sigMask, sigMaskSz: sigMaskSz, ts: ts}
}
