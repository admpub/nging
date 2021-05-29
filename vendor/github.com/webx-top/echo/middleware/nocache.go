package middleware

// Ported from Goji's middleware, source:
// https://github.com/zenazn/goji/tree/master/web/middleware

import (
	"time"

	"github.com/webx-top/echo"
)

// Unix epoch time
var epoch = time.Unix(0, 0).Format(time.RFC1123)

// Taken from https://github.com/mytrile/nocache
var noCacheHeaders = map[string]string{
	"Expires":         epoch,
	"Cache-Control":   "no-cache, private, max-age=0",
	"Pragma":          "no-cache",
	"X-Accel-Expires": "0",
}

var etagHeaders = []string{
	"ETag",
	"If-Modified-Since",
	"If-Match",
	"If-None-Match",
	"If-Range",
	"If-Unmodified-Since",
}

// NoCache is a simple piece of middleware that sets a number of HTTP headers to prevent
// a router (or subrouter) from being cached by an upstream proxy and/or client.
//
// As per http://wiki.nginx.org/HttpProxyModule - NoCache sets:
//      Expires: Thu, 01 Jan 1970 00:00:00 UTC
//      Cache-Control: no-cache, private, max-age=0
//      X-Accel-Expires: 0
//      Pragma: no-cache (for HTTP/1.0 proxies/clients)
func NoCache(skippers ...echo.Skipper) echo.MiddlewareFuncd {
	var skipper echo.Skipper
	if len(skippers) > 0 {
		skipper = skippers[0]
	}
	if skipper == nil {
		skipper = echo.DefaultSkipper
	}
	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next.Handle(c)
			}
			SetNoCacheHeader(c)
			return next.Handle(c)
		}
	}
}

func SetNoCacheHeader(c echo.Context) {
	reqHdr := c.Request().Header()
	resHdr := c.Response().Header()

	// Delete any ETag headers that may have been set
	for _, v := range etagHeaders {
		if len(reqHdr.Get(v)) > 0 {
			reqHdr.Del(v)
		}
	}

	// Set our NoCache headers
	for k, v := range noCacheHeaders {
		resHdr.Set(k, v)
	}
}
