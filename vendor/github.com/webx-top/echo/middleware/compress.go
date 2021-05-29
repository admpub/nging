package middleware

import (
	"bufio"
	"compress/gzip"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

type (
	// GzipConfig defines the config for Gzip middleware.
	GzipConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper `json:"-"`

		// Gzip compression level.
		// Optional. Default value -1.
		Level int `json:"level"`
	}

	gzipWriter struct {
		io.Writer
		engine.Response
	}
)

const (
	gzipScheme = "gzip"
)

var (
	// DefaultGzipConfig is the default Gzip middleware config.
	DefaultGzipConfig = &GzipConfig{
		Skipper: echo.DefaultSkipper,
		Level:   -1,
	}
)

func (w *gzipWriter) WriteHeader(code int) {
	if code == http.StatusNoContent {
		w.Header().Del(echo.HeaderContentEncoding)
	}
	w.Header().Del(echo.HeaderContentLength)
	w.WriteHeader(code)
}

func (w *gzipWriter) Write(b []byte) (int, error) {
	if len(w.Header().Get(echo.HeaderContentType)) == 0 {
		w.Header().Set(echo.HeaderContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w *gzipWriter) Flush() {
	w.Writer.(*gzip.Writer).Flush()
	if flusher, ok := w.Response.(http.Flusher); ok {
		flusher.Flush()
		return
	}
	if flusher, ok := w.StdResponseWriter().(http.Flusher); ok {
		flusher.Flush()
	}
}

func (w *gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hijacker, ok := w.Response.(http.Hijacker); ok {
		return hijacker.Hijack()
	}
	return w.StdResponseWriter().(http.Hijacker).Hijack()
}

func (w *gzipWriter) CloseNotify() <-chan bool {
	if closeNotifiler, ok := w.Response.(http.CloseNotifier); ok {
		return closeNotifiler.CloseNotify()
	}
	return w.StdResponseWriter().(http.CloseNotifier).CloseNotify()
}

func (w *gzipWriter) Push(target string, opts *http.PushOptions) error {
	if p, ok := w.Response.(http.Pusher); ok {
		return p.Push(target, opts)
	}
	return http.ErrNotSupported
}

// Gzip returns a middleware which compresses HTTP response using gzip compression
// scheme.
func Gzip(config ...*GzipConfig) echo.MiddlewareFunc {
	if len(config) < 1 || config[0] == nil {
		return GzipWithConfig(DefaultGzipConfig)
	}
	return GzipWithConfig(config[0])
}

// GzipWithConfig return Gzip middleware with config.
// See: `Gzip()`.
func GzipWithConfig(config *GzipConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultGzipConfig.Skipper
	}
	if config.Level == 0 {
		config.Level = DefaultGzipConfig.Level
	}
	pool := gzipCompressPool(config)
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if config.Skipper(c) {
				return h.Handle(c)
			}
			resp := c.Response()
			resp.Header().Add(echo.HeaderVary, echo.HeaderAcceptEncoding)
			if strings.Contains(c.Request().Header().Get(echo.HeaderAcceptEncoding), gzipScheme) {
				resp.Header().Add(echo.HeaderContentEncoding, gzipScheme)
				i := pool.Get()
				w, ok := i.(*gzip.Writer)
				if !ok {
					return echo.NewHTTPError(http.StatusInternalServerError, i.(error).Error()).SetRaw(i.(error))
				}
				rw := resp.Writer()
				w.Reset(rw)
				defer func() {
					if resp.Size() == 0 {
						if resp.Header().Get(echo.HeaderContentEncoding) == gzipScheme {
							resp.Header().Del(echo.HeaderContentEncoding)
						}
						// We have to reset response to it's pristine state when
						// nothing is written to body or error is returned.
						// See issue #424, #407.
						resp.SetWriter(rw)
						w.Reset(ioutil.Discard)
					}
					w.Close()
					pool.Put(w)
				}()
				resp.SetWriter(&gzipWriter{Writer: w, Response: resp})
			}
			return h.Handle(c)
		})
	}
}

func gzipCompressPool(config *GzipConfig) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			w, err := gzip.NewWriterLevel(ioutil.Discard, config.Level)
			if err != nil {
				return err
			}
			return w
		},
	}
}
