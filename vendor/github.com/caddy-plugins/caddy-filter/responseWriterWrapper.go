package filter

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func newResponseWriterWrapperFor(delegate http.ResponseWriter, beforeFirstWrite func(*responseWriterWrapper) bool) *responseWriterWrapper {
	wrapper := &responseWriterWrapper{
		skipped:             false,
		delegate:            delegate,
		beforeFirstWrite:    beforeFirstWrite,
		statusSetAtDelegate: 0,
		bodyAllowed:         true,
		maximumBufferSize:   -1,
		header:              http.Header{},
	}
	for key, values := range delegate.Header() {
		for i, value := range values {
			if i == 0 {
				wrapper.header.Set(key, value)
			} else {
				wrapper.header.Add(key, value)
			}
		}
	}
	return wrapper
}

type responseWriterWrapper struct {
	skipped             bool
	delegate            http.ResponseWriter
	buffer              *bytes.Buffer
	beforeFirstWrite    func(*responseWriterWrapper) bool
	bodyAllowed         bool
	firstContentWritten bool
	headerSetAtDelegate bool
	statusSetAtDelegate int
	maximumBufferSize   int
	header              http.Header
}

func (instance *responseWriterWrapper) Header() http.Header {
	if instance.skipped {
		return instance.delegate.Header()
	}
	return instance.header
}

func (instance *responseWriterWrapper) WriteHeader(status int) {
	if instance.skipped {
		instance.delegate.WriteHeader(status)
	}
	instance.bodyAllowed = bodyAllowedForStatus(status)
	instance.statusSetAtDelegate = status
}

func (instance *responseWriterWrapper) Write(content []byte) (int, error) {
	if instance.skipped {
		return instance.delegate.Write(content)
	}

	if len(content) <= 0 {
		return 0, nil
	}

	if !instance.firstContentWritten {
		if instance.beforeFirstWrite(instance) {
			instance.buffer = new(bytes.Buffer)
		} else {
			instance.skipped = true
			instance.buffer = nil
		}
		instance.firstContentWritten = true
	}

	if instance.buffer == nil {
		if err := instance.writeHeadersToDelegate(200); err != nil {
			return 0, err
		}
		return instance.delegate.Write(content)
	}

	if (instance.maximumBufferSize >= 0) &&
		((instance.buffer.Len() + len(content)) > instance.maximumBufferSize) {
		_, err := instance.delegate.Write(instance.buffer.Bytes())
		if err != nil {
			return 0, err
		}
		instance.buffer = nil
		return instance.delegate.Write(content)
	}

	return instance.buffer.Write(content)
}

func (instance *responseWriterWrapper) selectStatus(def int) int {
	if instance.statusSetAtDelegate > 0 {
		return instance.statusSetAtDelegate
	}
	if def > 0 {
		return def
	}
	return 200
}

func (instance *responseWriterWrapper) writeToDelegate(content []byte, defStatus int) (int, error) {
	if !instance.headerSetAtDelegate {
		err := instance.writeHeadersToDelegate(defStatus)
		if err != nil {
			return 0, err
		}
	}
	return instance.delegate.Write(content)
}

func (instance *responseWriterWrapper) writeRecordedToDelegate(defStatus int) (int, error) {
	recorded := instance.recorded()
	return instance.writeToDelegate(recorded, defStatus)
}

func (instance *responseWriterWrapper) writeToDelegateAndEncodeIfRequired(content []byte, defStatus int) (int, error) {
	if !instance.isGzipEncoded() {
		return instance.writeToDelegate(content, defStatus)
	}
	if !instance.headerSetAtDelegate {
		err := instance.writeHeadersToDelegate(defStatus)
		if err != nil {
			return 0, err
		}
	}
	writer, err := gzip.NewWriterLevel(instance.delegate, gzip.BestCompression)
	if err != nil {
		return instance.writeToDelegate(content, defStatus)
	}
	return writer.Write(content)
}

func (instance *responseWriterWrapper) writeHeadersToDelegate(defStatus int) error {
	if instance.headerSetAtDelegate {
		return errors.New("headers already set at response")
	}
	instance.headerSetAtDelegate = true
	w := instance.delegate
	for key, values := range instance.header {
		for i, value := range values {
			if i == 0 {
				w.Header().Set(key, value)
			} else {
				w.Header().Add(key, value)
			}
		}
	}
	w.WriteHeader(instance.selectStatus(defStatus))
	return nil
}

func (instance *responseWriterWrapper) isBodyAllowed() bool {
	return instance.bodyAllowed
}

func (instance *responseWriterWrapper) isGzipEncoded() bool {
	contentEncoding := instance.Header().Get("Content-Encoding")
	return strings.ToLower(contentEncoding) == "gzip"
}

func (instance *responseWriterWrapper) wasSomethingRecorded() bool {
	return instance.buffer != nil && instance.buffer.Len() > 0
}

func (instance *responseWriterWrapper) isInterceptingRequired() bool {
	return !instance.skipped && instance.wasSomethingRecorded()
}

func (instance *responseWriterWrapper) recorded() []byte {
	buffer := instance.buffer
	if buffer == nil {
		return []byte{}
	}
	return buffer.Bytes()
}

// Hijack implements http.Hijacker. It simply wraps the underlying
// ResponseWriter's Hijack method if there is one, or returns an error.
func (instance *responseWriterWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := instance.delegate.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, httpserver.NonHijackerError{Underlying: instance.delegate}
}

// CloseNotify implements http.CloseNotifier.
// It just inherits the underlying ResponseWriter's CloseNotify method.
// It panics if the underlying ResponseWriter is not a CloseNotifier.
func (instance *responseWriterWrapper) CloseNotify() <-chan bool {
	if cn, ok := instance.delegate.(http.CloseNotifier); ok {
		return cn.CloseNotify()
	}
	panic(httpserver.NonCloseNotifierError{Underlying: instance.delegate})
}

func (instance *responseWriterWrapper) recordedAndDecodeIfRequired() []byte {
	result := instance.recorded()
	if !instance.isGzipEncoded() {
		return result
	}
	src := bytes.NewBuffer(result)
	gzipSrc, err := gzip.NewReader(src)
	if err != nil {
		return result
	}
	result, err = ioutil.ReadAll(gzipSrc)
	if err != nil {
		return result
	}
	instance.Header().Del("Content-Encoding")
	return result
}

func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == 204:
		return false
	case status == 304:
		return false
	}
	return true
}
