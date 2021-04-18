package proxy

import (
	"net/url"

	"github.com/webx-top/echo"
	mw "github.com/webx-top/echo/middleware"
)

func NewTarget(t *mw.ProxyTarget) mw.ProxyTargeter {
	return &Target{ProxyTarget: t}
}

type Target struct {
	*mw.ProxyTarget
}

func (t *Target) GetURL(c echo.Context) *url.URL {
	address := c.Internal().String(`frp.admin.address`)
	u := *t.ProxyTarget.URL
	u.Host = address
	//echo.Dump(u)
	return &u
}
