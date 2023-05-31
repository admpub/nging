package cmder

import "io"

func NewSimple() *Simple {
	return &Simple{}
}

var _ Cmder = &Simple{}

type Simple struct {
}

func (s *Simple) Boot() error {
	return nil
}

func (s *Simple) StopHistory(...string) error {
	return nil
}

func (s *Simple) Start(writer ...io.Writer) error {
	return nil
}

func (s *Simple) Stop() error {
	return nil
}

func (s *Simple) Reload() error {
	return nil
}

func (s *Simple) Restart(writer ...io.Writer) error {
	return nil
}
