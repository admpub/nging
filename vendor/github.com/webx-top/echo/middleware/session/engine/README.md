# Session

Middleware support for echo, utilizing by
[admpub/sessions](https://github.com/admpub/sessions).

## Installation

```shell
go get github.com/webx-top/echo
```

## Usage

```go
package main

import (
    "github.com/webx-top/echo"
    "github.com/webx-top/echo/engine/standard"
    "github.com/webx-top/echo/middleware/session"
    cookieStore "github.com/webx-top/echo/middleware/session/engine/cookie"
)

func index(c echo.Context) error {
    session := c.Session()

    var count int
    v := session.Get("count")

    if v == nil {
        count = 0
    } else {
        count = v.(int)
        count += 1
    }

    session.Set("count", count)

    data := struct {
        Visit int
    }{
        Visit: count,
    }

    return c.JSON(http.StatusOK, data)
}

func main() {
    sessionOptions:=&echo.SessionOptions{
        Name:   `GOSESSIONID`,
        Engine: `cookie`,
        CookieOptions: &echo.CookieOptions{
            Path:     `/`,
            HttpOnly: true,
        },
    }
    cookieStore.RegWithOptions(&cookieStore.CookieOptions{
        KeyPairs: [][]byte{
            []byte("secret-key"),
        },
        SessionOptions: sessionOptions,
    })

    e := echo.New()

    // Attach middleware
    e.Use(session.Middleware(sessionOptions))

    // Routes
    e.Get("/", index)

    e.Run(standard.New(":8080"))
}
```
