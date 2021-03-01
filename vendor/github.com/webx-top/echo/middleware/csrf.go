package middleware

import (
	"crypto/subtle"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/random"
)

type (
	// CSRFConfig defines the config for CSRF middleware.
	CSRFConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper `json:"-"`

		// TokenLength is the length of the generated token.
		TokenLength uint8 `json:"token_length"`
		// Optional. Default value 32.

		// TokenLookup is a string in the form of "<source>:<key>" that is used
		// to extract token from the request.
		// Optional. Default value "header:X-CSRF-Token".
		// Possible values:
		// - "header:<name>"
		// - "form:<name>"
		// - "query:<name>"
		TokenLookup string `json:"token_lookup"`

		// Context key to store generated CSRF token into context.
		// Optional. Default value "csrf".
		ContextKey string `json:"context_key"`

		// Name of the CSRF cookie. This cookie will store CSRF token.
		// Optional. Default value "csrf".
		CookieName string `json:"cookie_name"`

		// Domain of the CSRF cookie.
		// Optional. Default value none.
		CookieDomain string `json:"cookie_domain"`

		// Path of the CSRF cookie.
		// Optional. Default value none.
		CookiePath string `json:"cookie_path"`

		// Max age (in seconds) of the CSRF cookie.
		// Optional. Default value 86400 (24hr).
		CookieMaxAge int `json:"cookie_max_age"`

		// Indicates if CSRF cookie is secure.
		// Optional. Default value false.
		CookieSecure bool `json:"cookie_secure"`

		// Indicates if CSRF cookie is HTTP only.
		// Optional. Default value false.
		CookieHTTPOnly bool `json:"cookie_http_only"`
	}

	// csrfTokenExtractor defines a function that takes `echo.Context` and returns
	// either a token or an error.
	csrfTokenExtractor func(echo.Context) (string, error)
)

var (
	// DefaultCSRFConfig is the default CSRF middleware config.
	DefaultCSRFConfig = CSRFConfig{
		Skipper:      echo.DefaultSkipper,
		TokenLength:  32,
		TokenLookup:  "header:" + echo.HeaderXCSRFToken,
		ContextKey:   "csrf",
		CookieName:   "_csrf",
		CookieMaxAge: 86400,
	}
	ErrCSRFTokenInvalid        = errors.New("csrf token is invalid")
	ErrCSRFTokenIsEmpty        = errors.New("empty csrf token")
	ErrCSRFTokenIsEmptyInForm  = fmt.Errorf("%w in form param", ErrCSRFTokenIsEmpty)
	ErrCSRFTokenIsEmptyInQuery = fmt.Errorf("%w in query param", ErrCSRFTokenIsEmpty)
)

// CSRF returns a Cross-Site Request Forgery (CSRF) middleware.
// See: https://en.wikipedia.org/wiki/Cross-site_request_forgery
func CSRF() echo.MiddlewareFuncd {
	c := DefaultCSRFConfig
	return CSRFWithConfig(c)
}

// CSRFWithConfig returns a CSRF middleware with config.
// See `CSRF()`.
func CSRFWithConfig(config CSRFConfig) echo.MiddlewareFuncd {
	if config.Skipper == nil {
		config.Skipper = DefaultCSRFConfig.Skipper
	}
	if config.TokenLength == 0 {
		config.TokenLength = DefaultCSRFConfig.TokenLength
	}
	if config.TokenLookup == "" {
		config.TokenLookup = DefaultCSRFConfig.TokenLookup
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultCSRFConfig.ContextKey
	}
	if config.CookieName == "" {
		config.CookieName = DefaultCSRFConfig.CookieName
	}
	if config.CookieMaxAge == 0 {
		config.CookieMaxAge = DefaultCSRFConfig.CookieMaxAge
	}

	// Initialize
	parts := strings.Split(config.TokenLookup, ":")
	extractor := csrfTokenFromHeader(parts[1])
	switch parts[0] {
	case "form":
		extractor = csrfTokenFromForm(parts[1])
	case "query":
		extractor = csrfTokenFromQuery(parts[1])
	}

	return func(next echo.Handler) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}
			req := c.Request()
			token := c.GetCookie(config.CookieName)

			if len(token) == 0 {
				// Generate token
				token = random.String(config.TokenLength)
			}

			switch req.Method() {
			case echo.GET, echo.HEAD, echo.OPTIONS, echo.TRACE:
			default:
				// Validate token only for requests which are not defined as 'safe' by RFC7231
				clientToken, err := extractor(c)
				if err != nil {
					return err
				}
				if !validateCSRFToken(token, clientToken) {
					return echo.NewHTTPError(http.StatusForbidden, ErrCSRFTokenInvalid.Error()).SetRaw(ErrCSRFTokenInvalid)
				}
			}
			// Set CSRF cookie
			cookie := echo.NewCookie(config.CookieName, token, c.CookieOptions())
			if len(config.CookiePath) > 0 {
				cookie.Path = config.CookiePath
			}
			if len(config.CookieDomain) > 0 {
				cookie.Domain = config.CookieDomain
			}
			echo.CookieMaxAge(cookie, config.CookieMaxAge)
			cookie.Secure = config.CookieSecure
			cookie.HttpOnly = config.CookieHTTPOnly

			// Store token in the context
			c.Set(config.ContextKey, token)

			// Protect clients from caching the response
			c.Response().Header().Add(echo.HeaderVary, echo.HeaderCookie)

			return next.Handle(c)
		}
	}
}

// csrfTokenFromForm returns a `csrfTokenExtractor` that extracts token from the
// provided request header.
func csrfTokenFromHeader(header string) csrfTokenExtractor {
	return func(c echo.Context) (string, error) {
		return c.Request().Header().Get(header), nil
	}
}

// csrfTokenFromForm returns a `csrfTokenExtractor` that extracts token from the
// provided form parameter.
func csrfTokenFromForm(param string) csrfTokenExtractor {
	return func(c echo.Context) (string, error) {
		token := c.Form(param)
		if token == "" {
			return "", ErrCSRFTokenIsEmptyInForm
		}
		return token, nil
	}
}

// csrfTokenFromQuery returns a `csrfTokenExtractor` that extracts token from the
// provided query parameter.
func csrfTokenFromQuery(param string) csrfTokenExtractor {
	return func(c echo.Context) (string, error) {
		token := c.Query(param)
		if token == "" {
			return "", ErrCSRFTokenIsEmptyInQuery
		}
		return token, nil
	}
}

func validateCSRFToken(token, clientToken string) bool {
	return subtle.ConstantTimeCompare([]byte(token), []byte(clientToken)) == 1
}
