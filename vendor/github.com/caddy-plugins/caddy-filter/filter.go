package filter

import (
	"github.com/caddyserver/caddy/caddyhttp/fastcgi"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"io"
	"net/http"
	"strconv"
)

const defaultMaxBufferSize = 10 * 1024 * 1024

type filterHandler struct {
	next              httpserver.Handler
	rules             []*rule
	maximumBufferSize int
}

func (instance filterHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) (int, error) {
	// Do not intercept if this is a websocket upgrade request.
	if request.Method == "GET" && request.Header.Get("Upgrade") == "websocket" {
		return instance.next.ServeHTTP(writer, request)
	}

	wrapper := newResponseWriterWrapperFor(writer, func(wrapper *responseWriterWrapper) bool {
		header := wrapper.Header()
		for _, rule := range instance.rules {
			if rule.matches(request, &header) {
				return true
			}
		}
		return false
	})
	wrapper.maximumBufferSize = instance.maximumBufferSize
	result, err := instance.next.ServeHTTP(wrapper, request)
	if wrapper.skipped {
		return result, err
	}
	var logError error
	if err != nil {
		var ok bool
		// This handles https://github.com/echocat/caddy-filter/issues/4
		// If the fastcgi module is used and the FastCGI server produces log output
		// this is send (by the FastCGI module) as an error. We have to check this and
		// handle this case of error in a special way.
		if logError, ok = err.(fastcgi.LogError); !ok {
			return result, err
		}
	}
	if !wrapper.isInterceptingRequired() || !wrapper.isBodyAllowed() {
		wrapper.writeHeadersToDelegate(result)
		return result, logError
	}
	if !wrapper.isBodyAllowed() {
		return result, logError
	}
	header := wrapper.Header()
	var body []byte
	bodyRetrieved := false
	for _, rule := range instance.rules {
		if rule.matches(request, &header) {
			if !bodyRetrieved {
				body = wrapper.recordedAndDecodeIfRequired()
				bodyRetrieved = true
			}
			body = rule.execute(request, &header, body)
		}
	}
	var n int
	if bodyRetrieved {
		oldContentLength := wrapper.Header().Get("Content-Length")
		if len(oldContentLength) > 0 {
			newContentLength := strconv.Itoa(len(body))
			wrapper.Header().Set("Content-Length", newContentLength)
		}
		n, err = wrapper.writeToDelegateAndEncodeIfRequired(body, result)
	} else {
		n, err = wrapper.writeRecordedToDelegate(result)
	}
	if err != nil {
		return result, err
	}
	if n < len(body) {
		return result, io.ErrShortWrite
	}
	return result, logError
}
