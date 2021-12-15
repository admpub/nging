// +build go1.11

package engine

import (
	"net"

	reuseport "github.com/admpub/go-reuseport"
)

func newListener(scheme, address string, reuse bool) (l net.Listener, err error) {
	if reuse {
		l, err = reuseport.Listen(scheme, address)
	} else {
		l, err = net.Listen(scheme, address)
	}
	return
}
