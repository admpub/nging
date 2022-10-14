package gopty

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/admpub/gopty/interfaces"
	"github.com/iamacarpet/go-winpty"
)

// Do the interface allocations only once for common
// Errno values.
const (
	errnoERROR_IO_PENDING = 997
)

var (
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
)

//go:embed winpty/*
var winpty_deps embed.FS

// errnoErr returns common boxed Errno values, to prevent
// allocations at runtime.
func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return nil
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}

var _ interfaces.Console = (*consoleWindows)(nil)

type consoleWindows struct {
	initialCols int
	initialRows int

	file *winpty.WinPTY

	cwd string
	env []string
}

func newNative(cols int, rows int) (Console, error) {
	return &consoleWindows{
		initialCols: cols,
		initialRows: rows,

		file: nil,

		cwd: ".",
		env: os.Environ(),
	}, nil
}

func (c *consoleWindows) Start(args []string) error {
	dllDir, err := c.UnloadEmbeddedDeps()
	if err != nil {
		return err
	}

	opts := winpty.Options{
		DLLPrefix:   dllDir,
		InitialCols: uint32(c.initialCols),
		InitialRows: uint32(c.initialRows),
		Command:     strings.Join(args, " "),
		Dir:         c.cwd,
		Env:         c.env,
	}

	cmd, err := winpty.OpenWithOptions(opts)
	if err != nil {
		return err
	}

	c.file = cmd
	return nil
}

func (c *consoleWindows) UnloadEmbeddedDeps() (string, error) {

	executableName, err := os.Executable()
	if err != nil {
		return "", err
	}
	executableName = filepath.Base(executableName)

	dllDir := filepath.Join(os.TempDir(), fmt.Sprintf("%s_winpty", executableName))

	if err := os.MkdirAll(dllDir, 0755); err != nil {
		return "", err
	}

	files := []string{"winpty.dll", "winpty-agent.exe"}
	for _, file := range files {
		filenameEmbedded := fmt.Sprintf("winpty/%s", file)
		filenameDisk := filepath.Join(dllDir, file)

		_, statErr := os.Stat(filenameDisk)
		if statErr == nil {
			// file is already there, skip it
			continue
		}

		data, err := winpty_deps.ReadFile(filenameEmbedded)
		if err != nil {
			return "", err
		}

		if err := ioutil.WriteFile(filenameDisk, data, 0644); err != nil {
			return "", err
		}
	}

	return dllDir, nil
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

	c.file.SetSize(uint32(c.initialCols), uint32(c.initialRows))
	return nil
}

func (c *consoleWindows) GetSize() (int, int, error) {
	return c.initialCols, c.initialRows, nil
}

// At this point, sys/windows does not yet contain the method GetProcessID
// this was copied and pasted from: https://github.com/golang/sys/blob/master/windows/zsyscall_windows.go#L2226
func (c *consoleWindows) getProcessIDFromHandle(process uintptr) (id uint32, err error) {
	modkernel32 := syscall.NewLazyDLL("kernel32.dll")
	procGetProcessId := modkernel32.NewProc("GetProcessId")

	r0, _, e1 := syscall.Syscall(procGetProcessId.Addr(), 1, process, 0, 0)
	id = uint32(r0)
	if id == 0 {
		if e1 != 0 {
			err = errnoErr(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func (c *consoleWindows) Wait() (*os.ProcessState, error) {
	if c.file == nil {
		return nil, ErrProcessNotStarted
	}

	handle := c.file.GetProcHandle()
	pid, err := c.getProcessIDFromHandle(handle)
	if err != nil {
		return nil, err
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return nil, err
	}

	return proc.Wait()
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

	handle := c.file.GetProcHandle()
	pid, err := c.getProcessIDFromHandle(handle)

	return int(pid), err
}

func (c *consoleWindows) Kill() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	handle := c.file.GetProcHandle()
	pid, err := c.getProcessIDFromHandle(handle)
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}

	return proc.Kill()
}

func (c *consoleWindows) Signal(sig os.Signal) error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	handle := c.file.GetProcHandle()
	pid, err := c.getProcessIDFromHandle(handle)
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(int(pid))
	if err != nil {
		return err
	}

	return proc.Signal(sig)
}
