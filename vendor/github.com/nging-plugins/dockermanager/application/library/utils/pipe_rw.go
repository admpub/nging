package utils

import (
	"io"
	"sync/atomic"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/common"
)

func NewPipeRWBy(r *io.PipeReader, w *io.PipeWriter) *PipeRW {
	return &PipeRW{r: r, w: w, err: make(chan error)}
}

func NewPipe() *PipeRW {
	r, w := io.Pipe()
	return NewPipeRWBy(r, w)
}

type PipeRW struct {
	r      io.ReadCloser
	w      io.WriteCloser
	err    chan error
	hasErr int32
}

func (p *PipeRW) DoWrite(f func(io.Writer) error) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Errorf(`%v`, err)
			}
		}()
		err := f(p.w)
		p.w.Close()
		p.SendErr(err)
	}()
}

func (p *PipeRW) DoRead(f func(io.Reader) error) error {
	err := f(p.r)
	if err != nil {
		return err
	}
	err = <-p.Err()
	return err
}

func (p *PipeRW) Close() error {
	errs := common.NewErrors()
	err1 := p.r.Close()
	err2 := p.w.Close()
	if err1 != nil {
		errs.Add(err1)
	}
	if err2 != nil {
		errs.Add(err2)
	}
	close(p.err)
	return errs.ToError()
}

func (p *PipeRW) SendErr(err error) {
	atomic.AddInt32(&p.hasErr, 1)
	p.err <- err
}

func (p *PipeRW) Err() <-chan error {
	v := atomic.SwapInt32(&p.hasErr, 0)
	if v == 0 {
		r := make(chan error)
		close(r)
		return r
	}
	return p.err
}
