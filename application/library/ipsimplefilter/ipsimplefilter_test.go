package ipsimplefilter

import (
	"net"
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestContains(t *testing.T) {
	i := IPR{}
	i.Parse(`127.0.0.1-129.0.0.1`)
	ip := net.ParseIP(`127.0.0.2`)
	y := i.Contains(ip)
	test.True(t, y)
	ip = net.ParseIP(`128.0.0.2`)
	y = i.Contains(ip)
	test.True(t, y)
	ip = net.ParseIP(`129.0.0.2`)
	y = i.Contains(ip)
	test.False(t, y)
}
