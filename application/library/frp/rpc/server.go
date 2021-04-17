package rpc

import (
	"context"
	"crypto/tls"
	"errors"

	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
)

func NewServer(addr string, token string, tls *tls.Config) *Server {
	return &Server{addr: addr, token: token, tls: tls}
}

type Server struct {
	token    string
	addr     string
	tls      *tls.Config
	s        *server.Server
	services []interface{}
}

func (r *Server) Register(services ...interface{}) *Server {
	r.services = append(r.services, services...)
	return r
}

func (r *Server) Start() error {
	r.s = server.NewServer(server.WithTLSConfig(r.tls))
	r.s.AuthFunc = r.auth
	//s.RegisterName("Arith", new(example.Arith), "")
	for _, s := range r.services {
		r.s.Register(s, ``)
	}
	return r.s.Serve("tcp", r.addr)
}

func (r *Server) auth(ctx context.Context, req *protocol.Message, token string) error {
	if token == "bearer "+r.token {
		return nil
	}

	return errors.New("invalid token")
}
