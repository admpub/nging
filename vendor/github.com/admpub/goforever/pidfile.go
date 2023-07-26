package goforever

import (
	"os"
	"strconv"
)

type Pidfile string

// Read the pidfile.
func (f *Pidfile) Read() int {
	data, err := os.ReadFile(string(*f))
	if err != nil {
		return 0
	}
	pid, err := strconv.ParseInt(string(data), 0, 32)
	if err != nil {
		return 0
	}
	return int(pid)
}

// Write the pidfile.
func (f *Pidfile) Write(data int) error {
	return os.WriteFile(string(*f), []byte(strconv.Itoa(data)), 0660)
}

// Delete the pidfile
func (f *Pidfile) Delete() bool {
	_, err := os.Stat(string(*f))
	if err != nil {
		return true
	}
	err = os.Remove(string(*f))
	return err == nil
}
