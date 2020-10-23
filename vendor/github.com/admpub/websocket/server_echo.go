// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package websocket

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

// EchoUpgrader specifies parameters for upgrading an HTTP connection to a
// WebSocket connection.
type EchoUpgrader struct {
	// Handler receives a websocket connection after the handshake has been
	// completed. This must be provided.
	Handler func(*Conn) error

	// HandshakeTimeout specifies the duration for the handshake to complete.
	HandshakeTimeout time.Duration

	// ReadBufferSize and WriteBufferSize specify I/O buffer sizes. If a buffer
	// size is zero, then a default value of 4096 is used. The I/O buffer sizes
	// do not limit the size of the messages that can be sent or received.
	ReadBufferSize, WriteBufferSize int

	// Subprotocols specifies the server's supported protocols in order of
	// preference. If this field is set, then the Upgrade method negotiates a
	// subprotocol by selecting the first match in this list with a protocol
	// requested by the client.
	Subprotocols []string

	// Error specifies the function for generating HTTP error responses. If Error
	// is nil, then http.Error is used to generate the HTTP response.
	Error func(ctx echo.Context, status int, reason error)

	// CheckOrigin returns true if the request Origin header is acceptable. If
	// CheckOrigin is nil, the host in the Origin header must not be set or
	// must match the host of the request.
	CheckOrigin func(r engine.Request) bool

	// EnableCompression specify if the server should attempt to negotiate per
	// message compression (RFC 7692). Setting this value to true does not
	// guarantee that compression will be supported. Currently only "no context
	// takeover" modes are supported.
	EnableCompression bool
}

func (u *EchoUpgrader) returnError(ctx echo.Context, status int, reason string) error {
	err := echo.NewHTTPError(status, reason)
	ctx.Response().Header().Set("Sec-Websocket-Version", "13")
	if u.Error != nil {
		u.Error(ctx, status, err)
	}
	return err
}

// checkSameOrigin returns true if the origin is not set or is equal to the request host.
func echoCheckSameOrigin(r engine.Request) bool {
	origin := r.Header().Get("Origin")
	if len(origin) == 0 {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	return u.Host == r.Host()
}

func (u *EchoUpgrader) selectSubprotocol(r engine.Request, responseHeader http.Header) string {
	if u.Subprotocols != nil {
		clientProtocols := EchoSubprotocols(r)
		for _, serverProtocol := range u.Subprotocols {
			for _, clientProtocol := range clientProtocols {
				if clientProtocol == serverProtocol {
					return clientProtocol
				}
			}
		}
	} else if responseHeader != nil {
		return responseHeader.Get("Sec-Websocket-Protocol")
	}
	return ""
}

// Upgrade upgrades the HTTP server connection to the WebSocket protocol.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie) and the
// application negotiated subprotocol (Sec-Websocket-Protocol).
//
// If the upgrade fails, then Upgrade replies to the client with an HTTP error
// response.
func (u *EchoUpgrader) Upgrade(ctx echo.Context, handler func(*Conn) error, responseHeader http.Header) error {
	r := ctx.Request()
	w := ctx.Response()
	if r.Method() != "GET" {
		return u.returnError(ctx, http.StatusMethodNotAllowed, "websocket: method not GET")
	}

	if _, ok := responseHeader["Sec-Websocket-Extensions"]; ok {
		return u.returnError(ctx, http.StatusInternalServerError, "websocket: application specific Sec-Websocket-Extensions headers are unsupported")
	}

	reqHeader := r.Header().Std()

	if !tokenListContainsValue(reqHeader, "Sec-Websocket-Version", "13") {
		return u.returnError(ctx, http.StatusBadRequest, "websocket: version != 13")
	}

	if !tokenListContainsValue(reqHeader, "Connection", "upgrade") {
		return u.returnError(ctx, http.StatusBadRequest, "websocket: could not find connection header with token 'upgrade'")
	}

	if !tokenListContainsValue(reqHeader, "Upgrade", "websocket") {
		return u.returnError(ctx, http.StatusBadRequest, "websocket: could not find upgrade header with token 'websocket'")
	}

	checkOrigin := u.CheckOrigin
	if checkOrigin == nil {
		checkOrigin = echoCheckSameOrigin
	}
	if !checkOrigin(r) {
		return u.returnError(ctx, http.StatusForbidden, "websocket: origin not allowed")
	}

	challengeKey := r.Header().Get("Sec-Websocket-Key")
	if challengeKey == "" {
		return u.returnError(ctx, http.StatusBadRequest, "websocket: key missing or blank")
	}

	subprotocol := u.selectSubprotocol(r, responseHeader)

	// Negotiate PMCE
	var compress bool
	if u.EnableCompression {
		for _, ext := range parseExtensions(reqHeader) {
			if ext[""] != "permessage-deflate" {
				continue
			}
			compress = true
			break
		}
	}
	var err error

	if handler == nil {
		handler = u.Handler
	}

	writerHeader := w.Header()
	writerHeader.Set("Upgrade", "websocket")
	writerHeader.Set("Connection", "Upgrade")
	writerHeader.Set("Sec-WebSocket-Accept", computeAcceptKey(challengeKey))

	if subprotocol == "" {
		// Find the best protocol, if any
		clientProtocols := EchoSubprotocols(r)
		if len(clientProtocols) != 0 {
			subprotocol = matchSubprotocol(clientProtocols, u.Subprotocols)
			if subprotocol != "" {
				writerHeader.Set("Sec-Websocket-Protocol", subprotocol)
			}
		}
	}

	if compress {
		writerHeader.Set("Sec-Websocket-Extensions", "permessage-deflate; server_no_context_takeover; client_no_context_takeover")
	}
	for k, vs := range responseHeader {
		if k == "Sec-Websocket-Protocol" {
			continue
		}
		writerHeader.Set(k, strings.Join(vs, "; "))
	}
	w.WriteHeader(http.StatusSwitchingProtocols)

	err = w.Hijacker(func(netConn net.Conn) {
		c := newConn(netConn, true, u.ReadBufferSize, u.WriteBufferSize)
		c.subprotocol = subprotocol
		if compress {
			c.newCompressionWriter = compressNoContextTakeover
			c.newDecompressionReader = decompressNoContextTakeover
		}

		// Clear deadlines set by HTTP server.
		netConn.SetDeadline(time.Time{})

		if u.HandshakeTimeout > 0 {
			netConn.SetWriteDeadline(time.Now().Add(u.HandshakeTimeout))
		}
		if u.HandshakeTimeout > 0 {
			netConn.SetWriteDeadline(time.Time{})
		}
		if handler != nil {
			err = handler(c)
		}
	})
	return err
}

// EchoUpgrade upgrades the HTTP server connection to the WebSocket protocol.
//
// This function is deprecated, use websocket.Upgrader instead.
//
// The application is responsible for checking the request origin before
// calling Upgrade. An example implementation of the same origin policy is:
//
//	if req.Header.Get("Origin") != "http://"+req.Host {
//		http.Error(w, "Origin not allowed", 403)
//		return
//	}
//
// If the endpoint supports subprotocols, then the application is responsible
// for negotiating the protocol used on the connection. Use the Subprotocols()
// function to get the subprotocols requested by the client. Use the
// Sec-Websocket-Protocol response header to specify the subprotocol selected
// by the application.
//
// The responseHeader is included in the response to the client's upgrade
// request. Use the responseHeader to specify cookies (Set-Cookie) and the
// negotiated subprotocol (Sec-Websocket-Protocol).
//
// The connection buffers IO to the underlying network connection. The
// readBufSize and writeBufSize parameters specify the size of the buffers to
// use. Messages can be larger than the buffers.
//
// If the request is not a valid WebSocket handshake, then Upgrade returns an
// error of type HandshakeError. Applications should handle this error by
// replying to the client with an HTTP error response.
func EchoUpgrade(ctx echo.Context, handler func(*Conn) error, responseHeader http.Header, readBufSize, writeBufSize int) error {
	u := EchoUpgrader{ReadBufferSize: readBufSize, WriteBufferSize: writeBufSize}
	u.Error = func(ctx echo.Context, status int, reason error) {
		// don't return errors to maintain backwards compatibility
	}
	u.CheckOrigin = func(r engine.Request) bool {
		// allow all connections by default
		return true
	}
	return u.Upgrade(ctx, handler, responseHeader)
}

// EchoSubprotocols returns the subprotocols requested by the client in the
// Sec-Websocket-Protocol header.
func EchoSubprotocols(r engine.Request) []string {
	h := strings.TrimSpace(r.Header().Get("Sec-Websocket-Protocol"))
	if h == "" {
		return nil
	}
	protocols := strings.Split(h, ",")
	for i := range protocols {
		protocols[i] = strings.TrimSpace(protocols[i])
	}
	return protocols
}

// EchoIsWebSocketUpgrade returns true if the client requested upgrade to the
// WebSocket protocol.
func EchoIsWebSocketUpgrade(r engine.Request) bool {
	reqHeader := r.Header().Std()
	return tokenListContainsValue(reqHeader, "Connection", "upgrade") &&
		tokenListContainsValue(reqHeader, "Upgrade", "websocket")
}

func matchSubprotocol(clientProtocols, serverProtocols []string) string {
	for _, serverProtocol := range serverProtocols {
		for _, clientProtocol := range clientProtocols {
			if clientProtocol == serverProtocol {
				return clientProtocol
			}
		}
	}

	return ""
}
