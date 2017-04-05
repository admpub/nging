package middleware

import (
	"fmt"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
)

type VisitorInfo struct {
	RealIP       string
	Time         time.Time
	Elapsed      time.Duration
	URI          string
	Method       string
	UserAgent    string
	Referer      string
	RequestSize  int64
	ResponseSize int64
	ResponseCode int
}

func Log(recv ...func(*VisitorInfo)) echo.MiddlewareFunc {
	var logging func(*VisitorInfo)
	if len(recv) > 0 {
		logging = recv[0]
	}
	if logging == nil {
		logger := log.GetLogger(`HTTP`)
		logging = func(v *VisitorInfo) {
			icon := "●"
			switch {
			case v.ResponseCode >= 500:
				icon = "▣"
			case v.ResponseCode >= 400:
				icon = "■"
			case v.ResponseCode >= 300:
				icon = "▲"
			}
			logger.Info(" " + icon + " " + fmt.Sprint(v.ResponseCode) + " " + v.RealIP + " " + v.Method + " " + v.URI + " " + v.Elapsed.String() + " " + fmt.Sprint(v.ResponseSize))
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
			info.URI = req.URI()
			info.ResponseSize = res.Size()
			info.ResponseCode = res.Status()
			logging(info)
			return nil
		})
	}
}
