// Copyright 2013 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package imageproxy provides an image proxy server.  For typical use of
// creating and using a Proxy, see cmd/imageproxy/main.go.
package imageproxy

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	imageLib "git.webx.top/coscms/app/base/lib/image"
	"github.com/admpub/httpcache"
	"github.com/golang/glog"
)

func text2Image(content string, buf *bytes.Buffer) {
	pngBgImgPath := "data/fonts/background.png"
	fontPath := ""
	fontPath = "Courier New.ttf"
	if len(content) != len([]rune(content)) { //含多字节文字
		fontPath = "FangZhengBoldBlack.ttf"
	}
	fontPath = "data/fonts/" + fontPath
	img := imageLib.TextToImage(content, fontPath, pngBgImgPath)
	b := new(bytes.Buffer)
	png.Encode(b, img)
	imgB := b.Bytes()

	fmt.Fprintf(buf, "HTTP/1.1 200 OK\n")
	fmt.Fprintf(buf, "Content-Type: image/png\n")
	fmt.Fprintf(buf, "\n")
	buf.Write(imgB)
}

var DefaultResponse = func(resp *http.Response, err error, req *http.Request) (*http.Response, error) {

	buf := new(bytes.Buffer)

	req.URL.RawQuery = ""
	req.URL.Path = "/ERR.png"
	req.URL.Host = "coscms.com"
	req.URL.Scheme = ""
	req.URL.Fragment = ""

	if resp != nil {
		if resp.StatusCode == http.StatusNotFound {
			/*
				req.URL.Path = "/ERR404.png"
				textContent := "[404] Not Found.\nhttp://www.CosCMS.com"
				textContent += "\n" + "OK"
				text2Image(textContent, buf)
				return http.ReadResponse(bufio.NewReader(buf), req)
			*/

			img := []byte("got 404 while fetching remote image.")
			fmt.Fprintf(buf, "HTTP/1.1 404 Not Found\n")
			fmt.Fprintf(buf, "Content-Type: text/html; charset=utf-8\n")
			fmt.Fprintf(buf, "Content-Length: %d\n", len(img))
			fmt.Fprintf(buf, "\n")
			buf.Write(img)
			return http.ReadResponse(bufio.NewReader(buf), req)

		}
		switch resp.Header.Get("X-Valid-Image") {
		case "1":
		default:
			req.URL.Path = "/ERR501.png"
			fmt.Fprintf(buf, "HTTP/1.1 501 Not Implemented\n")
			fmt.Fprintf(buf, "Content-Type: text/html; charset=utf-8\n")
			fmt.Fprintf(buf, "\n")
			buf.Write([]byte{})
			return http.ReadResponse(bufio.NewReader(buf), req)
		}
	}

	txt := []byte(fmt.Sprintf("%v", err))
	fmt.Fprintf(buf, "HTTP/1.1 500 OK\n")
	fmt.Fprintf(buf, "Content-Type: text/html; charset=utf-8\n")
	fmt.Fprintf(buf, "Content-Length: %d\n", len(txt))
	fmt.Fprintf(buf, "\n")
	buf.Write(txt)
	return http.ReadResponse(bufio.NewReader(buf), req)

}

// Proxy serves image requests.
//
// Note that a Proxy should not be run behind a http.ServeMux, since the
// ServeMux aggressively cleans URLs and removes the double slash in the
// embedded request URL.
type Proxy struct {
	Client *http.Client // client used to fetch remote URLs
	Cache  Cache        // cache used to cache responses

	// Whitelist specifies a list of remote hosts that images can be
	// proxied from.  An empty list means all hosts are allowed.
	Whitelist []string

	// Referrers, when given, requires that requests to the image
	// proxy come from a referring host. An empty list means all
	// hosts are allowed.
	Referrers []string

	// DefaultBaseURL is the URL that relative remote URLs are resolved in
	// reference to.  If nil, all remote URLs specified in requests must be
	// absolute.
	DefaultBaseURL *url.URL

	// SignatureKey is the HMAC key used to verify signed requests.
	SignatureKey []byte

	// Allow images to scale beyond their original dimensions.
	ScaleUp bool

	//Clean prefix
	CleanPrefix string

	//网址查询字符串中要取值的变量名。如果设为“-”则采用整个查询字符串作为值（比如?aaa值为aaa）
	ResrcVKey string

	//验证所有参数加密后的令牌
	ValidToken bool

	//是否向资源网址传递参数
	CrossParams bool
}

// NewProxy constructs a new proxy.  The provided http RoundTripper will be
// used to fetch remote URLs.  If nil is provided, http.DefaultTransport will
// be used.
func NewProxy(transport http.RoundTripper, cache Cache) *Proxy {
	if transport == nil {
		transport = http.DefaultTransport
	}
	if cache == nil {
		cache = NopCache
	}

	proxy := Proxy{
		Cache: cache,
	}

	client := new(http.Client)
	client.Transport = &httpcache.Transport{
		Transport:           &TransformingTransport{transport, client, &proxy},
		Cache:               cache,
		MarkCachedResponses: true,
		ResponseErrorFunc:   DefaultResponse,
	}

	proxy.Client = client

	return &proxy
}

// ServeHTTP handles image requests.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if len(p.CleanPrefix) > 0 {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, p.CleanPrefix)
	}

	if r.URL.Path == "/favicon.ico" {
		return // ignore favicon requests
	}

	req, err := NewRequest(r, p.DefaultBaseURL, p.ResrcVKey, p.ValidToken, p.CrossParams)
	if err != nil {
		msg := fmt.Sprintf("invalid request URL: %v", err)
		glog.Error(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if !p.allowed(req) {
		msg := fmt.Sprintf("request does not contain an allowed host or valid signature")
		glog.Error(msg)
		http.Error(w, msg, http.StatusForbidden)
		return
	}
	reqURL := req.String()
	resp, err := p.Client.Get(reqURL)
	if err != nil {
		msg := fmt.Sprintf("error fetching remote image: %v", err)
		glog.Error(msg)
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	cached := resp.Header.Get(httpcache.XFromCache)
	glog.Infof("request: %v (served from cache: %v)", *req, cached == "1")

	copyHeader(w, resp, "Cache-Control")
	copyHeader(w, resp, "Last-Modified")
	copyHeader(w, resp, "Expires")
	copyHeader(w, resp, "Etag")

	if is304 := check304(r, resp); is304 {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	copyHeader(w, resp, "Content-Length")
	copyHeader(w, resp, "Content-Type")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (p *Proxy) BuildURL(option, rurl string) (r string) {
	if len(p.SignatureKey) <= 0 {
		r = option + "/" + url.QueryEscape(rurl)
		return
	}
	if p.ValidToken {
		text := option + "/" + rurl
		token := GenerateSignature(p.SignatureKey, []byte(text))
		r = token + "/" + option + "/" + url.QueryEscape(rurl)
		return
	}

	token := GenerateSignature(p.SignatureKey, []byte(rurl))
	r = option + "," + optSignaturePrefix + token + "/" + url.QueryEscape(rurl)
	return
}

func copyHeader(w http.ResponseWriter, r *http.Response, header string) {
	key := http.CanonicalHeaderKey(header)
	if value, ok := r.Header[key]; ok {
		w.Header()[key] = value
	}
}

// allowed returns whether the specified request is allowed because it matches
// a host in the proxy whitelist or it has a valid signature.
func (p *Proxy) allowed(r *Request) bool {
	if len(p.Referrers) > 0 && !validReferrer(p.Referrers, r.Original) {
		glog.Infof("request not coming from allowed referrer: %v", r)
		return false
	}

	if len(p.Whitelist) == 0 && len(p.SignatureKey) == 0 {
		return true // no whitelist or signature key, all requests accepted
	}

	if len(p.Whitelist) > 0 {
		if validHost(p.Whitelist, r.URL) {
			return true
		}
		glog.Infof("request is not for an allowed host: %v", r)
	}

	if len(p.SignatureKey) > 0 {
		if p.ValidToken {
			if validToken(p.SignatureKey, r) {
				return true
			}
			glog.Infof("token invalid: %v for %v", r.Token, r.Setting)
		} else if validSignature(p.SignatureKey, r) {
			return true
		}
		glog.Infof("request contains invalid signature: %v for %v", r.Options.Signature, r.URL.String())
	}

	return false
}

// validHost returns whether the host in u matches one of hosts.
func validHost(hosts []string, u *url.URL) bool {
	for _, host := range hosts {
		if u.Host == host {
			return true
		}
		if strings.HasPrefix(host, "*.") && strings.HasSuffix(u.Host, host[2:]) {
			return true
		}
	}

	return false
}

// returns whether the referrer from the request is in the host list.
func validReferrer(hosts []string, r *http.Request) bool {
	parsed, err := url.Parse(r.Header.Get("Referer"))
	if err != nil { // malformed or blank header, just deny
		return false
	}

	return validHost(hosts, parsed)
}

//资源网址签名，用于验证资源网址是否被人为修改，对于图片尺寸等参数不限制（有风险）
func validSignature(key []byte, r *Request) bool {
	sig := r.Options.Signature
	if m := len(sig) % 4; m != 0 { // add padding if missing
		sig += strings.Repeat("=", 4-m)
	}

	got, err := base64.URLEncoding.DecodeString(sig)
	if err != nil {
		glog.Errorf("error base64 decoding signature %q", r.Options.Signature)
		return false
	}
	want := Crypt(key, []byte(r.URL.String()))
	return hmac.Equal(got, want)
}

//所有参数(包含资源网址)加密后的令牌，用于验证参数是否被人为修改（防止恶意访问，安全性高）
func validToken(key []byte, r *Request) bool {
	//r.Token = GenerateSignature(key, []byte(r.Setting))
	//fmt.Println("token:", r.Token)

	sig := r.Token
	if m := len(sig) % 4; m != 0 { // add padding if missing
		sig += strings.Repeat("=", 4-m)
	}

	got, err := base64.URLEncoding.DecodeString(sig)
	if err != nil {
		glog.Errorf("error base64 decoding token %q", r.Token)
		return false
	}
	want := Crypt(key, []byte(r.Setting))
	//fmt.Println(r.Setting, "=>", string(want))

	return hmac.Equal(got, want)
}

func Crypt(key, text []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(text)
	want := mac.Sum(nil)
	return want
}

//GenerateSignature (p.SignatureKey,"http://www.coscms.com/logo.png")
func GenerateSignature(key, url []byte) string {
	want := Crypt(key, url)
	return strings.TrimRight(base64.URLEncoding.EncodeToString(want), "=")
}

// check304 checks whether we should send a 304 Not Modified in response to
// req, based on the response resp.  This is determined using the last modified
// time and the entity tag of resp.
func check304(req *http.Request, resp *http.Response) bool {
	// TODO(willnorris): if-none-match header can be a comma separated list
	// of multiple tags to be matched, or the special value "*" which
	// matches all etags
	etag := resp.Header.Get("Etag")
	if len(etag) > 0 && etag == req.Header.Get("If-None-Match") {
		return true
	}

	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get("Last-Modified"))
	if err != nil {
		return false
	}
	ifModSince, err := time.Parse(time.RFC1123, req.Header.Get("If-Modified-Since"))
	if err != nil {
		return false
	}
	if lastModified.Before(ifModSince) {
		return true
	}

	return false
}

// TransformingTransport is an implementation of http.RoundTripper that
// optionally transforms images using the options specified in the request URL
// fragment.
type TransformingTransport struct {
	// Transport is the underlying http.RoundTripper used to satisfy
	// non-transform requests (those that do not include a URL fragment).
	Transport http.RoundTripper

	// CachingClient is used to fetch images to be resized.  This client is
	// used rather than Transport directly in order to ensure that
	// responses are properly cached.
	CachingClient *http.Client

	// Proxy is used to access command line flag settings during roundtripping.
	Proxy *Proxy
}

func (t *TransformingTransport) Valid(resp *http.Response) error {
	if resp.StatusCode == http.StatusNotFound {
		msg := "got 404 while fetching remote image."
		glog.Error(msg)
		return errors.New(msg)
	}

	ct := strings.SplitN(resp.Header.Get("Content-Type"), ";", 2)
	switch ct[0] {
	case "image/jpeg", "image/png", "image/gif": //,"image/bmp","image/tiff"
		resp.Header.Set("X-Valid-Image", "1")
	default:
		resp.Header.Set("X-Valid-Image", "")
		msg := "This content type is not supported: " + ct[0]
		glog.Error(msg)
		return errors.New(msg)
	}
	return nil
}

// RoundTrip implements the http.RoundTripper interface.
func (t *TransformingTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	if req.URL.Fragment == "" {
		// normal requests pass through
		glog.Infof("fetching remote URL: %v", req.URL)
		resp, err := t.Transport.RoundTrip(req)
		if err == nil {
			err = t.Valid(resp)
		}
		return resp, err
	}

	u := *req.URL
	u.Fragment = ""
	reqURL := u.String()
	resp, err := t.CachingClient.Get(reqURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	err = t.Valid(resp)
	if err != nil {
		return resp, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	opt := ParseOptions(req.URL.Fragment)

	// assign static settings from proxy to options
	if t.Proxy != nil {
		opt.ScaleUp = t.Proxy.ScaleUp
	}
	img, err := Transform(b, opt)
	if err != nil {
		glog.Errorf("error transforming image: %v", err)
		img = b
	}

	// replay response with transformed image and updated content length
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%s %s\n", resp.Proto, resp.Status)
	resp.Header.WriteSubset(buf, map[string]bool{"Content-Length": true})
	fmt.Fprintf(buf, "Content-Length: %d\n\n", len(img))
	buf.Write(img)

	return http.ReadResponse(bufio.NewReader(buf), req)
}
