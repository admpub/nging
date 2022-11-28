//go:build conpty
// +build conpty

package gopty

import (
	"github.com/admpub/conpty"
)

func IsConPtyAvailable() bool {
	return conpty.IsConPtyAvailable()
}

var _ interfaces.Console = (*conPtyWindows)(nil)

type conPtyWindows struct {
	initialCols int
	initialRows int

	file *conpty.ConPty

	cwd string
	env []string
}

func newNative(cols int, rows int) (Console, error) {
	cwd, _ := os.UserHomeDir()
	if len(cwd) == 0 {
		cwd = `.`
	}
	return &consoleWindows{
		initialCols: cols,
		initialRows: rows,

		file: nil,

		cwd: cwd,
		env: os.Environ(),
	}, nil
}

func (c *consoleWindows) Start(args []string) error {
	command := strings.Join(args, " ")
	cmd, err := conpty.Start(command, conpty.ConPtyDimensions(c.initialCols, c.initialRows))
	if err != nil {
		return err
	}

	c.file = cmd
	return nil
}

func (c *consoleWindows) Read(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	n, err := c.file.StdOut.Read(b)

	return n, err
}

func (c *consoleWindows) Write(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.file.StdIn.Write(b)
}

func (c *consoleWindows) Close() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	c.file.Close()
	return nil
}

func (c *consoleWindows) SetSize(cols int, rows int) error {
	c.initialRows = rows
	c.initialCols = cols

	if c.file == nil {
		return nil
	}

	c.file.Resize(c.initialCols, c.initialRows)
	return nil
}

func (c *consoleWindows) GetSize() (int, int, error) {
	return c.initialCols, c.initialRows, nil
}

func (c *consoleWindows) Wait() (*os.ProcessState, error) {
	if c.file == nil {
		return nil, ErrProcessNotStarted
	}

	_, err := c.file.Wait()
	if err != nil {
		return nil, err
	}

	return os.FindProcess(int(c.file.ProcessInformation().ProcessId))
}

func (c *consoleWindows) SetCWD(cwd string) error {
	c.cwd = cwd
	return nil
}

func (c *consoleWindows) SetENV(environ []string) error {
	c.env = append(os.Environ(), environ...)
	return nil
}

func (c *consoleWindows) Pid() (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return int(c.file.ProcessInformation().ProcessId), err
}

func (c *consoleWindows) Kill() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	proc, err := os.FindProcess(int(c.file.ProcessInformation().ProcessId))
	if err != nil {
		return err
	}

	return proc.Kill()
}

func (c *consoleWindows) Signal(sig os.Signal) error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	proc, err := os.FindProcess(int(c.file.ProcessInformation().ProcessId))
	if err != nil {
		return err
	}

	return proc.Signal(sig)
}
