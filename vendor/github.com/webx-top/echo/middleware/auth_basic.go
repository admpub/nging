package middleware

import (
	"encoding/base64"
	"net/http"

	"github.com/webx-top/echo"
)

type (
	BasicValidateFunc func(string, string) bool
)

const (
	basic = "Basic"
)

// BasicAuth returns an HTTP basic authentication middleware.
//
// For valid credentials it calls the next handler.
// For invalid credentials, it sends "401 - Unauthorized" response.
func BasicAuth(fn BasicValidateFunc, skipper ...echo.Skipper) echo.MiddlewareFunc {
	var isSkiped echo.Skipper
	if len(skipper) > 0 {
		isSkiped = skipper[0]
	} else {
		isSkiped = echo.DefaultSkipper
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if isSkiped(c) {
				return h.Handle(c)
			}
			auth := c.Request().Header().Get(echo.HeaderAuthorization)
			l := len(basic)

			if len(auth) > l+1 && auth[:l] == basic {
				b, err := base64.StdEncoding.DecodeString(auth[l+1:])
				if err == nil {
					cred := string(b)
					for i := 0; i < len(cred); i++ {
						if cred[i] == ':' {
							// Verify credentials
							if fn(cred[:i], cred[i+1:]) {
								return h.Handle(c)
							}
						}
					}
				}
			}
			c.Response().Header().Set(echo.HeaderWWWAuthenticate, basic+" realm=Restricted")
			return echo.NewHTTPError(http.StatusUnauthorized).SetRaw(echo.ErrUnauthorized)
		})
	}
}
