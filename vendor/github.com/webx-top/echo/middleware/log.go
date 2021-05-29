package middleware

import (
	"fmt"
	"io"
	std "log"
	"time"

	"github.com/webx-top/echo"
)

type VisitorInfo struct {
	RealIP       string
	Time         time.Time
	Elapsed      time.Duration
	Scheme       string
	Host         string
	URI          string
	Method       string
	UserAgent    string
	Referer      string
	RequestSize  int64
	ResponseSize int64
	ResponseCode int
}

var DefaultLogWriter = GetDefaultLogWriter()

func Log(recv ...func(*VisitorInfo)) echo.MiddlewareFunc {
	return LogWithWriter(nil, recv...)
}

func LogWithWriter(writer io.Writer, recv ...func(*VisitorInfo)) echo.MiddlewareFunc {
	var logging func(*VisitorInfo)
	if len(recv) > 0 {
		logging = recv[0]
	}
	if writer == nil {
		writer = DefaultLogWriter
	}
	logger := std.New(writer, ``, 0)
	if logging == nil {
		logging = func(v *VisitorInfo) {
			logger.Println(":" + fmt.Sprint(v.ResponseCode) + ": " + v.RealIP + " " + v.Method + " " + v.Scheme + " " + v.Host + " " + v.URI + " " + v.Elapsed.String() + " " + fmt.Sprint(v.ResponseSize))
		}
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			req := c.Request()
			res := c.Response()
			info := &VisitorInfo{Time: time.Now()}
			if err := h.Handle(c); err != nil {
				c.Error(err)
			}
			info.RealIP = req.RealIP()
			info.UserAgent = req.UserAgent()
			info.Referer = req.Referer()
			info.RequestSize = req.Size()
			info.Elapsed = time.Now().Sub(info.Time)
			info.Method = req.Method()
			info.Host = req.Host()
			info.Scheme = req.Scheme()
			info.URI = req.URI()
			info.ResponseSize = res.Size()
			info.ResponseCode = res.Status()
			logging(info)
			return nil
		})
	}
}
