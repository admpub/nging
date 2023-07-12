//go:build linux

package net

import (
	"errors"
	"fmt"
	reactor "github.com/godzie44/go-uring/reactor"
	"github.com/godzie44/go-uring/uring"
	"io"
	"net"
	"os"
	"runtime"
	"sync"
	"syscall"
	"time"
)

type Conn struct {
	fd           int
	lAddr, rAddr net.Addr
	reactor      *reactor.NetworkReactor

	readChan, writeChan chan uring.CQEvent

	readDeadline, writeDeadline time.Time

	lastReadOpID, lastWriteOpID reactor.RequestID

	readOp *uring.RecvOp
	sendOp *uring.SendOp

	readLock, writeLock sync.Mutex
}

func newConn(fd int, lAddr, rAddr net.Addr, r *reactor.NetworkReactor) *Conn {
	conn := &Conn{
		lAddr:     lAddr,
		rAddr:     rAddr,
		fd:        fd,
		reactor:   r,
		readChan:  make(chan uring.CQEvent),
		writeChan: make(chan uring.CQEvent),
		readOp:    uring.Recv(uintptr(fd), nil, 0),
		sendOp:    uring.Send(uintptr(fd), nil, 0),
	}

	return conn
}

func (c *Conn) Read(b []byte) (n int, err error) {
	c.readLock.Lock()
	defer c.readLock.Unlock()

	op := c.readOp
	op.SetBuffer(b)

	c.lastReadOpID = c.reactor.QueueWithDeadline(op, func(event uring.CQEvent) {
		c.readChan <- event
	}, c.readDeadline)

	cqe := <-c.readChan

	if err = cqe.Error(); err != nil {
		if errors.Is(err, syscall.ECANCELED) {
			err = fmt.Errorf("%w: %s", os.ErrDeadlineExceeded, err.Error())
		}

		return 0, &net.OpError{Op: "read", Net: "tcp", Source: c.lAddr, Addr: c.rAddr, Err: err}
	}

	if cqe.Res == 0 {
		err = &net.OpError{Op: "read", Net: "tcp", Source: c.lAddr, Addr: c.rAddr, Err: io.EOF}
	}

	runtime.KeepAlive(b)
	return int(cqe.Res), err
}

func (c *Conn) Write(b []byte) (n int, err error) {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	op := c.sendOp
	op.SetBuffer(b)

	c.lastWriteOpID = c.reactor.QueueWithDeadline(op, func(event uring.CQEvent) {
		c.writeChan <- event
	}, c.writeDeadline)

	cqe := <-c.writeChan

	if err = cqe.Error(); err != nil {
		if errors.Is(err, syscall.ECANCELED) {
			err = fmt.Errorf("%w: %s", os.ErrDeadlineExceeded, err.Error())
		}

		return 0, &net.OpError{Op: "write", Net: "tcp", Source: c.lAddr, Addr: c.rAddr, Err: err}
	}
	if cqe.Res == 0 {
		err = &net.OpError{Op: "write", Net: "tcp", Source: c.lAddr, Addr: c.rAddr, Err: io.ErrUnexpectedEOF}
	}

	runtime.KeepAlive(b)
	return int(cqe.Res), err
}

func (c *Conn) Close() error {
	err := syscall.Close(c.fd)
	if err != nil {
		err = &net.OpError{Op: "close", Net: "tcp", Source: c.lAddr, Addr: c.rAddr, Err: err}
	}
	return err
}

func (c *Conn) LocalAddr() net.Addr {
	return c.lAddr
}

func (c *Conn) RemoteAddr() net.Addr {
	return c.rAddr
}

func (c *Conn) SetDeadline(t time.Time) error {
	_ = c.SetReadDeadline(t)
	_ = c.SetWriteDeadline(t)
	return nil
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	if !t.IsZero() && time.Until(t) < 0 {
		c.reactor.Cancel(c.lastReadOpID)
	}

	c.readDeadline = t
	return nil
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	if !t.IsZero() && time.Until(t) < 0 {
		c.reactor.Cancel(c.lastWriteOpID)
	}

	c.writeDeadline = t
	return nil
}

func (c *Conn) Fd() uintptr {
	return uintptr(c.fd)
}
