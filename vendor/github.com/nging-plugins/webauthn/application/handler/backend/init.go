package backend

import "github.com/webx-top/echo"

func Register(r echo.RouteRegister) {
	g := r.Group(`/user`)
	g.Get(`/webauthn`, WebAuthn)
}
