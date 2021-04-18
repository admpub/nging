package net

import (
	"net"
	"net/url"
	"time"

	"golang.org/x/net/websocket"
)

// addr: domain:port
func ConnectWSSServer(addr string) (net.Conn, error) {
	addr = "wss://" + addr + FrpWebsocketPath
	uri, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}

	origin := "https://" + uri.Host
	cfg, err := websocket.NewConfig(addr, origin)
	if err != nil {
		return nil, err
	}
	cfg.Dialer = &net.Dialer{
		Timeout: 10 * time.Second,
	}

	conn, err := websocket.DialConfig(cfg)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
