# hydra
[Hydra](https://github.com/ory-am/hydra) middleware for [Echo](https://github.com/webx-top/echo) framework.
It uses Hydra's API to extract and validate auth token.

## Example

``` go
import (
    "github.com/webx-top/echo"
    "github.com/webx-top/echo/engine/standard"
    "github.com/ory-am/hydra/firewall"
    hydraMW "github.com/webx-top/echo/middleware/hydra"
)

func handler(c echo.Context) error {
	ctx := c.Get("hydra").(*firewall.Context) // or hydraMW.GetContext(c)
	// Now you can access ctx.Subject etc.
	return nil
}

func main(){
	// Initialize Hydra
	hc, err := hydraMW.Connect(hydraMW.Options{
		ClientID     : "...",
		ClientSecret : "...",
		ClusterURL   : "",
	})
	if err != nil {
		panic(err)
	}

	// Use the middleware
 	e := echo.New()
	e.Get("/", handler, hydraMW.ScopesRequired(hc, nil, "scope1", "scope2"))
	e.Run(standard.New(":4444"))
}
```
