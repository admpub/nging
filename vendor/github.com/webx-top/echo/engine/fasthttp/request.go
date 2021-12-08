//go:build !appengine
// +build !appengine

package fasthttp

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/admpub/fasthttp"
	"github.com/admpub/realip"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

type Request struct {
	requestMu  sync.RWMutex
	context    *fasthttp.RequestCtx
	url        engine.URL
	header     engine.Header
	value      *Value
	realIP     string
	stdRequest *http.Request
}

func NewRequest(c *fasthttp.RequestCtx) *Request {
	req := &Request{
		context: c,
		url:     &URL{url: c.URI()},
		header: &RequestHeader{
			header: &c.Request.Header,
			stdhdr: nil,
		},
	}
	req.value = NewValue(req)
	return req
}

func (r *Request) Context() context.Context {
	return r.context
}

func (r *Request) WithContext(ctx context.Context) *http.Request {
	return r.StdRequest().WithContext(ctx)
}

func (r *Request) SetValue(key string, value interface{}) {
	r.requestMu.Lock()
	r.context.SetUserValue(key, value)
	r.requestMu.Unlock()
}

func (r *Request) Host() string {
	return engine.Bytes2str(r.context.Host())
}

func (r *Request) URI() string {
	return engine.Bytes2str(r.context.RequestURI())
}

// SetURI implements `engine.Request#SetURI` function.
func (r *Request) SetURI(uri string) {
	r.context.Request.Header.SetRequestURI(uri)
}

func (r *Request) URL() engine.URL {
	return r.url
}

func (r *Request) Header() engine.Header {
	return r.header
}

func (r *Request) Proto() string {
	return "HTTP/1.1"
}

func (r *Request) RemoteAddress() string {
	return r.context.RemoteAddr().String()
}

// RealIP implements `engine.Request#RealIP` function.
func (r *Request) RealIP() string {
	if len(r.realIP) > 0 {
		return r.realIP
	}
	r.realIP = realip.XRealIP(r.header.Get(echo.HeaderXRealIP), r.header.Get(echo.HeaderXForwardedFor), r.RemoteAddress())
	return r.realIP
}

func (r *Request) Method() string {
	return engine.Bytes2str(r.context.Method())
}

func (r *Request) SetMethod(method string) {
	r.context.Request.Header.SetMethod(method)
}

func (r *Request) Body() io.ReadCloser {
	return ioutil.NopCloser(bytes.NewBuffer(r.context.PostBody()))
}

// SetBody implements `engine.Request#SetBody` function.
func (r *Request) SetBody(reader io.Reader) {
	r.context.Request.SetBodyStream(reader, 0)
}

func (r *Request) FormValue(name string) string {
	//return string(r.context.FormValue(name))
	return r.Form().Get(name)
}

func (r *Request) Form() engine.URLValuer {
	return r.value
}

func (r *Request) PostForm() engine.URLValuer {
	return r.value.postArgs
}

func (r *Request) MultipartForm() *multipart.Form {
	if !strings.HasPrefix(string(r.context.Request.Header.ContentType()), echo.MIMEMultipartForm) {
		return nil
	}
	re, err := r.context.MultipartForm()
	if err != nil {
		r.context.Logger().Printf(err.Error())
	}
	return re
}

func (r *Request) IsTLS() bool {
	return r.context.IsTLS()
}

func (r *Request) Cookie(key string) string {
	return engine.Bytes2str(r.context.Request.Header.Cookie(key))
}

func (r *Request) Referer() string {
	return engine.Bytes2str(r.context.Referer())
}

func (r *Request) UserAgent() string {
	return engine.Bytes2str(r.context.UserAgent())
}

func (r *Request) Object() interface{} {
	return r.context
}

func (r *Request) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	fileHeader, err := r.context.FormFile(key)
	if err != nil {
		return nil, nil, err
	}
	var file multipart.File
	file, err = fileHeader.Open()
	return file, fileHeader, err
}

func (r *Request) Scheme() string {
	if r.IsTLS() {
		return echo.SchemeHTTPS
	}
	scheme := engine.Bytes2str(r.context.URI().Scheme())
	if len(scheme) > 0 {
		return scheme
	}
	return echo.SchemeHTTP
}

// Size implements `engine.Request#ContentLength` function.
func (r *Request) Size() int64 {
	return int64(r.context.Request.Header.ContentLength())
}

func (r *Request) reset(c *fasthttp.RequestCtx, h engine.Header, u engine.URL) {
	r.requestMu = sync.RWMutex{}
	r.context = c
	r.header = h
	r.url = u
	r.value = NewValue(r)
	r.realIP = ``
	r.stdRequest = nil
}

// BasicAuth returns the username and password provided in the request's
// Authorization header, if the request uses HTTP Basic Authentication.
// See RFC 2617, Section 2.
func (r *Request) BasicAuth() (username, password string, ok bool) {
	auth := r.Header().Get(echo.HeaderAuthorization)
	if auth == "" {
		return
	}
	return parseBasicAuth(auth)
}

// SetHost implements `engine.Request#SetHost` function.
func (r *Request) SetHost(host string) {
	r.context.Request.SetHost(host)
}

func (r *Request) StdRequest() *http.Request {
	if r.stdRequest != nil {
		return r.stdRequest
	}
	var req http.Request
	ctx := r.context
	req.Method = r.Method()
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1
	req.RequestURI = r.URI()
	req.ContentLength = r.Size()
	req.Host = r.Host()
	req.RemoteAddr = r.RemoteAddress()
	req.TLS = ctx.TLSConnectionState()

	hdr := make(http.Header)
	ctx.Request.Header.VisitAll(func(k, v []byte) {
		sk := engine.Bytes2str(k)
		sv := engine.Bytes2str(v)
		switch sk {
		case "Transfer-Encoding":
			req.TransferEncoding = append(req.TransferEncoding, sv)
		default:
			hdr.Set(sk, sv)
		}
	})
	req.Header = hdr
	req.Body = r.Body()
	rURL, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		ctx.Logger().Printf("cannot parse requestURI %q: %s", req.RequestURI, err)
	}
	req.URL = rURL
	r.stdRequest = req.WithContext(r.context)
	return r.stdRequest
}

// parseBasicAuth parses an HTTP Basic Authentication string.
// "Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==" returns ("Aladdin", "open sesame", true).
func parseBasicAuth(auth string) (username, password string, ok bool) {
	const prefix = "Basic "
	if !strings.HasPrefix(auth, prefix) {
		return
	}
	c, err := base64.StdEncoding.DecodeString(auth[len(prefix):])
	if err != nil {
		return
	}
	cs := string(c)
	s := strings.IndexByte(cs, ':')
	if s < 0 {
		return
	}
	return cs[:s], cs[s+1:], true
}
