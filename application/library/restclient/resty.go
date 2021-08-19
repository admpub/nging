package restclient

import (
	"context"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"github.com/admpub/resty/v2"
	"github.com/webx-top/com"
	"github.com/webx-top/com/httpClientOptions"
	"golang.org/x/net/proxy"
	"golang.org/x/net/publicsuffix"
)

var (
	// DefaultTimeout 默认超时时间
	DefaultTimeout     = 10 * time.Second
	restyClient        *resty.Client
	restyRetryable     *resty.Client
	restyOnce          sync.Once
	restyRetryableOnce sync.Once
)

func initRestyClient() {
	restyClient = resty.New().SetTimeout(DefaultTimeout)
}

func ResetResty() {
	restyOnce = sync.Once{}
}

func Resty() *resty.Request {
	restyOnce.Do(initRestyClient)
	return restyClient.R()
}

// - retryable -

func initRetryable() {
	restyRetryable = resty.New().SetRetryCount(3).SetTimeout(DefaultTimeout)
}

func ResetRestyRetryable() {
	restyRetryableOnce = sync.Once{}
}

func RestyRetryable() *resty.Request {
	restyRetryableOnce.Do(initRetryable)
	return restyRetryable.R()
}

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

// SetProxy Proxy client
func SetProxy(client *http.Client, proxyString string) error {
	proxyURL, err := url.Parse(proxyString)
	if err != nil {
		return err
	}
	if client.Transport == nil {
		client.Transport = &http.Transport{}
	}
	switch proxyURL.Scheme {
	case "http", "https":
		client.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
	case "socks5", "socks5h":
		fallthrough
	default: // proxy.RegisterDialerType(``)
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return err
		}
		if dial, ok := dialer.(proxy.ContextDialer); ok {
			client.Transport.(*http.Transport).DialContext = dial.DialContext
		} else {
			client.Transport.(*http.Transport).DialContext = func(ctx context.Context, network string, address string) (net.Conn, error) {
				return dialer.Dial(network, address)
			}
		}
	}
	return nil
}
