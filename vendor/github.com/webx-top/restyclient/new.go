package restyclient

import (
	"net/http"
	"net/http/cookiejar"

	"github.com/admpub/log"
	"github.com/admpub/resty/v2"
	"github.com/webx-top/com"
	"github.com/webx-top/com/httpClientOptions"
	"golang.org/x/net/publicsuffix"
)

func NewCookiejar() (*cookiejar.Jar, error) {
	return cookiejar.New(&cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	})
}

func ProxyURL(proxy string) func(c *http.Client) {
	return func(c *http.Client) {
		if err := SetProxy(c, proxy); err != nil {
			log.Error(err)
		}
	}
}

func NewWithOptions(options ...com.HTTPClientOptions) *resty.Client {
	hclient := com.HTTPClientWithTimeout(
		DefaultTimeout,
		options...,
	)
	return resty.NewWithClient(hclient)
}

func New(proxy ...string) *resty.Client {
	cookieJar, _ := NewCookiejar()
	options := []com.HTTPClientOptions{
		httpClientOptions.Timeout(DefaultTimeout),
		httpClientOptions.CookieJar(cookieJar),
	}
	if len(proxy) > 0 && len(proxy[0]) > 0 {
		options = append(options, ProxyURL(proxy[0]))
	}
	return NewWithOptions(options...)
}
