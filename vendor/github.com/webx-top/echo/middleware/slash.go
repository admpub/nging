package middleware

import (
	"github.com/webx-top/echo"
)

type (
	// TrailingSlashConfig defines the config for TrailingSlash middleware.
	TrailingSlashConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper `json:"-"`

		// Status code to be used when redirecting the request.
		// Optional, but when provided the request is redirected using this code.
		RedirectCode int `json:"redirect_code"`
	}
)

var (
	// DefaultTrailingSlashConfig is the default TrailingSlash middleware config.
	DefaultTrailingSlashConfig = TrailingSlashConfig{
		Skipper: echo.DefaultSkipper,
	}
)

// AddTrailingSlash returns a root level (before router) middleware which adds a
// trailing slash to the request `URL#Path`.
//
// Usage `Echo#Pre(AddTrailingSlash())`
func AddTrailingSlash() echo.MiddlewareFuncd {
	return AddTrailingSlashWithConfig(DefaultTrailingSlashConfig)
}

// AddTrailingSlashWithConfig returns a AddTrailingSlash middleware with config.
// See `AddTrailingSlash()`.
func AddTrailingSlashWithConfig(config TrailingSlashConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			url := req.URL()
			path := url.Path()
			qs := url.RawQuery()
			if path != "/" && path[len(path)-1] != '/' {
				path += "/"
				uri := path
				if qs != "" {
					uri += "?" + qs
				}

				// Redirect
				if config.RedirectCode != 0 {
					return c.Redirect(uri, config.RedirectCode)
				}

				// Forward
				req.SetURI(uri)
				url.SetPath(path)
			}
			return next.Handle(c)
		}
	}
}

// RemoveTrailingSlash returns a root level (before router) middleware which removes
// a trailing slash from the request URI.
//
// Usage `Echo#Pre(RemoveTrailingSlash())`
func RemoveTrailingSlash() echo.MiddlewareFuncd {
	return RemoveTrailingSlashWithConfig(TrailingSlashConfig{})
}

// RemoveTrailingSlashWithConfig returns a RemoveTrailingSlash middleware with config.
// See `RemoveTrailingSlash()`.
func RemoveTrailingSlashWithConfig(config TrailingSlashConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			url := req.URL()
			path := url.Path()
			qs := url.RawQuery()
			l := len(path) - 1
			if l >= 0 && path != "/" && path[l] == '/' {
				path = path[:l]
				uri := path
				if qs != "" {
					uri += "?" + qs
				}

				// Redirect
				if config.RedirectCode != 0 {
					return c.Redirect(uri, config.RedirectCode)
				}

				// Forward
				req.SetURI(uri)
				url.SetPath(path)
			}
			return next.Handle(c)
		}
	}
}
