package httpClientOptions

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/webx-top/com"
)

func InsecureSkipVerify(skips ...bool) com.HTTPClientOptions {
	skip := true
	if len(skips) > 0 {
		skip = skips[0]
	}
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		if tr.TLSClientConfig == nil {
			tr.TLSClientConfig = &tls.Config{}
		}
		tr.TLSClientConfig.InsecureSkipVerify = skip
	}
}

func Proxy(fn func(*http.Request) (*url.URL, error)) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.Proxy = fn
	}
}

func DialContext(fn func(ctx context.Context, network, addr string) (net.Conn, error)) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.DialContext = fn
	}
}

func Dial(fn func(network, addr string) (net.Conn, error)) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.Dial = fn
	}
}

func DialTLSContext(fn func(ctx context.Context, network, addr string) (net.Conn, error)) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.DialTLSContext = fn
	}
}

func DialTLS(fn func(network, addr string) (net.Conn, error)) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.Dial = fn
	}
}

func TLSClientConfig(cfg *tls.Config) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.TLSClientConfig = cfg
	}
}

func TLSHandshakeTimeout(t time.Duration) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.TLSHandshakeTimeout = t
	}
}

func DisableKeepAlives(disabled bool) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.DisableKeepAlives = disabled
	}
}

func DisableCompression(disabled bool) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.DisableCompression = disabled
	}
}

func MaxIdleConns(max int) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.MaxIdleConns = max
	}
}

func MaxIdleConnsPerHost(max int) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.MaxIdleConnsPerHost = max
	}
}

func MaxConnsPerHost(max int) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.MaxConnsPerHost = max
	}
}

func IdleConnTimeout(t time.Duration) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.IdleConnTimeout = t
	}
}

func ResponseHeaderTimeout(t time.Duration) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.ResponseHeaderTimeout = t
	}
}

func ExpectContinueTimeout(t time.Duration) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.ExpectContinueTimeout = t
	}
}

func TLSNextProto(next map[string]func(authority string, c *tls.Conn) http.RoundTripper) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.TLSNextProto = next
	}
}

func ProxyConnectHeader(header http.Header) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.ProxyConnectHeader = header
	}
}

func MaxResponseHeaderBytes(max int64) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.MaxResponseHeaderBytes = max
	}
}

func WriteBufferSize(size int) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.WriteBufferSize = size
	}
}

func ReadBufferSize(size int) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.ReadBufferSize = size
	}
}

func ForceAttemptHTTP2(on bool) com.HTTPClientOptions {
	return func(c *http.Client) {
		if c.Transport == nil {
			c.Transport = &http.Transport{}
		}
		tr := c.Transport.(*http.Transport)
		tr.ForceAttemptHTTP2 = on
	}
}

func Transport(tr *http.Transport) com.HTTPClientOptions {
	return func(c *http.Client) {
		c.Transport = tr
	}
}

func CheckRedirect(fn func(req *http.Request, via []*http.Request) error) com.HTTPClientOptions {
	return func(c *http.Client) {
		c.CheckRedirect = fn
	}
}

func CookieJar(cookieJar http.CookieJar) com.HTTPClientOptions {
	return func(c *http.Client) {
		c.Jar = cookieJar
	}
}

func Timeout(timeout time.Duration) com.HTTPClientOptions {
	return func(c *http.Client) {
		c.Timeout = timeout
	}
}
