package render

import (
	"net/http"

	"github.com/webx-top/echo"
)

var (
	DefaultDataKey    = `data`
	DefaultTmplKey    = `tmpl`
	DefaultTmplName   = `index`
	DefaultErrorTmpl  = `error`
	JSONPCallbackName = `callback`
	DefaultErrorFunc  = OutputError
	DefaultOutputFunc = Output
)

// Middleware set renderer
func Middleware(d echo.Renderer) echo.MiddlewareFunc {
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			c.SetRenderer(d)
			return h.Handle(c)
		})
	}
}

// AutoOutput Outputs the specified format
func AutoOutput(outputFunc func(string, echo.Context) error, skipper ...echo.Skipper) echo.MiddlewareFunc {
	isSkiped := echo.DefaultSkipper
	if len(skipper) > 0 {
		isSkiped = skipper[0]
	}
	if outputFunc == nil {
		outputFunc = DefaultOutputFunc
	}
	return func(h echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if isSkiped(c) {
				return h.Handle(c)
			}
			format := c.Format()
			if err := h.Handle(c); err != nil {
				return DefaultErrorFunc(err, format, c)
			}
			return outputFunc(format, c)
		})
	}
}

// Output Outputs the specified format
func Output(format string, c echo.Context) error {
	switch format {
	case `json`:
		return c.JSON(c.Get(DefaultDataKey))
	case `jsonp`:
		return c.JSONP(c.Query(JSONPCallbackName), c.Get(DefaultDataKey))
	case `xml`:
		return c.XML(c.Get(DefaultDataKey))
	default:
		tmpl, ok := c.Get(DefaultTmplKey).(string)
		if !ok {
			tmpl = DefaultTmplName
		}
		data := c.Get(DefaultDataKey)
		if v, y := data.(*echo.Data); y {
			SetFuncs(c, v)
			return c.Render(tmpl, v.Data)
		}
		if h, y := data.(echo.H); y {
			v := h.ToData()
			SetFuncs(c, v)
			return c.Render(tmpl, v.Data)
		}
		return c.Render(tmpl, data)
	}
}

// SetFuncs register template function
func SetFuncs(c echo.Context, v *echo.Data) {
	c.SetFunc(`Info`, func() interface{} {
		return v.Info
	})
	c.SetFunc(`Code`, func() interface{} {
		return v.Code
	})
	c.SetFunc(`Zone`, func() interface{} {
		return v.Zone
	})
}

// OutputError Outputs the specified format
func OutputError(err error, format string, c echo.Context) error {
	if apiData, ok := err.(*echo.Data); ok {
		c.Set(DefaultDataKey, apiData)
	} else {
		c.Set(DefaultDataKey, echo.NewData(c.Code(), err.Error()))
	}
	c.Set(DefaultTmplKey, DefaultErrorTmpl)
	c.SetCode(http.StatusOK)
	return Output(format, c)
}
