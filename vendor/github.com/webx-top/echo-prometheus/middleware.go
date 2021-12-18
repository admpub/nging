package echoprometheus

import (
	"reflect"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/webx-top/echo"
)

// Config responsible to configure middleware
type Config struct {
	Skipper             echo.Skipper
	Namespace           string
	Buckets             []float64
	Subsystem           string
	NormalizeHTTPStatus bool
	OnlyRoutePath       bool
}

const (
	httpRequestsCount    = "requests_total"
	httpRequestsDuration = "request_duration_seconds"
	notFoundPath         = "/not-found"
)

// DefaultConfig has the default instrumentation config
var DefaultConfig = Config{
	Skipper: func(c echo.Context) bool {
		return c.Request().URL().Path() == `/metrics`
	},
	Namespace: "echo",
	Subsystem: "http",
	Buckets: []float64{
		0.0005,
		0.001, // 1ms
		0.002,
		0.005,
		0.01, // 10ms
		0.02,
		0.05,
		0.1, // 100 ms
		0.2,
		0.5,
		1.0, // 1s
		2.0,
		5.0,
		10.0, // 10s
		15.0,
		20.0,
		30.0,
	},
	NormalizeHTTPStatus: true,
	OnlyRoutePath:       true,
}

func normalizeHTTPStatus(status int) string {
	if status < 200 {
		return "1xx"
	} else if status < 300 {
		return "2xx"
	} else if status < 400 {
		return "3xx"
	} else if status < 500 {
		return "4xx"
	}
	return "5xx"
}

var notFoundHandlerPointer = reflect.ValueOf(echo.NotFoundHandler).Pointer()

func isNotFoundHandler(handler echo.Handler) bool {
	return reflect.ValueOf(handler).Pointer() == notFoundHandlerPointer
}

// NewConfig returns a new config with default values
func NewConfig() Config {
	return DefaultConfig
}

// MetricsMiddleware returns an echo middleware with default config for instrumentation.
func MetricsMiddleware() echo.MiddlewareFuncd {
	return MetricsMiddlewareWithConfig(DefaultConfig)
}

// MetricsMiddlewareWithConfig returns an echo middleware for instrumentation.
func MetricsMiddlewareWithConfig(config Config) echo.MiddlewareFuncd {
	httpRequests := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.Namespace,
		Subsystem: config.Subsystem,
		Name:      httpRequestsCount,
		Help:      "Number of HTTP operations",
	}, []string{"status", "method", "handler"})
	prometheus.DefaultRegisterer.Unregister(httpRequests)
	prometheus.DefaultRegisterer.MustRegister(httpRequests)

	httpDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.Namespace,
		Subsystem: config.Subsystem,
		Name:      httpRequestsDuration,
		Help:      "Spend time by processing a route",
		Buckets:   config.Buckets,
	}, []string{"method", "handler"})
	prometheus.DefaultRegisterer.Unregister(httpDuration)
	prometheus.DefaultRegisterer.MustRegister(httpDuration)

	skipper := config.Skipper
	if skipper == nil {
		skipper = echo.DefaultSkipper
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if skipper(c) {
				return next.Handle(c)
			}
			req := c.Request()
			var path string

			// to avoid attack high cardinality of 404
			if isNotFoundHandler(next) {
				path = notFoundPath
			}

			if len(path) == 0 {
				if config.OnlyRoutePath {
					path = c.Path()
				} else {
					path = req.URL().Path()
				}
			}
			//c.Logger().Debug(path)

			timer := prometheus.NewTimer(httpDuration.WithLabelValues(req.Method(), path))
			err := next.Handle(c)
			timer.ObserveDuration()

			if err != nil {
				c.Error(err)
			}

			var status string
			if config.NormalizeHTTPStatus {
				status = normalizeHTTPStatus(c.Response().Status())
			} else {
				status = strconv.Itoa(c.Response().Status())
			}

			httpRequests.WithLabelValues(status, req.Method(), path).Inc()

			return err
		}
	}
}
