package standard

import (
	"net/url"
)

type URL struct {
	url   *url.URL
	query url.Values
}

func (u *URL) SetPath(path string) {
	u.url.Path = path
}

func (u *URL) RawPath() string {
	return u.url.EscapedPath()
}

func (u *URL) Path() string {
	return u.url.Path
}

func (u *URL) QueryValue(name string) string {
	if u.query == nil {
		u.query = u.url.Query()
	}
	return u.query.Get(name)
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
		u.query = u.url.Query()
	}
	return u.query
}

func (u *URL) reset(url *url.URL) {
	u.url = url
	u.query = nil
}

func (u *URL) RawQuery() string {
	return u.url.RawQuery
}

func (u *URL) SetRawQuery(rawQuery string) {
	u.url.RawQuery = rawQuery
}

func (u *URL) String() string {
	return u.url.String()
}

func (u *URL) Object() interface{} {
	return u.url
}
