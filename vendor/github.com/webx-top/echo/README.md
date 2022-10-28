# Echo
[![Build Status](https://travis-ci.org/webx-top/echo.svg?branch=master)](https://travis-ci.org/webx-top/echo) [![Go Report Card](https://goreportcard.com/badge/github.com/webx-top/echo)](https://goreportcard.com/report/github.com/webx-top/echo)
#### Echo is a fast and unfancy web framework for Go (Golang). Up to 10x faster than the rest.
This package need >= **go 1.13**

## Features

- Optimized HTTP router which smartly prioritize routes.
- Build robust and scalable RESTful APIs.
- Run with standard HTTP server or FastHTTP server.
- Group APIs.
- Extensible middleware framework.
- Define middleware at root, group or route level.
- Handy functions to send variety of HTTP responses.
- Centralized HTTP error handling.
- Template rendering with any template engine.
- Define your format for the logger.
- Highly customizable.

## Quick Start

### Installation

```sh
$ go get github.com/webx-top/echo
```

### Hello, World!

Create `server.go`

```go
package main

import (
	"net/http"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine/standard"
)

func main() {
	e := echo.New()
	e.Get("/", func(c echo.Context) error {
		return c.String("Hello, World!", http.StatusOK)
	})
	e.Run(standard.New(":1323"))
}
```

Start server

```sh
$ go run server.go
```

Browse to [http://localhost:1323](http://localhost:1323) and you should see
Hello, World! on the page.

### Routing

```go
e.Post("/users", saveUser)
e.Get("/users/:id", getUser)
e.Put("/users/:id", updateUser)
e.Delete("/users/:id", deleteUser)
e.Get("/user/<id:[\\d]+>", getUser)
```

### Path Parameters

```go
func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	id := c.Param("id")
	// or id := c.Paramx("id").Uint64()
}
```

### Query Parameters

`/show?team=x-men&member=wolverine`

```go
func show(c echo.Context) error {
	// Get team and member from the query string
	team := c.Query("team")
	member := c.Query("member")
	age := c.Queryx("age").Uint()
}
```

### Form `application/x-www-form-urlencoded`

`POST` `/save`

name | value
:--- | :---
name | Joe Smith
email | joe@labstack.com


```go
func save(c echo.Context) error {
	// Get name and email
	name := c.Form("name")
	email := c.Form("email")
	age := c.Formx("age").Uint()
}
```

### Form `multipart/form-data`

`POST` `/save`

name | value
:--- | :---
name | Joe Smith
email | joe@labstack.com
avatar | avatar

```go
func save(c echo.Context) error {
	// Get name and email
	name := c.Form("name")
	email := c.Form("email")

	//------------
	// Get avatar
	//------------
	_, err := c.SaveUploadedFile("avatar","./")
	return err
}
```

### Handling Request

- Bind `JSON` or `XML` payload into Go struct based on `Content-Type` request header.
- Render response as `JSON` or `XML` with status code.

```go
type User struct {
	Name  string `json:"name" xml:"name"`
	Email string `json:"email" xml:"email"`
}

e.Post("/users", func(c echo.Context) error {
	u := new(User)
	if err := c.MustBind(u); err != nil {
		return err
	}
	return c.JSON(u, http.StatusCreated)
	// or
	// return c.XML(u, http.StatusCreated)
})
```

### Static Content

Server any file from static directory for path `/static/*`.

```go
e.Use(mw.Static(&mw.StaticOptions{
	Root:"static", //存放静态文件的物理路径
	Path:"/static/", //网址访问静态文件的路径
	Browse:true, //是否显示文件列表
}))
```

### Middleware

```go
// Root level middleware
e.Use(middleware.Log())
e.Use(middleware.Recover())

// Group level middleware
g := e.Group("/admin")
g.Use(middleware.BasicAuth(func(username, password string) bool {
	if username == "joe" && password == "secret" {
		return true
	}
	return false
}))

// Route level middleware
track := func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		println("request to /users")
		return next.Handle(c)
	}
}
e.Get("/users", func(c echo.Context) error {
	return c.String("/users", http.StatusOK)
}, track)
```

### Cookie
```go
e.Get("/setcookie", func(c echo.Context) error {
	c.SetCookie("uid","1")
	return c.String("/setcookie: uid="+c.GetCookie("uid"), http.StatusOK)
})
```

### Session
```go
...
import (
	...
	"github.com/webx-top/echo/middleware/session"
	//boltStore "github.com/webx-top/echo/middleware/session/engine/bolt"
	cookieStore "github.com/webx-top/echo/middleware/session/engine/cookie"
)
...
sessionOptions := &echo.SessionOptions{
	Engine: `cookie`,
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
})

e.Use(session.Middleware(sessionOptions))

e.Get("/session", func(c echo.Context) error {
	c.Session().Set("uid",1).Save()
	return c.String(fmt.Sprintf("/session: uid=%v",c.Session().Get("uid")))
})
```

### Websocket
```go
...
import (
	...
	"github.com/admpub/websocket"
	"github.com/webx-top/echo"
	ws "github.com/webx-top/echo/handler/websocket"
)
...

e.AddHandlerWrapper(ws.HanderWrapper)

e.Get("/websocket", func(c *websocket.Conn, ctx echo.Context) error {
	//push(writer)
	go func() {
		var counter int
		for {
			if counter >= 10 { //测试只推10条
				return
			}
			time.Sleep(5 * time.Second)
			message := time.Now().String()
			ctx.Logger().Info(`Push message: `, message)
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				ctx.Logger().Error(`Push error: `, err.Error())
				return
			}
			counter++
		}
	}()

	//echo
	ws.DefaultExecuter(c, ctx)
	return nil
})
```
[More...](https://github.com/webx-top/echo/blob/master/handler/websocket/example/main.go)

### Sockjs
```go
...
import (
	...
	"github.com/webx-top/echo"
	"github.com/admpub/sockjs-go/v3/sockjs"
	ws "github.com/webx-top/echo/handler/sockjs"
)
...

options := ws.Options{
	Handle: func(c sockjs.Session) error {
		//push(writer)
		go func() {
			var counter int
			for {
				if counter >= 10 { //测试只推10条
					return
				}
				time.Sleep(5 * time.Second)
				message := time.Now().String()
				log.Info(`Push message: `, message)
				if err := c.Send(message); err != nil {
					log.Error(`Push error: `, err.Error())
					return
				}
				counter++
			}
		}()

		//echo
		ws.DefaultExecuter(c)
		return nil
	},
	Options: &sockjs.DefaultOptions,
	Prefix:  "/websocket",
}
options.Wrapper(e)
```
[More...](https://github.com/webx-top/echo/blob/master/handler/sockjs/example/main.go)

### Other Example

```go
package main

import (
	"net/http"

	"github.com/webx-top/echo"
	// "github.com/webx-top/echo/engine/fasthttp"
	"github.com/webx-top/echo/engine/standard"
	mw "github.com/webx-top/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(mw.Log())

	e.Get("/", func(c echo.Context) error {
		return c.String("Hello, World!")
	})
	e.Get("/echo/:name", func(c echo.Context) error {
		return c.String("Echo " + c.Param("name"))
	})
	
	e.Get("/std", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`standard net/http handleFunc`))
		w.WriteHeader(200)
	})

	// FastHTTP
	// e.Run(fasthttp.New(":4444"))

	// Standard
	e.Run(standard.New(":4444"))
}
```

[See other examples...](https://github.com/admpub/echo-example/blob/master/_v2/main.go)

## Middleware list
Middleware  | Import path | Description
:-----------|:------------|:-----------
[BasicAuth](https://github.com/webx-top/echo/blob/master/middleware/auth.go)  | github.com/webx-top/echo/middleware |HTTP basic authentication
[BodyLimit](https://github.com/webx-top/echo/blob/master/middleware/bodylimit.go)  | github.com/webx-top/echo/middleware |Limit request body
[Gzip](https://github.com/webx-top/echo/blob/master/middleware/compress.go)  | github.com/webx-top/echo/middleware |Send gzip HTTP response
[Secure](https://github.com/webx-top/echo/blob/master/middleware/secure.go)  | github.com/webx-top/echo/middleware |Protection against attacks
[CORS](https://github.com/webx-top/echo/blob/master/middleware/cors.go)  | github.com/webx-top/echo/middleware |Cross-Origin Resource Sharing
[CSRF](https://github.com/webx-top/echo/blob/master/middleware/csrf.go)  | github.com/webx-top/echo/middleware |Cross-Site Request Forgery
[Log](https://github.com/webx-top/echo/blob/master/middleware/log.go)  | github.com/webx-top/echo/middleware |Log HTTP requests
[MethodOverride](https://github.com/webx-top/echo/blob/master/middleware/methodOverride.go)  | github.com/webx-top/echo/middleware |Override request method
[Recover](https://github.com/webx-top/echo/blob/master/middleware/recover.go)  | github.com/webx-top/echo/middleware |Recover from panics
[HTTPSRedirect](https://github.com/webx-top/echo/blob/master/middleware/redirect.go)  | github.com/webx-top/echo/middleware |Redirect HTTP requests to HTTPS
[HTTPSWWWRedirect](https://github.com/webx-top/echo/blob/master/middleware/redirect.go)  | github.com/webx-top/echo/middleware |Redirect HTTP requests to WWW HTTPS
[WWWRedirect](https://github.com/webx-top/echo/blob/master/middleware/redirect.go)  | github.com/webx-top/echo/middleware |Redirect non WWW requests to WWW
[NonWWWRedirect](https://github.com/webx-top/echo/blob/master/middleware/redirect.go)  | github.com/webx-top/echo/middleware |Redirect WWW requests to non WWW
[AddTrailingSlash](https://github.com/webx-top/echo/blob/master/middleware/slash.go)  | github.com/webx-top/echo/middleware |Add trailing slash to the request URI
[RemoveTrailingSlash](https://github.com/webx-top/echo/blob/master/middleware/slash.go)  | github.com/webx-top/echo/middleware |Remove trailing slash from the request URI
[Static](https://github.com/webx-top/echo/blob/master/middleware/static.go)  | github.com/webx-top/echo/middleware |Serve static files
[MaxAllowed](https://github.com/webx-top/echo/blob/master/middleware/limit.go) | github.com/webx-top/echo/middleware | MaxAllowed limits simultaneous requests; can help with high traffic load
[RateLimit](https://github.com/webx-top/echo/tree/master/middleware/ratelimit) | github.com/webx-top/echo/middleware/ratelimit | Rate limiting HTTP requests
[Language](https://github.com/webx-top/echo/tree/master/middleware/language) | github.com/webx-top/echo/middleware/language | Multi-language support
[Session](https://github.com/webx-top/echo/blob/master/middleware/session/middleware.go)  | github.com/webx-top/echo/middleware/session | Sessions Manager
[JWT](https://github.com/webx-top/echo/blob/master/middleware/jwt/jwt.go)  | github.com/webx-top/echo/middleware/jwt | JWT authentication
[Markdown](https://github.com/webx-top/echo/blob/master/middleware/markdown/markdown.go)  | github.com/webx-top/echo/middleware/markdown | Markdown rendering
[Render](https://github.com/webx-top/echo/blob/master/middleware/render/middleware.go)  | github.com/webx-top/echo/middleware/render | HTML template rendering
[ReverseProxy](https://github.com/webx-top/reverseproxy/blob/master/middleware.go)  | github.com/webx-top/reverseproxy | Reverse proxy


## Handler Wrapper list
Wrapper     | Import path | Description
:-----------|:------------|:-----------
Websocket   |github.com/webx-top/echo/handler/websocket | [Example](https://github.com/webx-top/echo/blob/master/handler/websocket/example/main.go)
Sockjs      |github.com/webx-top/echo/handler/sockjs | [Example](https://github.com/webx-top/echo/blob/master/handler/sockjs/example/main.go)
Oauth2      |github.com/webx-top/echo/handler/oauth2 | [Example](https://github.com/webx-top/echo/blob/master/handler/oauth2/example/main.go)
Pprof      |github.com/webx-top/echo/handler/pprof | -


## Cases
- [Nging](https://github.com/admpub/nging)

## Credits
- [Vishal Rana](https://github.com/vishr) - Author
- [Hank Shen](https://github.com/admpub) - Author
- [Nitin Rana](https://github.com/nr17) - Consultant
- [Contributors](https://github.com/webx-top/echo/graphs/contributors)

## License

[Apache 2](https://github.com/webx-top/echo/blob/master/LICENSE)
