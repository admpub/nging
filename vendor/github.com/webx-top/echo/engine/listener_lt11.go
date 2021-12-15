// +build !go1.11

package engine

import (
	"net"
)

func newListener(scheme, address string, reuse bool) (l net.Listener, err error) {
	l, err = net.Listen(scheme, address)
	return
}
