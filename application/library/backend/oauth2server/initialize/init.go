package initialize

import (

	//"github.com/golang-jwt/jwt"
	"encoding/gob"

	"github.com/admpub/oauth2/v4"
	"github.com/coscms/oauth2s"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"
	"github.com/admpub/nging/v5/application/library/backend/oauth2server"
	"github.com/admpub/nging/v5/application/library/config"
)

func init() {
	func() {
		defer recover()
		gob.Register(map[string][]string{})
	}()
	oauth2server.RoutePrefix = `/oauth2`
	handler.Register(func(r echo.RouteRegister) {
		oauth2server.Debug = !config.FromFile().Sys.IsEnv(`prod`)
		var tokenStore oauth2.TokenStore
		oauth2server.Default.Init(
			oauth2s.JWTMethod(nil),
			//oauth2s.JWTKey([]byte(config.FromFile().Cookie.HashKey)),
			//oauth2s.JWTMethod(jwt.SigningMethodHS512),
			oauth2s.ClientStore(oauth2server.DefaultAppClient),
			oauth2s.SetStore(tokenStore),
			oauth2s.SetHandler(&oauth2s.HandlerInfo{
				PasswordAuthorization: oauth2server.PasswordAuthorizationHandler,
				UserAuthorize:         oauth2server.UserAuthorizeHandler,
				InternalError:         oauth2server.InternalErrorHandler,
				ResponseError:         oauth2server.ResponseErrorHandler,
				RefreshingScope:       oauth2server.RefreshingScopeHandler,
				RefreshingValidation:  oauth2server.RefreshingValidationHandler,
			}),
		)
		oauth2server.Route(r)
	})
}
