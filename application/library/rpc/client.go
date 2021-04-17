package rpc

import (
	"context"
	"crypto/tls"
	"sync"

	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/share"
)

func NewClient(addr string, token string, tls *tls.Config) *Client {
	return &Client{
		addr:    addr,
		token:   token,
		tls:     tls,
		clients: sync.Map{},
	}
}

type Client struct {
	token   string
	addr    string
	tls     *tls.Config
	d       client.ServiceDiscovery
	once    sync.Once
	clients sync.Map
}

func (r *Client) connect() (err error) {
	r.d, err = client.NewPeer2PeerDiscovery("tcp@"+r.addr, "")
	return
}

func (r *Client) Client(servicePath string) client.XClient {
	r.once.Do(func() {
		if err := r.connect(); err != nil {
			panic(err)
		}
	})
	options := client.DefaultOption
	options.TLSConfig = r.tls
	return client.NewXClient(servicePath, client.Failtry, client.RandomSelect, r.d, options)
}

func (r *Client) Call(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}) error {
	c := r.Client(servicePath)
	if len(r.token) > 0 {
		c.Auth("bearer " + r.token)
		ctx = context.WithValue(ctx, share.ReqMetaDataKey, make(map[string]string))
	}
	defer c.Close()
	return c.Call(ctx, serviceMethod, args, reply)
}

func (r *Client) Close() error {
	if r.d != nil {
		r.d.Close()
	}
	return nil
}
