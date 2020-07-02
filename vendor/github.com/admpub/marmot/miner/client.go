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
	"github.com/webx-top/com"
	"github.com/webx-top/com/httpClientOptions"
	"golang.org/x/net/proxy"
)

var (
	Debugf        = log.Debugf
	CheckRedirect = func(req *http.Request, via []*http.Request) error {
		Debugf("[GoWorker] Redirect:%v", req.URL)
		return nil
	}
)

// NewJar Cookie record Jar
func NewJar() *cookiejar.Jar {
	cookieJar, _ := cookiejar.New(nil)
	return cookieJar
}

// Default Client
var (
	// Save Cookie, No timeout!
	Client = NewClient()

	// Not Save Cookie
	NoCookieClient = NewHTTPClient(
		httpClientOptions.Timeout(time.Second*time.Duration(DefaultTimeOut)),
		httpClientOptions.CheckRedirect(CheckRedirect),
	)
)

// NewProxyClient New a Proxy client, Default save cookie, Can timeout
// We should support some proxy way such as http(s) or socks
func NewProxyClient(proxystring string, timeouts ...time.Duration) (*http.Client, error) {
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

	timeout := time.Second * time.Duration(DefaultTimeOut)
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}
	return NewHTTPClient(
		httpClientOptions.Transport(httpTransport),
		httpClientOptions.Timeout(timeout),
		httpClientOptions.CheckRedirect(CheckRedirect),
		httpClientOptions.CookieJar(NewJar()),
	), nil
}

// NewClient New a client, diff from proxy client
func NewClient(timeouts ...time.Duration) *http.Client {
	timeout := time.Second * time.Duration(DefaultTimeOut)
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}
	return NewHTTPClient(
		httpClientOptions.Timeout(timeout),
		httpClientOptions.CheckRedirect(CheckRedirect),
		httpClientOptions.CookieJar(NewJar()),
	)
}

// NewHTTPClient New a client
func NewHTTPClient(options ...com.HTTPClientOptions) *http.Client {
	client := &http.Client{}
	for _, option := range options {
		option(client)
	}
	return client
}
