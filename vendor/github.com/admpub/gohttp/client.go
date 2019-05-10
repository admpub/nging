package gohttp

import (
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"

	"golang.org/x/net/publicsuffix"
)

type Option struct {
	Address        []string
	ConnectTimeout time.Duration
	TLSTimeout     time.Duration
	Timeout        time.Duration
	Agent          string
	Delay          time.Duration
	MaxRedirects   int
	MaxIdleConns   int
}

type clientResource struct {
	Transport http.RoundTripper
	Jar       http.CookieJar
}

type useInfo struct {
	Index    int
	LastTime time.Time
}

var defaultOption = &Option{
	ConnectTimeout: 30000 * time.Millisecond,
	TLSTimeout:     15 * time.Second,
	Agent:          "gohttp v1.0",
	Address:        make([]string, 0),
	MaxRedirects:   -1,
	MaxIdleConns:   0,
}

var debug = false
var defaultDialer = &net.Dialer{Timeout: defaultOption.ConnectTimeout}
var defaultTransport = MakeTransport("0.0.0.0")
var defaultCookiejar = MakeCookiejar()
var proxyTransport *http.Transport

var hostDelay = make(map[string]time.Duration)
var hostDelayLock sync.RWMutex

var defaultGetter = NewIpRollClient(defaultOption.Address...)

func MakeCookiejar() http.CookieJar {
	cookiejarOptions := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, _ := cookiejar.New(&cookiejarOptions)

	return jar
}

func MakeClient(transport http.RoundTripper, jar http.CookieJar) *http.Client {
	return &http.Client{Jar: jar, Transport: transport, Timeout: 60 * time.Second}
}

func MakeTransport(ip string) *http.Transport {
	addr, _ := net.ResolveTCPAddr("tcp", ip+":0")
	dialer := &net.Dialer{
		Timeout:   defaultOption.ConnectTimeout,
		LocalAddr: addr,
	}
	transport := &http.Transport{
		Dial:                dialer.Dial,
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: defaultOption.MaxIdleConns,
		TLSHandshakeTimeout: defaultOption.TLSTimeout,
	}

	if defaultOption.MaxIdleConns <= 0 {
		transport.DisableKeepAlives = true
	}

	return transport
}

func SetDebug(d bool) {
	defer hostDelayLock.Unlock()
	hostDelayLock.Lock()

	debug = d
}

func IsDebug() bool {
	defer hostDelayLock.RUnlock()
	hostDelayLock.RLock()
	return debug
}

func SetHostDelay(host string, delay time.Duration) {
	defer hostDelayLock.Unlock()
	hostDelayLock.Lock()
	if d, ok := hostDelay[host]; ok && delay > d {
		hostDelay[host] = delay
		return
	}
	hostDelay[host] = delay
}

func GetHostDelay(host string) time.Duration {
	defer hostDelayLock.RUnlock()
	hostDelayLock.RLock()

	if d, ok := hostDelay[host]; ok {
		return d
	}

	return defaultOption.Delay
}

func SetOption(option *Option) {
	if option.Agent != "" {
		defaultOption.Agent = option.Agent
	}

	if option.ConnectTimeout > 0 {
		defaultOption.ConnectTimeout = option.ConnectTimeout
	}

	if option.TLSTimeout > 0 {
		defaultOption.TLSTimeout = option.TLSTimeout
	}

	if option.Delay > 0 {
		defaultOption.Delay = option.Delay
	}

	if option.Address != nil && len(option.Address) > 0 {
		defaultOption.Address = make([]string, 0)
		defaultOption.Address = append(defaultOption.Address, option.Address...)
		defaultGetter = NewIpRollClient(defaultOption.Address...)
	}

	if option.MaxRedirects > 0 {
		defaultOption.MaxRedirects = option.MaxRedirects
	}

	if option.MaxIdleConns > 0 {
		defaultOption.MaxIdleConns = option.MaxIdleConns
		defaultTransport.MaxIdleConnsPerHost = option.MaxIdleConns
	}
}

func ResetCookie(urlstr string) error {
	uri, err := url.Parse(urlstr)
	if err != nil {
		return err
	}
	cookies := defaultCookiejar.Cookies(uri)
	for _, c := range cookies {
		c.Expires = time.Now().Add(-1 * time.Hour)
	}
	defaultCookiejar.SetCookies(uri, cookies)

	defaultGetter.ResetCookie(uri)

	return nil
}

func GetDefaultDialer() *net.Dialer {
	return defaultDialer
}

func GetDefaultTransport() *http.Transport {
	return defaultTransport
}

func GetDefaultClient() *http.Client {
	return MakeClient(defaultTransport, defaultCookiejar)
}

func GetDefaultGetter() ClientGetter {
	return defaultGetter
}
