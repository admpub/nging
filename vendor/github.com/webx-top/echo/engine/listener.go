package engine

import (
	"net"
	"strings"
)

func NewListener(address string, reuse bool) (net.Listener, error) {
	scheme := "tcp"
	delim := "://"
	if pos := strings.Index(address, delim); pos > 0 {
		scheme = address[0:pos]
		address = address[pos+len(delim):]
	}
	return newListener(scheme, address, reuse)
}
