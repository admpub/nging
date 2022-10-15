//go:build !windows
// +build !windows

package gopty

import (
	"os"
	"os/exec"
	"strings"

	"github.com/creack/pty"

	"github.com/admpub/gopty/interfaces"
)

var _ interfaces.Console = (*consoleNix)(nil)

type consoleNix struct {
	file *os.File
	cmd  *exec.Cmd

	initialCols int
	initialRows int

	cwd string
	env []string
}

func newNative(cols int, rows int) (Console, error) {
	cwd, _ := os.UserHomeDir()
	if len(cwd) == 0 {
		cwd = `.`
	}
	env := os.Environ()
	var hasEnvTERM bool
	for _, ev := range env {
		if strings.HasPrefix(ev, `TERM=`) {
			hasEnvTERM = true
			break
		}
	}
	if !hasEnvTERM {
		term := os.Getenv(`TERM`)
		if len(term) == 0 {
			term = `xterm`
		}
		env = append(env, `TERM=`+term)
	}
	return &consoleNix{
		initialCols: cols,
		initialRows: rows,

		file: nil,

		cwd: cwd,
		env: env,
	}, nil
}

// Start starts a process and wraps in a console
func (c *consoleNix) Start(args []string) error {
	cmd, err := c.buildCmd(args)
	if err != nil {
		return err
	}
	c.cmd = cmd

	cmd.Dir = c.cwd
	cmd.Env = c.env

	f, err := pty.StartWithSize(cmd, &pty.Winsize{Cols: uint16(c.initialCols), Rows: uint16(c.initialRows)})
	if err != nil {
		return err
	}

	c.file = f
	return nil
}

func (c *consoleNix) buildCmd(args []string) (*exec.Cmd, error) {
	if len(args) < 1 {
		return nil, ErrInvalidCmd
	}
	cmd := exec.Command(args[0], args[1:]...)
	return cmd, nil
}

func (c *consoleNix) Read(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.file.Read(b)
}

func (c *consoleNix) Write(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.file.Write(b)
}

func (c *consoleNix) Close() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	return c.file.Close()
}

func (c *consoleNix) SetSize(cols int, rows int) error {
	if c.file == nil {
		c.initialRows = rows
		c.initialCols = cols
		return nil
	}

	return pty.Setsize(c.file, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

func (c *consoleNix) GetSize() (int, int, error) {
	if c.file == nil {
		return c.initialCols, c.initialRows, nil
	}

	rows, cols, err := pty.Getsize(c.file)
	return cols, rows, err
}

func (c *consoleNix) Wait() (*os.ProcessState, error) {
	if c.cmd == nil {
		return nil, ErrProcessNotStarted
	}

	return c.cmd.Process.Wait()
}

func (c *consoleNix) SetCWD(cwd string) error {
	c.cwd = cwd
	return nil
}

func (c *consoleNix) SetENV(environ []string) error {
	c.env = append(os.Environ(), environ...)
	return nil
}

func (c *consoleNix) Pid() (int, error) {
	if c.cmd == nil {
		return 0, ErrProcessNotStarted
	}

	return c.cmd.Process.Pid, nil
}

func (c *consoleNix) Kill() error {
	if c.cmd == nil {
		return ErrProcessNotStarted
	}

	return c.cmd.Process.Kill()
}

func (c *consoleNix) Signal(sig os.Signal) error {
	if c.cmd == nil {
		return ErrProcessNotStarted
	}

	return c.cmd.Process.Signal(sig)
}
