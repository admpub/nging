package restyclient

import (
	"context"
	"net"
	"net/http"
	"net/url"

	"golang.org/x/net/proxy"
)

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
