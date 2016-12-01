// +build !appengine

package fasthttp

import (
	"net/url"

	"github.com/admpub/fasthttp"
	"github.com/webx-top/echo/engine"
)

type UrlValue struct {
	*fasthttp.Args
	initFn func() *fasthttp.Args
	values *url.Values
}

func (u *UrlValue) Add(key string, value string) {
	u.All()
	u.values.Add(key, value)
}

func (u *UrlValue) Del(key string) {
	u.init()
	u.Args.Del(key)
}

func (u *UrlValue) Get(key string) string {
	u.init()
	return engine.Bytes2str(u.Args.Peek(key))
}

func (u *UrlValue) Gets(key string) []string {
	u.init()
	u.All()
	if v, ok := (*u.values)[key]; ok {
		return v
	}
	return []string{}
}

func (u *UrlValue) Set(key string, value string) {
	u.init()
	u.Args.Set(key, value)
}

func (u *UrlValue) Reset(data url.Values) {
	a := &fasthttp.Args{}
	for key, values := range data {
		for _, value := range values {
			a.Set(key, value)
		}
	}
	a.CopyTo(u.Args)
	u.values = &data
}

func (u *UrlValue) init() {
	if u.Args != nil {
		return
	}
	u.Args = u.initFn()
}

func (u *UrlValue) Encode() string {
	u.init()
	return u.Args.String()
}

func (u *UrlValue) All() map[string][]string {
	if u.values != nil {
		return *u.values
	}
	r := url.Values{}
	u.init()
	u.Args.VisitAll(func(k, v []byte) {
		key := engine.Bytes2str(k)
		r.Add(key, engine.Bytes2str(v))
	})
	u.values = &r
	return *u.values
}

func NewValue(c *Request) *Value {
	v := &Value{
		queryArgs: &UrlValue{initFn: func() *fasthttp.Args {
			return c.context.QueryArgs()
		}},
		postArgs: &UrlValue{initFn: func() *fasthttp.Args {
			return c.context.PostArgs()
		}},
		request: c,
	}
	return v
}

type Value struct {
	queryArgs *UrlValue
	postArgs  *UrlValue
	form      *url.Values
	request   *Request
}

func (v *Value) Add(key string, value string) {
	v.init()
	v.form.Add(key, value)
}

func (v *Value) Del(key string) {
	v.init()
	v.form.Del(key)
}

func (v *Value) Get(key string) string {
	v.init()
	return v.form.Get(key)
}

func (v *Value) Gets(key string) []string {
	v.init()
	form := *v.form
	if v, ok := form[key]; ok {
		return v
	}
	return []string{}
}

func (v *Value) Set(key string, value string) {
	v.init()
	v.form.Set(key, value)
}

func (v *Value) Encode() string {
	v.init()
	return v.form.Encode()
}

func (v *Value) init() {
	if v.form != nil {
		return
	}
	form := url.Values(v.queryArgs.All())
	for key, vals := range v.postArgs.All() {
		form[key] = vals
	}
	mf := v.request.MultipartForm()
	if mf != nil && mf.Value != nil {
		for key, vals := range mf.Value {
			form[key] = vals
		}
	}
	v.form = &form
}

func (v *Value) All() map[string][]string {
	v.init()
	return *v.form
}

func (v *Value) Reset(data url.Values) {
	v.form = &data
}
