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
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	r "math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unsafe"
)

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// Md5 md5 hash string
func Md5(str string) string {
	m := md5.New()
	io.WriteString(m, str)
	return hex.EncodeToString(m.Sum(nil))
}

func ByteMd5(b []byte) string {
	m := md5.New()
	m.Write(b)
	return hex.EncodeToString(m.Sum(nil))
}

func Md5file(file string) string {
	barray, _ := ioutil.ReadFile(file)
	return ByteMd5(barray)
}

func Token(key string, val []byte, args ...string) string {
	hm := hmac.New(sha1.New, []byte(key))
	hm.Write(val)
	for _, v := range args {
		hm.Write([]byte(v))
	}
	return base64.URLEncoding.EncodeToString(hm.Sum(nil))
}

func Token256(key string, val []byte, args ...string) string {
	hm := hmac.New(sha256.New, []byte(key))
	hm.Write(val)
	for _, v := range args {
		hm.Write([]byte(v))
	}
	return base64.URLEncoding.EncodeToString(hm.Sum(nil))
}

func Encode(data interface{}, args ...string) ([]byte, error) {
	if len(args) > 0 && args[0] == `JSON` {
		return JSONEncode(data)
	}
	return GobEncode(data)
}

func Decode(data []byte, to interface{}, args ...string) error {
	if len(args) > 0 && args[0] == `JSON` {
		return JSONDecode(data, to)
	}
	return GobDecode(data, to)
}

func GobEncode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(&data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GobDecode(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

func JSONEncode(data interface{}, indents ...string) ([]byte, error) {
	if len(indents) > 0 && len(indents[0]) > 0 {
		return json.MarshalIndent(data, ``, indents[0])
	}
	return json.Marshal(data)
}

func JSONDecode(data []byte, to interface{}) error {
	return json.Unmarshal(data, to)
}

func sha(m hash.Hash, str string) string {
	io.WriteString(m, str)
	return hex.EncodeToString(m.Sum(nil))
}

// Sha1 sha1 hash string
func Sha1(str string) string {
	return sha(sha1.New(), str)
}

// Sha256 sha256 hash string
func Sha256(str string) string {
	return sha(sha256.New(), str)
}

// Ltrim trim space on left
func Ltrim(str string) string {
	return strings.TrimLeftFunc(str, unicode.IsSpace)
}

// Rtrim trim space on right
func Rtrim(str string) string {
	return strings.TrimRightFunc(str, unicode.IsSpace)
}

// Trim trim space in all string length
func Trim(str string) string {
	return strings.TrimSpace(str)
}

// StrRepeat repeat string times
func StrRepeat(str string, times int) string {
	return strings.Repeat(str, times)
}

// StrReplace replace find all occurs to string
func StrReplace(str string, find string, to string) string {
	return strings.Replace(str, find, to, -1)
}

// IsLetter returns true if the 'l' is an English letter.
func IsLetter(l uint8) bool {
	n := (l | 0x20) - 'a'
	if n >= 0 && n < 26 {
		return true
	}
	return false
}

// Expand replaces {k} in template with match[k] or subs[atoi(k)] if k is not in match.
func Expand(template string, match map[string]string, subs ...string) string {
	var p []byte
	var i int
	for {
		i = strings.Index(template, "{")
		if i < 0 {
			break
		}
		p = append(p, template[:i]...)
		template = template[i+1:]
		i = strings.Index(template, "}")
		if s, ok := match[template[:i]]; ok {
			p = append(p, s...)
		} else {
			j, _ := strconv.Atoi(template[:i])
			if j >= len(subs) {
				p = append(p, []byte("Missing")...)
			} else {
				p = append(p, subs[j]...)
			}
		}
		template = template[i+1:]
	}
	p = append(p, template...)
	return string(p)
}

// Reverse s string, support unicode
func Reverse(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, rune := range s {
		n--
		runes[n] = rune
	}
	return string(runes[n:])
}

// RandomCreateBytes generate random []byte by specify chars.
func RandomCreateBytes(n int, alphabets ...byte) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	var bytes = make([]byte, n)
	var randby bool
	if num, err := rand.Read(bytes); num != n || err != nil {
		r.Seed(time.Now().UnixNano())
		randby = true
	}
	for i, b := range bytes {
		if len(alphabets) == 0 {
			if randby {
				bytes[i] = alphanum[r.Intn(len(alphanum))]
			} else {
				bytes[i] = alphanum[b%byte(len(alphanum))]
			}
		} else {
			if randby {
				bytes[i] = alphabets[r.Intn(len(alphabets))]
			} else {
				bytes[i] = alphabets[b%byte(len(alphabets))]
			}
		}
	}
	return bytes
}

// Substr returns the substr from start to length.
func Substr(s string, dot string, lengthAndStart ...int) string {
	var start, length, argsLen, ln int
	argsLen = len(lengthAndStart)
	if argsLen > 0 {
		length = lengthAndStart[0]
	}
	if argsLen > 1 {
		start = lengthAndStart[1]
	}
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	ln = len(bt)
	if start > ln {
		start = start % ln
	}
	end := start + length
	if end > (ln - 1) {
		end = ln
	}
	if dot == "" || end == ln {
		return string(bt[start:end])
	}
	return string(bt[start:end]) + dot
}

func IsASCIIUpper(r rune) bool {
	return 'A' <= r && r <= 'Z'
}

func ToASCIIUpper(r rune) rune {
	if 'a' <= r && r <= 'z' {
		r -= ('a' - 'A')
	}
	return r
}

func IsAlpha(r rune) bool {
	if ('Z' < r || r < 'A') && ('z' < r || r < 'a') {
		return false
	}
	return true
}

func IsAlphaNumeric(r rune) bool {
	if ('Z' < r || r < 'A') && ('z' < r || r < 'a') && ('9' < r || r < '0') {
		return false
	}
	return true
}

func IsNumeric(r rune) bool {
	if '9' < r || r < '0' {
		return false
	}
	return true
}

// GonicCase : webxTop => webx_top
func GonicCase(name string) string {
	s := make([]rune, 0, len(name)+3)
	for idx, chr := range name {
		if IsASCIIUpper(chr) && idx > 0 {
			if !IsASCIIUpper(s[len(s)-1]) {
				s = append(s, '_')
			}
		}
		if !IsASCIIUpper(chr) && idx > 1 {
			l := len(s)
			if IsASCIIUpper(s[l-1]) && IsASCIIUpper(s[l-2]) {
				s = append(s, s[l-1])
				s[l-1] = '_'
			}
		}
		s = append(s, chr)
	}
	return strings.ToLower(string(s))
}

// TitleCase : webx_top => Webx_Top
func TitleCase(name string) string {
	var s []rune
	upNextChar := true
	name = strings.ToLower(name)
	for _, chr := range name {
		switch {
		case upNextChar:
			upNextChar = false
			chr = ToASCIIUpper(chr)
		case chr == '_':
			upNextChar = true
			continue
		}
		s = append(s, chr)
	}
	return string(s)
}

// SnakeCase : WebxTop => webx_top
func SnakeCase(name string) string {
	var s []rune
	for idx, chr := range name {
		if isUpper := IsASCIIUpper(chr); isUpper {
			if idx > 0 {
				s = append(s, '_')
			}
			chr -= ('A' - 'a')
		}
		s = append(s, chr)
	}
	return string(s)
}

// CamelCase : webx_top => webxTop
func CamelCase(s string) string {
	n := ""
	var capNext bool
	for _, v := range s {
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
				capNext = false
			} else {
				n += string(v)
			}
			continue
		}
		if v == '_' || v == ' ' {
			capNext = true
		} else {
			capNext = false
			n += string(v)
		}
	}
	return n
}

// PascalCase : webx_top => WebxTop
func PascalCase(s string) string {
	n := ""
	capNext := true
	for _, v := range s {
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
				capNext = false
			} else {
				n += string(v)
			}
			continue
		}
		if v == '_' || v == ' ' {
			capNext = true
		} else {
			capNext = false
			n += string(v)
		}
	}
	return n
}

// UpperCaseFirst : webx => Webx
func UpperCaseFirst(name string) string {
	s := []rune(name)
	if len(s) > 0 {
		s[0] = unicode.ToUpper(s[0])
		name = string(s)
	}
	return name
}

// LowerCaseFirst : WEBX => wEBX
func LowerCaseFirst(name string) string {
	s := []rune(name)
	if len(s) > 0 {
		s[0] = unicode.ToLower(s[0])
		name = string(s)
	}
	return name
}

func AddSlashes(s string, args ...rune) string {
	b := []rune{'\''}
	if len(args) > 0 {
		b = append(b, args...)
	}
	return AddCSlashes(s, b...)
}

func AddCSlashes(s string, b ...rune) string {
	r := []rune{}
	for _, v := range []rune(s) {
		if v == '\\' {
			r = append(r, '\\')
		} else {
			for _, f := range b {
				if v == f {
					r = append(r, '\\')
					break
				}
			}
		}
		r = append(r, v)
	}
	s = string(r)
	return s
}

// MaskString 0123456789 => 012****789
func MaskString(v string, width ...float64) string {
	size := len(v)
	if size < 1 {
		return ``
	}
	if size == 1 {
		return `*`
	}
	show := 0.3
	if len(width) > 0 {
		show = width[0]
	}
	showSize := int(float64(size) * show)
	if showSize < 1 {
		showSize = 1
	}
	hideSize := size - showSize*2
	rights := showSize + hideSize
	if rights > 0 && hideSize > 0 && rights < size && showSize < size {
		return v[0:showSize] + strings.Repeat(`*`, hideSize) + v[rights:]
	}
	if show < 0.5 {
		showSize = int(float64(size) * 0.5)
		if showSize < 1 {
			showSize = 1
		}
		hideSize = size - showSize
		if hideSize > 0 && showSize < size {
			return v[0:showSize] + strings.Repeat(`*`, hideSize)
		}
	}
	return v[0:1] + strings.Repeat(`*`, size-1)
}

// LeftPadZero 字符串指定长度，长度不足的时候左边补零
func LeftPadZero(input string, padLength int) string {
	return fmt.Sprintf(`%0*s`, padLength, input)
}

var (
	reSpaceLine     = regexp.MustCompile("([\\t\\s\r]*\n){2,}")
	BreakLine       = []byte("\n")
	BreakLineString = "\n"
)

func CleanSpaceLine(b []byte) []byte {
	return reSpaceLine.ReplaceAll(b, BreakLine)
}

func CleanSpaceLineString(b string) string {
	return reSpaceLine.ReplaceAllString(b, BreakLineString)
}
