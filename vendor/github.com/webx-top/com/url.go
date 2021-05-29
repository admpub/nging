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

// Base64Encode base64 encode
func Base64Encode(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

// Base64Decode base64 decode
func Base64Decode(str string) (string, error) {
	s, e := base64.StdEncoding.DecodeString(str)
	return string(s), e
}

// URLSafeBase64 base64字符串编码为URL友好的字符串
func URLSafeBase64(str string, encode bool) string {
	if encode { // 编码后处理
		str = strings.TrimRight(str, `=`)
		str = strings.Replace(str, `/`, `_`, -1)
		str = strings.Replace(str, `+`, `-`, -1)
		return str
	}
	// 解码前处理
	str = strings.Replace(str, `_`, `/`, -1)
	str = strings.Replace(str, `-`, `+`, -1)
	var missing = (4 - len(str)%4) % 4
	str += strings.Repeat(`=`, missing)
	return str
}

// SafeBase64Encode base64 encode
func SafeBase64Encode(str string) string {
	str = Base64Encode(str)
	return URLSafeBase64(str, true)
}

// SafeBase64Decode base64 decode
func SafeBase64Decode(str string) (string, error) {
	str = URLSafeBase64(str, false)
	return Base64Decode(str)
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
	siteURL := urlInfo.Scheme + `://` + urlInfo.Host
	if strings.HasPrefix(relURL, `/`) {
		return siteURL + relURL
	}
	for strings.HasPrefix(relURL, `./`) {
		relURL = strings.TrimPrefix(relURL, `./`)
	}
	urlPath := path.Dir(urlInfo.Path)
	for strings.HasPrefix(relURL, `../`) {
		urlPath = path.Dir(urlPath)
		relURL = strings.TrimPrefix(relURL, `../`)
	}
	return siteURL + path.Join(urlPath, relURL)
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
