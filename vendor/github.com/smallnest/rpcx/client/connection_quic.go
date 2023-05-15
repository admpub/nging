//go:build quic
// +build quic

package client

import (
	"crypto/tls"
	"net"

	"github.com/quic-go/quic-go"
	"github.com/smallnest/quick"
)

func newDirectQuicConn(c *Client, network, address string) (net.Conn, error) {
	tlsConf := c.option.TLSConfig
	if tlsConf == nil {
		tlsConf = &tls.Config{InsecureSkipVerify: true}
	}

	if len(tlsConf.NextProtos) == 0 {
		tlsConf.NextProtos = []string{"rpcx"}
	}

	quicConfig := &quic.Config{}

	return quick.Dial(address, tlsConf, quicConfig)
}
