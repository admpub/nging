package goforever

import (
	"os"
)

var Default = New()

func New() *Process {
	return NewProcess("goforever", "")
}

func NewProcess(name string, command string, args ...string) *Process {
	if len(command) == 0 {
		command = os.Args[0]
	}
	p := &Process{
		Name:     name,
		Args:     args,
		Command:  command,
		Respawn:  1,
		Children: make(map[string]*Process, 0),
		Pidfile:  Pidfile(name + `.pid`),
	}
	return p
}

func StartChild(name string) (*Process, error) {
	return Default.StartChild(name)
}

func RestartChild(name string) (*Process, error) {
	return Default.RestartChild(name)
}

func StopChild(name string) error {
	return Default.StopChild(name)
}

func Child(name string) *Process {
	return Default.Children.Get(name)
}

func ChildKeys() []string {
	return Default.Children.Keys()
}

func Add(name string, procs *Process, run ...bool) *Process {
	return Default.Add(name, procs, run...)
}

func Run() {
	Default.Run()
}
