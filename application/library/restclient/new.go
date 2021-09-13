package restclient

import (
	"net/http/cookiejar"

	"github.com/admpub/resty/v2"
	"github.com/webx-top/com"
	"github.com/webx-top/com/httpClientOptions"
	"golang.org/x/net/publicsuffix"
)

func NewClient(proxy ...string) *resty.Client {
	cookieJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	hclient := com.HTTPClientWithTimeout(
		DefaultTimeout,
		httpClientOptions.Timeout(DefaultTimeout),
		httpClientOptions.CookieJar(cookieJar),
	)
	if len(proxy) > 0 && len(proxy[0]) > 0 {
		if err := SetProxy(hclient, proxy[0]); err != nil {
			panic(err)
		}
	}
	client := resty.NewWithClient(hclient)
	return client
}
