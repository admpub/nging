package collector

import "io"

type Context struct {
	Pages  []*PageConfig
	Values map[string]string
}

func (c *Context) Set(name string, value string) *Context {
	c.Values[name] = value
	return c
}

func (c *Context) Reader(conf *PageConfig) (r io.Reader) {
	return
}
