package rpc

import (
	"context"
	"crypto/tls"
	"reflect"

	"github.com/admpub/log"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"github.com/webx-top/com"
)

func NewServer(addr string, token string, tls *tls.Config) *Server {
	return &Server{
		addr:              addr,
		token:             token,
		tls:               tls,
		namedServices:     make(map[string][]interface{}),
		namedFuncServices: make(map[string]map[string]interface{}),
	}
}

type Server struct {
	token             string
	addr              string
	tls               *tls.Config
	s                 *server.Server
	services          []interface{}
	namedServices     map[string][]interface{}
	namedFuncServices map[string]map[string]interface{}
}

func (r *Server) Register(services ...interface{}) *Server {
	r.services = append(r.services, services...)
	return r
}

func (r *Server) RegisterName(name string, service interface{}) *Server {
	if _, ok := r.namedServices[name]; !ok {
		r.namedServices[name] = []interface{}{}
	}
	r.namedServices[name] = append(r.namedServices[name], service)
	return r
}

func (r *Server) RegisterFuncName(servicePath, name string, service interface{}) *Server {
	if _, ok := r.namedFuncServices[servicePath]; !ok {
		r.namedFuncServices[servicePath] = map[string]interface{}{}
	}
	r.namedFuncServices[servicePath][name] = service
	return r
}

func (r *Server) Start() error {
	r.s = server.NewServer(server.WithTLSConfig(r.tls))
	if len(r.token) > 0 {
		r.s.AuthFunc = r.auth
	}
	for _, s := range r.services {
		r.s.Register(s, ``)
	}
	for k, ss := range r.namedServices {
		for _, s := range ss {
			t := reflect.TypeOf(s)
			if t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			if t.Kind() == reflect.Func {
				// 自动获取函数名称（最终获取到的名称会比真实的函数名称多一个“-fm”后缀，这点在调用的时候要注意，别忘了加上此后缀）
				log.Debug(`[rpc] register function: `, com.FuncName(s))
				r.s.RegisterFunction(k, s, ``)
				continue
			}
			log.Debug(`[rpc] register name: `, com.FuncName(s), t.Kind)
			r.s.RegisterName(k, s, ``)
		}
	}
	for servicePath, ss := range r.namedFuncServices {
		for name, s := range ss {
			r.s.RegisterFunctionName(servicePath, name, s, ``)
		}
	}
	return r.s.Serve("tcp", r.addr)
}

func (r *Server) auth(ctx context.Context, req *protocol.Message, token string) error {
	if token == "bearer "+r.token {
		return nil
	}

	return ErrInvalidToken
}

func (r *Server) Close() error {
	if r.s != nil {
		return r.s.Close()
	}
	return nil
}
