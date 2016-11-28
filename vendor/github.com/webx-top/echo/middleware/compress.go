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

var (
	// DefaultGzipConfig is the default Gzip middleware config.
	DefaultGzipConfig = &GzipConfig{
		Skipper: echo.DefaultSkipper,
		Level:   -1,
	}
)

func (w gzipWriter) Write(b []byte) (int, error) {
	if w.Header().Get(echo.HeaderContentType) == `` {
		w.Header().Set(echo.HeaderContentType, http.DetectContentType(b))
	}
	return w.Writer.Write(b)
}

func (w gzipWriter) Flush() error {
	return w.Writer.(*gzip.Writer).Flush()
}

func (w gzipWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.Response.(http.Hijacker).Hijack()
}

func (w *gzipWriter) CloseNotify() <-chan bool {
	return w.Response.(http.CloseNotifier).CloseNotify()
}

var writerPool = sync.Pool{
	New: func() interface{} {
		return gzip.NewWriter(ioutil.Discard)
	},
}

// Gzip returns a middleware which compresses HTTP response using gzip compression
// scheme.
func Gzip() echo.MiddlewareFunc {
	return GzipWithConfig(DefaultGzipConfig)
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

	pool := gzipPool(config)
	scheme := `gzip`

	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if config.Skipper(c) {
				return h.Handle(c)
			}
			resp := c.Response()
			resp.Header().Add(echo.HeaderVary, echo.HeaderAcceptEncoding)
			if strings.Contains(c.Request().Header().Get(echo.HeaderAcceptEncoding), scheme) {
				rw := resp.Writer()
				w := pool.Get().(*gzip.Writer)
				w.Reset(rw)
				defer func() {
					if resp.Size() == 0 {
						// We have to reset response to it's pristine state when
						// nothing is written to body or error is returned.
						// See issue #424, #407.
						resp.SetWriter(rw)
						resp.Header().Del(echo.HeaderContentEncoding)
						w.Reset(ioutil.Discard)
					}
					w.Close()
					pool.Put(w)
				}()
				gw := gzipWriter{Writer: w, Response: resp}
				resp.Header().Set(echo.HeaderContentEncoding, scheme)
				resp.SetWriter(gw)
			}
			return h.Handle(c)
		})
	}
}

func gzipPool(config *GzipConfig) sync.Pool {
	return sync.Pool{
		New: func() interface{} {
			w, _ := gzip.NewWriterLevel(ioutil.Discard, config.Level)
			return w
		},
	}
}
