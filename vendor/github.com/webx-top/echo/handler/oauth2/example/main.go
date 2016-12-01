package main

import (
	"encoding/json"
	"flag"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/handler/oauth2"
	mw "github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/session"
	boltStore "github.com/webx-top/echo/middleware/session/engine/bolt"
	cookieStore "github.com/webx-top/echo/middleware/session/engine/cookie"
)

func main() {
	port := flag.String(`p`, "8080", "port")
	flag.Parse()
	e := echo.New()
	e.Use(mw.Log())
	sessionOptions := &echo.SessionOptions{
		Engine: `bolt`,
		Name:   `SESSIONID`,
		CookieOptions: &echo.CookieOptions{
			Path:     `/`,
			Domain:   ``,
			MaxAge:   0,
			Secure:   false,
			HttpOnly: true,
		},
	}

	cookieStore.RegWithOptions(&cookieStore.CookieOptions{
		KeyPairs: [][]byte{
			[]byte(`123456789012345678901234567890ab`),
		},
		SessionOptions: sessionOptions,
	})

	boltStore.RegWithOptions(&boltStore.BoltOptions{
		File: `./session.db`,
		KeyPairs: [][]byte{
			[]byte(`123456789012345678901234567890ab`),
		},
		BucketName:     `session`,
		SessionOptions: sessionOptions,
	})

	e.Use(session.Middleware(sessionOptions))

	e.Get("/", func(c echo.Context) error {
		return c.HTML(`Login: <a href="/oauth/login/github" target="_blank">github</a>`)
	})

	options := oauth2.New(`http://www.coscms.com`, oauth2.Config{
		GithubKey:    `9b168a10a77fbcafcdcf`,
		GithubSecret: `929bbf6136084052faf4f5768c14af805173ac27`,
	})
	options.Success(func(ctx echo.Context) error {
		user := options.User(ctx)
		b, e := json.MarshalIndent(user, "", "  ")
		if e != nil {
			return e
		}
		return ctx.String(string(b))
	})
	options.Wrapper(e)

	switch `` {
	case `fast`:
		// FastHTTP
		e.Run(fasthttp.New(":" + *port))

	default:
		// Standard
		e.Run(standard.New(":" + *port))
	}
}
