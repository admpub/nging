//go:build linux

package net

import (
	"context"
	reactor "github.com/godzie44/go-uring/reactor"
	"github.com/godzie44/go-uring/uring"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"syscall"
	"time"
)

// defaultTCPKeepAlive is a default constant value for TCPKeepAlive times
// See golang.org/issue/31510
const (
	defaultTCPKeepAlive = 15 * time.Second
)

//Listener tcp listener with uring reactor inside.
type Listener struct {
	lc net.ListenConfig

	sockFd int

	acceptChan chan uring.CQEvent
	acceptOp   *uring.AcceptOp

	reactor *reactor.NetworkReactor

	stopReactorFn func()

	addr net.Addr
}

func NewListener(lc net.ListenConfig, addr string, reactor *reactor.NetworkReactor) (*Listener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}

	sockFd, err := serverSocket(tcpAddr)
	if err != nil {
		return nil, err
	}

	l := &Listener{
		lc:         lc,
		sockFd:     sockFd,
		addr:       tcpAddr,
		reactor:    reactor,
		acceptChan: make(chan uring.CQEvent),
		acceptOp:   uring.Accept(uintptr(sockFd), 0),
	}

	ctx, cancel := context.WithCancel(context.Background())

	l.stopReactorFn = cancel
	go reactor.Run(ctx)

	return l, nil
}

func serverSocket(tcpAddr *net.TCPAddr) (int, error) {
	sockFd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC, 0)
	if err != nil {
		return 0, err
	}

	if err = setDefaultListenerSockopts(sockFd); err != nil {
		return 0, err
	}

	addr := syscall.SockaddrInet4{
		Port: tcpAddr.Port,
	}
	copy(addr.Addr[:], tcpAddr.IP.To4())

	if err = syscall.Bind(sockFd, &addr); err != nil {
		return 0, os.NewSyscallError("bind", err)
	}

	if err = syscall.Listen(sockFd, syscall.SOMAXCONN); err != nil {
		return 0, os.NewSyscallError("listen", err)
	}

	return sockFd, nil
}

func (l *Listener) Accept() (net.Conn, error) {
	l.reactor.Queue(l.acceptOp, func(event uring.CQEvent) {
		l.acceptChan <- event
	})
	cqe := <-l.acceptChan

	if err := cqe.Error(); err != nil {
		return nil, err
	}

	fd := int(cqe.Res)

	rAddr, _ := l.acceptOp.Addr()
	tc := newConn(fd, l.addr, rAddr, l.reactor)
	if l.lc.KeepAlive >= 0 {
		_ = setKeepAlive(fd, true)
		ka := l.lc.KeepAlive
		if l.lc.KeepAlive == 0 {
			ka = defaultTCPKeepAlive
		}
		_ = setKeepAlivePeriod(fd, ka)
	}
	return tc, nil
}

func setKeepAlive(fd int, keepalive bool) error {
	err := unix.SetsockoptInt(fd, unix.SOL_SOCKET, unix.SO_KEEPALIVE, boolint(keepalive))
	return wrapSyscallError("setsockopt", err)
}

func setKeepAlivePeriod(fd int, d time.Duration) error {
	secs := int(roundDurationUp(d, time.Second))
	if err := unix.SetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_KEEPINTVL, secs); err != nil {
		return wrapSyscallError("setsockopt", err)
	}
	err := unix.SetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_KEEPIDLE, secs)
	return wrapSyscallError("setsockopt", err)
}

func setDefaultListenerSockopts(s int) error {
	err := os.NewSyscallError("setsockopt", syscall.SetsockoptInt(s, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1))
	if err != nil {
		return err
	}

	return os.NewSyscallError("setsockopt", syscall.SetNonblock(s, false))
}

func roundDurationUp(d time.Duration, to time.Duration) time.Duration {
	return (d + to - 1) / to
}

func wrapSyscallError(name string, err error) error {
	if _, ok := err.(syscall.Errno); ok {
		err = os.NewSyscallError(name, err)
	}
	return err
}

func boolint(b bool) int {
	if b {
		return 1
	}
	return 0
}

func (l *Listener) Close() (err error) {
	err = syscall.Close(l.sockFd)
	l.stopReactorFn()
	return err
}

func (l *Listener) Addr() net.Addr {
	return l.addr
}
