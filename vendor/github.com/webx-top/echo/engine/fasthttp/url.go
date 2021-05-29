// +build !appengine

package fasthttp

import (
	"net/url"

	"github.com/admpub/fasthttp"
	"github.com/webx-top/echo/engine"
)

type URL struct {
	url   *fasthttp.URI
	query url.Values
}

func (u *URL) SetPath(path string) {
	u.url.SetPath(path)
}

func (u *URL) RawPath() string {
	return engine.Bytes2str(u.url.PathOriginal())
}

func (u *URL) Path() string {
	return engine.Bytes2str(u.url.Path())
}

func (u *URL) QueryValue(name string) string {
	return engine.Bytes2str(u.url.QueryArgs().Peek(name))
}

func (u *URL) QueryValues(name string) []string {
	u.Query()
	if v, ok := u.query[name]; ok {
		return v
	}
	return []string{}
}

func (u *URL) Query() url.Values {
	if u.query == nil {
		u.query = url.Values{}
		u.url.QueryArgs().VisitAll(func(key []byte, value []byte) {
			u.query.Set(string(key), string(value))
		})
	}
	return u.query
}

func (u *URL) RawQuery() string {
	return engine.Bytes2str(u.url.QueryString())
}

func (u *URL) SetRawQuery(rawQuery string) {
	u.url.SetQueryString(rawQuery)
}

func (u *URL) String() string {
	return u.url.String()
}

func (u *URL) Object() interface{} {
	return u.url
}

func (u *URL) reset(url *fasthttp.URI) {
	u.url = url
}
