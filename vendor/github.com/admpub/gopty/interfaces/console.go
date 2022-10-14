package interfaces

import (
	"io"
	"os"
)

// Console communication interface
type Console interface {
	io.Reader
	io.Writer
	io.Closer

	// SetSize sets the console size
	SetSize(cols int, rows int) error

	// GetSize gets the console size
	// cols, rows, error
	GetSize() (int, int, error)

	// Start starts the process with the supplied args
	Start(args []string) error

	// Wait waits the process to finish
	Wait() (*os.ProcessState, error)

	// SetCWD sets the current working dir of the process
	SetCWD(cwd string) error

	// SetENV sets environment variables to pass to the child process
	SetENV(environ []string) error

	// Pid returns the pid of the running process
	Pid() (int, error)

	// Kill kills the process. See exec/Process.Kill
	Kill() error

	// Signal sends a signal to the process. See exec/Process.Signal
	Signal(sig os.Signal) error
}
