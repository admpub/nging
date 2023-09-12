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
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

type RemoteError struct {
	Host string
	Err  error
}

func (e *RemoteError) Error() string {
	return e.Err.Error()
}

var UserAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/29.0.1541.0 Safari/537.36"

// HTTPGet gets the specified resource. ErrNotFound is returned if the
// server responds with status 404.
func HTTPGet(client *http.Client, url string, header http.Header) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	for k, vs := range header {
		req.Header[k] = vs
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &RemoteError{req.URL.Host, err}
	}
	if resp.StatusCode == 200 {
		return resp.Body, nil
	}
	resp.Body.Close()
	if resp.StatusCode == 404 { // 403 can be rate limit error.  || resp.StatusCode == 403 {
		err = NotFoundError{"Resource not found: " + url}
	} else {
		err = &RemoteError{req.URL.Host, fmt.Errorf("get %s -> %d", url, resp.StatusCode)}
	}
	return nil, err
}

// HTTPGetToFile gets the specified resource and writes to file.
// ErrNotFound is returned if the server responds with status 404.
func HTTPGetToFile(client *http.Client, url string, header http.Header, fileName string) error {
	rc, err := HTTPGet(client, url, header)
	if err != nil {
		return err
	}
	defer rc.Close()

	os.MkdirAll(filepath.Dir(fileName), os.ModePerm)
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, rc)
	if err != nil {
		return err
	}
	err = f.Sync()
	return err
}

// HTTPGetBytes gets the specified resource. ErrNotFound is returned if the server
// responds with status 404.
func HTTPGetBytes(client *http.Client, url string, header http.Header) ([]byte, error) {
	rc, err := HTTPGet(client, url, header)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	return ioutil.ReadAll(rc)
}

// HTTPGetJSON gets the specified resource and mapping to struct.
// ErrNotFound is returned if the server responds with status 404.
func HTTPGetJSON(client *http.Client, url string, v interface{}) error {
	rc, err := HTTPGet(client, url, nil)
	if err != nil {
		return err
	}
	defer rc.Close()
	err = json.NewDecoder(rc).Decode(v)
	if _, ok := err.(*json.SyntaxError); ok {
		err = NotFoundError{"JSON syntax error at " + url}
	}
	return err
}

// A RawFile describes a file that can be downloaded.
type RawFile interface {
	Name() string
	RawUrl() string
	Data() []byte
	SetData([]byte)
}

// FetchFiles fetches files specified by the rawURL field in parallel.
func FetchFiles(client *http.Client, files []RawFile, header http.Header) error {
	ch := make(chan error, len(files))
	for i := range files {
		go func(i int) {
			p, err := HTTPGetBytes(client, files[i].RawUrl(), nil)
			if err != nil {
				ch <- err
				return
			}
			files[i].SetData(p)
			ch <- nil
		}(i)
	}
	for range files {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

// FetchFilesCurl uses command `curl` to fetch files specified by the rawURL field in parallel.
func FetchFilesCurl(files []RawFile, curlOptions ...string) error {
	ch := make(chan error, len(files))
	for i := range files {
		go func(i int) {
			stdout, _, err := ExecCmd("curl", append(curlOptions, files[i].RawUrl())...)
			if err != nil {
				ch <- err
				return
			}

			files[i].SetData([]byte(stdout))
			ch <- nil
		}(i)
	}
	for range files {
		if err := <-ch; err != nil {
			return err
		}
	}
	return nil
}

// HTTPPost ==============================
func HTTPPost(client *http.Client, url string, body []byte, header http.Header) (io.ReadCloser, error) {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UserAgent)
	for k, vs := range header {
		req.Header[k] = vs
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &RemoteError{req.URL.Host, err}
	}
	if resp.StatusCode == 200 {
		return resp.Body, nil
	}
	resp.Body.Close()
	if resp.StatusCode == 404 { // 403 can be rate limit error.  || resp.StatusCode == 403 {
		err = NotFoundError{"Resource not found: " + url}
	} else {
		err = &RemoteError{req.URL.Host, fmt.Errorf("get %s -> %d", url, resp.StatusCode)}
	}
	return nil, err
}

func HTTPPostBytes(client *http.Client, url string, body []byte, header http.Header) ([]byte, error) {
	rc, err := HTTPPost(client, url, body, header)
	if err != nil {
		return nil, err
	}
	p, err := ioutil.ReadAll(rc)
	rc.Close()
	return p, nil
}

func HTTPPostJSON(client *http.Client, url string, body []byte, header http.Header) ([]byte, error) {
	if header == nil {
		header = http.Header{}
	}
	header.Add("Content-Type", "application/json")
	p, err := HTTPPostBytes(client, url, body, header)
	if err != nil {
		return []byte{}, err
	}
	return p, nil
}

// NewCookie is a helper method that returns a new http.Cookie object.
// Duration is specified in seconds. If the duration is zero, the cookie is permanent.
// This can be used in conjunction with ctx.SetCookie.
func NewCookie(name string, value string, args ...interface{}) *http.Cookie {
	var (
		alen     = len(args)
		age      int64
		path     string
		domain   string
		secure   bool
		httpOnly bool
	)
	switch alen {
	case 5:
		httpOnly, _ = args[4].(bool)
		fallthrough
	case 4:
		secure, _ = args[3].(bool)
		fallthrough
	case 3:
		domain, _ = args[2].(string)
		fallthrough
	case 2:
		path, _ = args[1].(string)
		fallthrough
	case 1:
		switch args[0].(type) {
		case int:
			age = int64(args[0].(int))
		case int64:
			age = args[0].(int64)
		case time.Duration:
			age = int64(args[0].(time.Duration))
		}
	}
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		MaxAge:   0,
		Secure:   secure,
		HttpOnly: httpOnly,
	}
	if age > 0 {
		cookie.Expires = time.Unix(time.Now().Unix()+age, 0)
	} else if age < 0 {
		cookie.Expires = time.Unix(1, 0)
	}
	return cookie
}

type HTTPClientOptions func(c *http.Client)

func HTTPClientWithTimeout(timeout time.Duration, options ...HTTPClientOptions) *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, timeout)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(timeout))
				return conn, nil
			},
			ResponseHeaderTimeout: timeout,
		},
	}
	for _, opt := range options {
		opt(client)
	}
	return client
}

// IsNetworkOrHostDown - if there was a network error or if the host is down.
// expectTimeouts indicates that *context* timeouts are expected and does not
// indicate a downed host. Other timeouts still returns down.
func IsNetworkOrHostDown(err error, expectTimeouts bool) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.Canceled) {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return !expectTimeouts
	}

	// We need to figure if the error either a timeout
	// or a non-temporary error.
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		switch urlErr.Err.(type) {
		case *net.DNSError, *net.OpError, net.UnknownNetworkError:
			return true
		}
	}
	var e net.Error
	if errors.As(err, &e) {
		if e.Timeout() {
			return true
		}
	}

	// Fallback to other mechanisms.
	switch {
	case strings.Contains(err.Error(), "Connection closed by foreign host"):
		return true
	case strings.Contains(err.Error(), "TLS handshake timeout"):
		// If error is - tlsHandshakeTimeoutError.
		return true
	case strings.Contains(err.Error(), "i/o timeout"):
		// If error is - tcp timeoutError.
		return true
	case strings.Contains(err.Error(), "connection timed out"):
		// If err is a net.Dial timeout.
		return true
	case strings.Contains(err.Error(), "connection refused"):
		// If err is connection refused
		return true

	case strings.Contains(strings.ToLower(err.Error()), "503 service unavailable"):
		// Denial errors
		return true
	}
	return false
}

func HTTPCanRetry(code int) bool {
	return code < 200 || (code > 299 && code < http.StatusInternalServerError)
}

func ParseHTTPRetryAfter(res http.ResponseWriter) time.Duration {
	r := res.Header().Get(`Retry-After`)
	return ParseRetryAfter(r)
}

func ParseRetryAfter(r string) time.Duration {
	if len(r) == 0 {
		return 0
	}
	if StrIsNumeric(r) {
		i := Int64(r)
		if i <= 0 {
			return 0
		}
		return time.Duration(i) * time.Second
	}
	t, err := time.Parse(time.RFC1123, r)
	if err != nil {
		log.Printf("failed to ParseRetryAfter(%q): %v\n", r, err)
		return 0
	}
	//fmt.Printf("%+v", t.String())
	if t.Before(time.Now()) {
		return 0
	}
	return time.Until(t)
}
