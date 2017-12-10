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
)

// UrlEncode url encode string, is + not %20
func UrlEncode(str string) string {
	return url.QueryEscape(str)
}

// UrlDecode url decode string
func UrlDecode(str string) (string, error) {
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
