package gerberos

import (
	"errors"
	"fmt"
	"os"
)

type Source interface {
	Initialize(r *Rule) error
	Matches() (chan *Match, error)
}

var sources = map[string]func() Source{
	"file":    NewFileSource,
	"systemd": NewSystemdSource,
	"kernel":  NewKernelSource,
	"test":    NewTestSource,
	"process": NewProcessSource,
}

func RegisterSource(name string, fn func() Source) {
	sources[name] = fn
}

func NewFileSource() Source {
	return &fileSource{}
}

func NewSystemdSource() Source {
	return &systemdSource{}
}

func NewKernelSource() Source {
	return &kernelSource{}
}

func NewTestSource() Source {
	return &testSource{}
}

func NewProcessSource() Source {
	return &processSource{}
}

type fileSource struct {
	Rule *Rule
	path string
}

func (s *fileSource) Initialize(r *Rule) error {
	s.Rule = r

	if len(r.Source) < 2 {
		return errors.New("missing path parameter")
	}
	s.path = r.Source[1]

	if fi, err := os.Stat(s.path); err == nil && fi.IsDir() {
		return fmt.Errorf(`"%s" is a directory`, s.path)
	}

	if len(r.Source) > 2 {
		return errors.New("superfluous parameter(s)")
	}

	return nil
}

func (s *fileSource) Matches() (chan *Match, error) {
	return s.Rule.ProcessScanner("tail", "-n", "0", "-F", s.path)
}

type systemdSource struct {
	Rule    *Rule
	service string
}

func (s *systemdSource) Initialize(r *Rule) error {
	s.Rule = r

	if len(r.Source) < 2 {
		return errors.New("missing service parameter")
	}
	s.service = r.Source[1]

	if len(r.Source) > 2 {
		return errors.New("superfluous parameter(s)")
	}

	return nil
}

func (s *systemdSource) Matches() (chan *Match, error) {
	return s.Rule.ProcessScanner("journalctl", "-n", "0", "-f", "-u", s.service)
}

type kernelSource struct {
	Rule *Rule
}

func (k *kernelSource) Initialize(r *Rule) error {
	k.Rule = r

	if len(r.Source) > 1 {
		return errors.New("superfluous parameter(s)")
	}

	return nil
}

func (k *kernelSource) Matches() (chan *Match, error) {
	return k.Rule.ProcessScanner("journalctl", "-kf", "-n", "0")
}

type testSource struct {
	Rule        *Rule
	matchesErr  error
	processPath string
}

func (s *testSource) Initialize(r *Rule) error {
	s.Rule = r

	return nil
}

func (s *testSource) Matches() (chan *Match, error) {
	if s.matchesErr != nil {
		return nil, s.matchesErr
	}

	p := "test/producer"
	if s.processPath != "" {
		p = s.processPath
	}
	return s.Rule.ProcessScanner(p)
}

type processSource struct {
	Rule *Rule
	name string
	args []string
}

func (s *processSource) Initialize(r *Rule) error {
	s.Rule = r

	if len(r.Source) < 2 {
		return errors.New("missing process name")
	}
	s.name = r.Source[1]

	s.args = r.Source[2:]

	return nil
}

func (s *processSource) Matches() (chan *Match, error) {
	return s.Rule.ProcessScanner(s.name, s.args...)
}
