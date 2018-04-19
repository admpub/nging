package middleware

import (
	"net/http"
	"strings"

	"github.com/webx-top/echo"
)

type (
	// RedirectConfig defines the config for Redirect middleware.
	RedirectConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper `json:"-"`

		// Status code to be used when redirecting the request.
		// Optional. Default value http.StatusMovedPermanently.
		Code int `json:"code"`
	}
)

var (
	// DefaultRedirectConfig is the default Redirect middleware config.
	DefaultRedirectConfig = RedirectConfig{
		Skipper: echo.DefaultSkipper,
		Code:    http.StatusMovedPermanently,
	}
)

// HTTPSRedirect redirects HTTP requests to HTTPS.
// For example, http://webx.top will be redirect to https://webx.top.
//
// Usage `Echo#Pre(HTTPSRedirect())`
func HTTPSRedirect() echo.MiddlewareFuncd {
	return HTTPSRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `HTTPSRedirect()`.
func HTTPSRedirectWithConfig(config RedirectConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			if !req.IsTLS() {
				host := req.Host()
				uri := req.URI()
				return c.Redirect("https://"+host+uri, config.Code)
			}
			return next.Handle(c)
		}
	}
}

// HTTPSWWWRedirect redirects HTTP requests to WWW HTTPS.
// For example, http://webx.top will be redirect to https://www.webx.top.
//
// Usage `Echo#Pre(HTTPSWWWRedirect())`
func HTTPSWWWRedirect() echo.MiddlewareFuncd {
	return HTTPSWWWRedirectWithConfig(DefaultRedirectConfig)
}

// HTTPSWWWRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `HTTPSWWWRedirect()`.
func HTTPSWWWRedirectWithConfig(config RedirectConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			host := req.Host()
			uri := req.URI()
			if !req.IsTLS() && !strings.HasPrefix(host, `www.`) {
				return c.Redirect("https://www."+host+uri, http.StatusMovedPermanently)
			}
			return next.Handle(c)
		}
	}
}

// WWWRedirect redirects non WWW requests to WWW.
// For example, http://webx.top will be redirect to http://www.webx.top.
//
// Usage `Echo#Pre(WWWRedirect())`
func WWWRedirect() echo.MiddlewareFuncd {
	return WWWRedirectWithConfig(DefaultRedirectConfig)
}

// WWWRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `WWWRedirect()`.
func WWWRedirectWithConfig(config RedirectConfig) echo.MiddlewareFuncd {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			scheme := req.Scheme()
			host := req.Host()
			if !strings.HasPrefix(host, `www.`) {
				uri := req.URI()
				return c.Redirect(scheme+"://www."+host+uri, http.StatusMovedPermanently)
			}
			return next.Handle(c)
		}
	}
}

// NonWWWRedirect redirects WWW requests to non WWW.
// For example, http://www.webx.top will be redirect to http://webx.top.
//
// Usage `Echo#Pre(NonWWWRedirect())`
func NonWWWRedirect() echo.MiddlewareFuncd {
	return NonWWWRedirectWithConfig(DefaultRedirectConfig)
}

// NonWWWRedirectWithConfig returns a HTTPSRedirect middleware with config.
// See `NonWWWRedirect()`.
func NonWWWRedirectWithConfig(config RedirectConfig) echo.MiddlewareFuncd {
	if config.Skipper == nil {
		config.Skipper = DefaultTrailingSlashConfig.Skipper
	}
	if config.Code == 0 {
		config.Code = DefaultRedirectConfig.Code
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			req := c.Request()
			scheme := req.Scheme()
			host := req.Host()
			if strings.HasPrefix(host, `www.`) {
				uri := req.URI()
				return c.Redirect(scheme+"://"+host[4:]+uri, http.StatusMovedPermanently)
			}
			return next.Handle(c)
		}
	}
}
