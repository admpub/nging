// +build !windows

package server

import "github.com/shirou/gopsutil/net"

func NetStatTCP() (<-chan net.ConnectionStat, error) {
	return nil, ErrNotImplemented
}

func NetStatUDP() (<-chan net.ConnectionStat, error) {
	return nil, ErrNotImplemented
}
