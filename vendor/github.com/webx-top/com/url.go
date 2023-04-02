// Copyright 2013 com authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package com

import (
	"encoding/base64"
	"math"
	"net/url"
	"path"
	"regexp"
	"strings"
)

// URLEncode url encode string, is + not %20
func URLEncode(str string) string {
	return url.QueryEscape(str)
}

// URLDecode url decode string
func URLDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// RawURLEncode rawurlencode()
func RawURLEncode(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}

// RawURLDecode rawurldecode()
func RawURLDecode(str string) (string, error) {
	return url.QueryUnescape(str)
}

// Base64Encode base64 encode
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64 decode
func Base64Decode(str string) (string, error) {
	s, e := base64.StdEncoding.DecodeString(str)
	return string(s), e
}

var urlSafeBase64EncodeReplacer = strings.NewReplacer(`/`, `_`, `+`, `-`)
var urlSafeBase64DecodeReplacer = strings.NewReplacer(`_`, `/`, `-`, `+`)

// URLSafeBase64 base64字符串编码为URL友好的字符串
func URLSafeBase64(str string, encode bool) string {
	if encode { // 编码后处理
		str = strings.TrimRight(str, `=`)
		str = urlSafeBase64EncodeReplacer.Replace(str)
		return str
	}
	// 解码前处理
	str = urlSafeBase64DecodeReplacer.Replace(str)
	var missing = (4 - len(str)%4) % 4
	if missing > 0 {
		str += strings.Repeat(`=`, missing)
	}
	return str
}

// SafeBase64Encode base64 encode
func SafeBase64Encode(str string) string {
	str = base64.URLEncoding.EncodeToString([]byte(str))
	str = strings.TrimRight(str, `=`)
	return str
}

// SafeBase64Decode base64 decode
func SafeBase64Decode(str string) (string, error) {
	var missing = (4 - len(str)%4) % 4
	if missing > 0 {
		str += strings.Repeat(`=`, missing)
	}
	b, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		return ``, err
	}
	return string(b), nil
}

// TotalPages 总页数
func TotalPages(totalRows uint, limit uint) uint {
	return uint(math.Ceil(float64(totalRows) / float64(limit)))
}

// Offset 根据页码计算偏移值
func Offset(page uint, limit uint) uint {
	if page == 0 {
		page = 1
	}
	return (page - 1) * limit
}

// AbsURL 获取页面内相对网址的绝对路径
func AbsURL(pageURL string, relURL string) string {
	if strings.Contains(relURL, `://`) {
		return relURL
	}
	urlInfo, err := url.Parse(pageURL)
	if err != nil {
		return ``
	}
	return AbsURLx(urlInfo, relURL, true)
}

func AbsURLx(pageURLInfo *url.URL, relURL string, onlyRelative ...bool) string {
	if (len(onlyRelative) == 0 || !onlyRelative[0]) && strings.Contains(relURL, `://`) {
		return relURL
	}
	siteURL := pageURLInfo.Scheme + `://` + pageURLInfo.Host
	if strings.HasPrefix(relURL, `/`) {
		return siteURL + relURL
	}
	for strings.HasPrefix(relURL, `./`) {
		relURL = strings.TrimPrefix(relURL, `./`)
	}
	urlPath := path.Dir(pageURLInfo.Path)
	for strings.HasPrefix(relURL, `../`) {
		urlPath = path.Dir(urlPath)
		relURL = strings.TrimPrefix(relURL, `../`)
	}
	return siteURL + path.Join(urlPath, relURL)
}

func URLSeparator(pageURL string) string {
	sep := `?`
	if strings.Contains(pageURL, sep) {
		sep = `&`
	}
	return sep
}

var localIPRegexp = regexp.MustCompile(`^127(?:\.(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)){3}$`)

// IsLocalhost 是否是本地主机
func IsLocalhost(host string) bool {
	switch host {
	case `localhost`:
		return true
	case `[::1]`:
		return true
	default:
		return localIPRegexp.MatchString(host)
	}
}

// SplitHost localhost:8080 => localhost
func SplitHost(hostport string) string {
	host, _ := SplitHostPort(hostport)
	return host
}

func SplitHostPort(hostport string) (host string, port string) {
	if strings.HasSuffix(hostport, `]`) {
		host = hostport
		return
	}
	sep := `]:`
	pos := strings.LastIndex(hostport, sep)
	if pos > -1 {
		host = hostport[0 : pos+1]
		if len(hostport) > pos+2 {
			port = hostport[pos+2:]
		}
		return
	}
	sep = `:`
	pos = strings.LastIndex(hostport, sep)
	if pos > -1 {
		host = hostport[0:pos]
		if len(hostport) > pos+1 {
			port = hostport[pos+1:]
		}
		return
	}
	host = hostport
	return
}

func WithURLParams(urlStr string, key string, value string, args ...string) string {
	if strings.Contains(urlStr, `?`) {
		urlStr += `&`
	} else {
		urlStr += `?`
	}
	urlStr += key + `=` + url.QueryEscape(value)
	var k string
	for i, j := 0, len(args); i < j; i++ {
		if i%2 == 0 {
			k = args[i]
			continue
		}
		urlStr += `&` + k + `=` + url.QueryEscape(args[i])
		k = ``
	}
	if len(k) > 0 {
		urlStr += `&` + k + `=`
	}
	return urlStr
}

func FullURL(domianURL string, myURL string) string {
	if IsFullURL(myURL) {
		return myURL
	}
	if !strings.HasPrefix(myURL, `/`) && !strings.HasSuffix(domianURL, `/`) {
		myURL = `/` + myURL
	}
	myURL = domianURL + myURL
	return myURL
}

func IsFullURL(purl string) bool {
	if len(purl) == 0 {
		return false
	}
	if purl[0] == '/' {
		return false
	}
	// find "://"
	firstPos := strings.Index(purl, `/`)
	if firstPos < 0 || firstPos == len(purl)-1 {
		return false
	}
	if firstPos > 1 && purl[firstPos-1] == ':' && purl[firstPos+1] == '/' {
		return true
	}
	return false
}
