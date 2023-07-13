//go:build linux

package uring

import (
	"github.com/libp2p/go-sockaddr"
	sockaddrnet "github.com/libp2p/go-sockaddr/net"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
	"time"
	"unsafe"
)

type OpCode uint8

const (
	NopCode OpCode = iota
	ReadVCode
	WriteVCode
	opFSync
	opReadFixed
	opWriteFixed
	opPollAdd
	opPollRemove
	opSyncFileRange
	opSendMsg
	opRecvMsg
	TimeoutCode
	opTimeoutRemove
	AcceptCode
	AsyncCancelCode
	LinkTimeoutCode
	ConnectCode
	opFAllocate
	opOpenAt
	CloseCode
	opFilesUpdate
	opStatX
	ReadCode
	WriteCode
	opFAdvise
	opMAdvise
	SendCode
	RecvCode
	openAt2Code
	epollCtlCode
	spliceCode
	ProvideBuffersCode
	removeBuffersCode
	teeCode
	shutdownCode
	renameAtCode
	unlinkAtCode
	mkdirAtCode
	symlinkAtCode
	linkAtCode
)

//NopOp - do not perform any I/O. This is useful for testing the performance of the io_uring implementation itself.
type NopOp struct {
}

func Nop() *NopOp {
	return &NopOp{}
}

func (op *NopOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(NopCode, -1, uintptr(unsafe.Pointer(nil)), 0, 0)
}

func (op *NopOp) Code() OpCode {
	return NopCode
}

//ReadVOp vectored read operation, similar to preadv2(2).
type ReadVOp struct {
	FD     uintptr
	Size   int64
	Offset uint64
	IOVecs []syscall.Iovec
}

//ReadV vectored read operation, similar to preadv2(2).
func ReadV(file *os.File, vectors [][]byte, offset uint64) *ReadVOp {
	buffs := make([]syscall.Iovec, len(vectors))
	for i, v := range vectors {
		buffs[i] = syscall.Iovec{
			Base: &v[0],
			Len:  uint64(len(v)),
		}
	}

	return &ReadVOp{FD: file.Fd(), IOVecs: buffs, Offset: offset}
}

func (op *ReadVOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(ReadVCode, int32(op.FD), uintptr(unsafe.Pointer(&op.IOVecs[0])), uint32(len(op.IOVecs)), op.Offset)
}

func (op *ReadVOp) Code() OpCode {
	return ReadVCode
}

//WriteVOp vectored write operation, similar to pwritev2(2).
type WriteVOp struct {
	FD     uintptr
	IOVecs []syscall.Iovec
	Offset uint64
}

//WriteV vectored writes bytes to file. Write starts from offset.
//If the file is not seekable, offset must be set to zero.
func WriteV(file *os.File, bytes [][]byte, offset uint64) *WriteVOp {
	buffs := make([]syscall.Iovec, len(bytes))
	for i := range bytes {
		buffs[i].SetLen(len(bytes[i]))
		buffs[i].Base = &bytes[i][0]
	}

	return &WriteVOp{FD: file.Fd(), IOVecs: buffs, Offset: offset}
}

func (op *WriteVOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(WriteVCode, int32(op.FD), uintptr(unsafe.Pointer(&op.IOVecs[0])), uint32(len(op.IOVecs)), op.Offset)
}

func (op *WriteVOp) Code() OpCode {
	return WriteVCode
}

//TimeoutOp timeout command.
type TimeoutOp struct {
	dur time.Duration
}

//Timeout - timeout operation.
func Timeout(duration time.Duration) *TimeoutOp {
	return &TimeoutOp{
		dur: duration,
	}
}

func (op *TimeoutOp) PrepSQE(sqe *SQEntry) {
	spec := syscall.NsecToTimespec(op.dur.Nanoseconds())
	sqe.fill(TimeoutCode, -1, uintptr(unsafe.Pointer(&spec)), 1, 0)
}

func (op *TimeoutOp) Code() OpCode {
	return TimeoutCode
}

//AcceptOp accept command.
type AcceptOp struct {
	fd    uintptr
	flags uint32
	addr  *unix.RawSockaddrAny
	len   uint32
}

//Accept - accept operation.
func Accept(fd uintptr, flags uint32) *AcceptOp {
	return &AcceptOp{
		addr:  &unix.RawSockaddrAny{},
		len:   unix.SizeofSockaddrAny,
		fd:    fd,
		flags: flags,
	}
}

func (op *AcceptOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(AcceptCode, int32(op.fd), uintptr(unsafe.Pointer(op.addr)), 0, uint64(uintptr(unsafe.Pointer(&op.len))))
	sqe.OpcodeFlags = op.flags
}

func (op *AcceptOp) Fd() int {
	return int(op.fd)
}

func (op *AcceptOp) Code() OpCode {
	return AcceptCode
}

func (op *AcceptOp) Addr() (net.Addr, error) {
	sAddr, err := sockaddr.AnyToSockaddr(op.addr)
	if err != nil {
		return nil, err
	}

	return sockaddrnet.SockaddrToTCPAddr(sAddr), nil
}

func (op *AcceptOp) AddrLen() uint32 {
	return op.len
}

//CancelOp attempt to cancel an already issued request.
type CancelOp struct {
	flags          uint32
	targetUserData uint64
}

//Cancel create CancelOp. Put in targetUserData value of user_data field of the request that should be cancelled.
func Cancel(targetUserData uint64, flags uint32) *CancelOp {
	return &CancelOp{flags: flags, targetUserData: targetUserData}
}

func (op *CancelOp) SetTargetUserData(ud uint64) {
	op.targetUserData = ud
}

func (op *CancelOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(AsyncCancelCode, int32(-1), uintptr(op.targetUserData), 0, 0)
	sqe.OpcodeFlags = op.flags
}

func (op *CancelOp) Code() OpCode {
	return AsyncCancelCode
}

//LinkTimeoutOp IORING_OP_LINK_TIMEOUT command.
type LinkTimeoutOp struct {
	dur time.Duration
}

//LinkTimeout - timeout operation for linked command.
//Note: previous queued SQE must be queued with flag SqeIOLinkFlag.
func LinkTimeout(duration time.Duration) *LinkTimeoutOp {
	return &LinkTimeoutOp{
		dur: duration,
	}
}

func (op *LinkTimeoutOp) PrepSQE(sqe *SQEntry) {
	spec := syscall.NsecToTimespec(op.dur.Nanoseconds())
	sqe.fill(LinkTimeoutCode, -1, uintptr(unsafe.Pointer(&spec)), 1, 0)
}

func (op *LinkTimeoutOp) Code() OpCode {
	return LinkTimeoutCode
}

//RecvOp receive a message from a socket operation.
type RecvOp struct {
	fd       uintptr
	buff     []byte
	msgFlags uint32
}

//Recv receive a message from a socket.
func Recv(socketFd uintptr, buff []byte, msgFlags uint32) *RecvOp {
	return &RecvOp{
		fd:       socketFd,
		buff:     buff,
		msgFlags: msgFlags,
	}
}

func (op *RecvOp) SetBuffer(buff []byte) {
	op.buff = buff
}

func (op *RecvOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(RecvCode, int32(op.fd), uintptr(unsafe.Pointer(&op.buff[0])), uint32(len(op.buff)), 0)
	sqe.OpcodeFlags = op.msgFlags
}

func (op *RecvOp) Fd() int {
	return int(op.fd)
}

func (op *RecvOp) Code() OpCode {
	return RecvCode
}

//SendOp send a message to a socket operation.
type SendOp struct {
	fd       uintptr
	buff     []byte
	msgFlags uint32
}

//Send send a message to a socket.
func Send(socketFd uintptr, buff []byte, msgFlags uint32) *SendOp {
	return &SendOp{
		fd:       socketFd,
		buff:     buff,
		msgFlags: msgFlags,
	}
}

func (op *SendOp) SetBuffer(buff []byte) {
	op.buff = buff
}

func (op *SendOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(SendCode, int32(op.fd), uintptr(unsafe.Pointer(&op.buff[0])), uint32(len(op.buff)), 0)
	sqe.OpcodeFlags = op.msgFlags
}

func (op *SendOp) Fd() int {
	return int(op.fd)
}

func (op *SendOp) Code() OpCode {
	return SendCode
}

//ProvideBuffersOp .
type ProvideBuffersOp struct {
	buff     []byte
	bufferId uint64
	groupId  uint16
}

//ProvideBuffers .
func ProvideBuffers(buff []byte, bufferId uint64, groupId uint16) *ProvideBuffersOp {
	return &ProvideBuffersOp{
		buff:     buff,
		bufferId: bufferId,
		groupId:  groupId,
	}
}

func (op *ProvideBuffersOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(ProvideBuffersCode, int32(1), uintptr(unsafe.Pointer(&op.buff[0])), uint32(len(op.buff)), op.bufferId)
	sqe.BufIG = op.groupId
}

func (op *ProvideBuffersOp) Code() OpCode {
	return ProvideBuffersCode
}

//CloseOp closes a file descriptor, equivalent of a close(2) system call.
type CloseOp struct {
	fd uintptr
}

//Close closes a file descriptor, equivalent of a close(2) system call.
func Close(fd uintptr) *CloseOp {
	return &CloseOp{
		fd: fd,
	}
}

func (op *CloseOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(CloseCode, int32(op.fd), 0, 0, 0)
}

func (op *CloseOp) Code() OpCode {
	return CloseCode
}

//ReadOp read operation, equivalent of a pread(2) system call.
type ReadOp struct {
	fd   uintptr
	buff []byte
	off  uint64
}

//Read - create read operation, equivalent of a pread(2) system call.
func Read(fd uintptr, buff []byte, offset uint64) *ReadOp {
	return &ReadOp{fd: fd, buff: buff, off: offset}
}

func (op *ReadOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(ReadCode, int32(op.fd), uintptr(unsafe.Pointer(&op.buff[0])), uint32(len(op.buff)), op.off)
}

func (op *ReadOp) Code() OpCode {
	return ReadCode
}

//WriteOp write operation, equivalent of a pwrite(2) system call.
type WriteOp struct {
	fd   uintptr
	buff []byte
	off  uint64
}

//Write - create write operation, equivalent of a pwrite(2) system call.
func Write(fd uintptr, buff []byte, offset uint64) *WriteOp {
	return &WriteOp{fd: fd, buff: buff, off: offset}
}

func (op *WriteOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(WriteCode, int32(op.fd), uintptr(unsafe.Pointer(&op.buff[0])), uint32(len(op.buff)), op.off)
}

func (op *WriteOp) Code() OpCode {
	return WriteCode
}

//ConnectOp connect operation, equivalent of a connect(2) system call.
type ConnectOp struct {
	fd   uintptr
	addr *sockaddrnet.RawSockaddrAny
	len  sockaddr.Socklen
}

//Connect operation, equivalent of a connect(2) system call.
func Connect(fd uintptr, addr *net.TCPAddr) *ConnectOp {
	sa := sockaddrnet.NetAddrToSockaddr(addr)
	rsa, l, err := sockaddr.SockaddrToAny(sa)
	if err != nil {
		panic(err)
	}

	return &ConnectOp{
		fd:   fd,
		addr: rsa,
		len:  l,
	}
}

func (op *ConnectOp) PrepSQE(sqe *SQEntry) {
	sqe.fill(ConnectCode, int32(op.fd), uintptr(unsafe.Pointer(op.addr)), 0, uint64(op.len))
}

func (op *ConnectOp) Code() OpCode {
	return ConnectCode
}
