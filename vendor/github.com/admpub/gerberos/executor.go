package gerberos

import (
	"bytes"
	"io"
	"os/exec"
)

type Executor interface {
	Execute(name string, args ...string) (string, int, error)
	ExecuteWithStd(stdin io.Reader, stdout io.Writer, name string, args ...string) (string, int, error)
}

func NewDefaultExecutor() *defaultExecutor {
	return &defaultExecutor{}
}

type defaultExecutor struct{}

func (e *defaultExecutor) Execute(name string, args ...string) (string, int, error) {
	return e.ExecuteWithStd(nil, nil, name, args...)
}

func (e *defaultExecutor) ExecuteWithStd(stdin io.Reader, stdout io.Writer, name string, args ...string) (string, int, error) {
	cmd := exec.Command(name, args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	var (
		b   []byte
		err error
	)
	if stdout == nil {
		b, err = cmd.CombinedOutput()
	} else {
		bf := &bytes.Buffer{}
		cmd.Stderr = bf
		b, err = bf.Bytes(), cmd.Run()
	}
	if err != nil {
		eerr, ok := err.(*exec.ExitError)
		if ok && eerr != nil {
			return string(b), eerr.ExitCode(), eerr
		}
		return "", -1, err
	}

	return string(b), 0, nil
}
