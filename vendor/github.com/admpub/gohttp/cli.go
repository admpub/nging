package gohttp

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type ClientGetter interface {
	GetHttpClient(httpurl string, proxyurl string, usejar bool) (*http.Client, error)
}

type IpRollClient struct {
	ips        []string
	useLock    sync.RWMutex
	useMap     map[string]*useInfo
	clientMap  map[string]*clientResource
	clientLock sync.RWMutex
}

func NewIpRollClient(ip ...string) *IpRollClient {
	if ip == nil {
		ip = make([]string, 0)
	}

	roll := &IpRollClient{
		ips:    ip,
		useMap: make(map[string]*useInfo),
	}

	if len(ip) > 0 {
		roll.clientMap = make(map[string]*clientResource)
	}

	return roll
}

func (s *IpRollClient) GetHttpClient(urlStr string, proxy string, usejar bool) (*http.Client, error) {

	var clientres *clientResource
	if proxy != "" {
		proxyuri, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		if proxyTransport == nil {
			proxyTransport = &http.Transport{
				Dial:                defaultDialer.Dial,
				Proxy:               http.ProxyURL(proxyuri),
				MaxIdleConnsPerHost: defaultOption.MaxIdleConns,
				TLSHandshakeTimeout: defaultOption.TLSTimeout,
			}
		} else {
			proxyTransport.Proxy = http.ProxyURL(proxyuri)
		}
		if IsDebug() {
			log.Printf("[gohttp] url = %s, use proxy = %s\n", urlStr, proxy)
		}
		clientres = &clientResource{proxyTransport, defaultCookiejar}
	} else {

		uri, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		delay := time.Duration(0)

		//并发取的时候锁定
		s.useLock.Lock()
		use, ok := s.useMap[uri.Host]
		need_delay := GetHostDelay(uri.Host)
		if ok {
			//need_delay
			lastIndex := use.Index
			if len(s.ips) != 0 {
				use.Index = (use.Index + 1) % len(s.ips)
			}

			//使用同一个IP，则需要延迟
			if lastIndex == use.Index && need_delay > 0 {
				sub := time.Now().Sub(use.LastTime)
				if sub < need_delay {
					delay = need_delay - sub
				}
			}
			use.LastTime = time.Now().Add(delay)
		} else {
			use = &useInfo{
				Index:    0,
				LastTime: time.Now(),
			}
		}
		s.useMap[uri.Host] = use
		s.useLock.Unlock()

		if IsDebug() {
			if len(s.ips) == 0 {
				log.Printf("[gohttp] url = %s, delay = %dms, use default setting\n", urlStr, delay/time.Millisecond)
			} else {
				log.Printf("[gohttp] url = %s, delay = %dms, use ip = %s\n", urlStr, delay, s.ips[use.Index])
			}
		}

		if delay > 0 {
			time.Sleep(delay)
		}

		if len(s.ips) == 0 {
			clientres = &clientResource{defaultTransport, defaultCookiejar}
		} else {
			//
			//加锁并发
			ip := s.ips[use.Index]
			s.clientLock.Lock()
			if v, ok := s.clientMap[ip]; ok {
				clientres = v
			} else {
				clientres = &clientResource{MakeTransport(ip), MakeCookiejar()}
				s.clientMap[ip] = clientres
			}
			s.clientLock.Unlock()
		}
	}

	if usejar {
		return MakeClient(clientres.Transport, clientres.Jar), nil
	}
	return MakeClient(clientres.Transport, MakeCookiejar()), nil
}

func (s *IpRollClient) ResetCookie(uri *url.URL) {
	s.clientLock.Lock()
	for _, client := range s.clientMap {
		if client.Jar == nil {
			continue
		}
		cookies := client.Jar.Cookies(uri)
		for _, c := range cookies {
			c.Expires = time.Now().Add(-1 * time.Hour)
		}
		client.Jar.SetCookies(uri, cookies)
	}
	s.clientLock.Unlock()
}
