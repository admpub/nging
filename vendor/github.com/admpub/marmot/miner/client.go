//
// 	Copyright 2017 by marmot author: gdccmcm14@live.com.
// 	Licensed under the Apache License, Version 2.0 (the "License");
// 	you may not use this file except in compliance with the License.
// 	You may obtain a copy of the License at
// 		http://www.apache.org/licenses/LICENSE-2.0
// 	Unless required by applicable law or agreed to in writing, software
// 	distributed under the License is distributed on an "AS IS" BASIS,
// 	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// 	See the License for the specific language governing permissions and
// 	limitations under the License
//

package miner

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/admpub/log"
	"golang.org/x/net/proxy"
)

var Debugf func(format string, a ...interface{}) = log.Debugf

// NewJar Cookie record Jar
func NewJar() *cookiejar.Jar {
	cookieJar, _ := cookiejar.New(nil)
	return cookieJar
}

// Default Client
var (
	// Save Cookie, No timeout!
	Client = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			Debugf("[GoWorker] Redirect:%v", req.URL)
			return nil
		},
		Jar: NewJar(),
	}

	// Not Save Cookie
	NoCookieClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			Debugf("[GoWorker] Redirect:%v", req.URL)
			return nil
		},
	}
)

// NewProxyClient New a Proxy client, Default save cookie, Can timeout
// We should support some proxy way such as http(s) or socks
func NewProxyClient(proxystring string) (*http.Client, error) {
	proxyURL, err := url.Parse(proxystring)
	if err != nil {
		return nil, err
	}

	prefix := strings.SplitN(proxystring, ":", 2)[0]

	// setup a http transport
	httpTransport := &http.Transport{}

	// http://
	// https://
	// socks5://
	switch prefix {
	case "http", "https":
		httpTransport.Proxy = http.ProxyURL(proxyURL)
	case "socks5":
		// create a socks5 dialer
		dialer, err := proxy.FromURL(proxyURL, proxy.Direct)
		if err != nil {
			return nil, err
		}
		httpTransport.Dial = dialer.Dial
	default:
		return nil, errors.New("this proxy way not allow:" + prefix)
	}

	// This a alone client, diff from global client.
	client := &http.Client{
		// Allow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			Debugf("[GoWorker] Redirect:%v", req.URL)
			return nil
		},
		// Allow proxy: http, https, socks5
		Transport: httpTransport,
		// Allow keep cookie
		Jar: NewJar(),
		// Allow Timeout
		Timeout: time.Second * time.Duration(DefaultTimeOut),
	}
	return client, nil
}

// NewClient New a client, diff from proxy client
func NewClient(timeout ...time.Duration) (*http.Client, error) {
	client := &http.Client{
		// Allow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			Debugf("[GoWorker] Redirect:%v", req.URL)
			return nil
		},
		Jar:     NewJar(),
		Timeout: time.Second * time.Duration(DefaultTimeOut),
	}
	if len(timeout) > 0 {
		client.Timeout = timeout[0]
	}
	return client, nil
}
