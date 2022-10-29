package render

import (
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/tplfunc"
)

type Config struct {
	TmplDir              string
	Theme                string
	Engine               string
	Style                string
	Reload               bool
	ParseStrings         map[string]string
	ParseStringFuncs     map[string]func() string
	ErrorPages           map[int]string
	ErrorProcessors      []ErrorProcessor
	DefaultHTTPErrorCode int
	StaticOptions        *middleware.StaticOptions
	Debug                bool
	renderer             driver.Driver
	errorPageFuncSetter  []echo.HandlerFunc
	FuncMapGlobal        map[string]interface{}
	RendererDo           []func(driver.Driver)
}

var DefaultFuncMapSkipper = func(c echo.Context) bool {
	return c.Format() != `html` && !c.IsAjax() && !c.IsPjax()
}

func (t *Config) SetRendererDo(rd ...func(driver.Driver)) *Config {
	t.RendererDo = rd
	return t
}

func (t *Config) AddRendererDo(rd ...func(driver.Driver)) *Config {
	if t.RendererDo == nil {
		t.RendererDo = []func(driver.Driver){}
	}
	t.RendererDo = append(t.RendererDo, rd...)
	return t
}

func (t *Config) Parser() func([]byte) []byte {
	if t.ParseStrings == nil {
		return nil
	}
	return func(b []byte) []byte {
		s := string(b)
		for oldVal, newVal := range t.ParseStrings {
			s = strings.Replace(s, oldVal, newVal, -1)
		}
		if t.ParseStringFuncs != nil {
			for oldVal, newVal := range t.ParseStringFuncs {
				s = strings.Replace(s, oldVal, newVal(), -1)
			}
		}
		return []byte(s)
	}
}

// NewRenderer 新建渲染接口
func (t *Config) NewRenderer(manager ...driver.Manager) driver.Driver {
	tmplDir := t.TmplDir
	if len(t.Theme) > 0 {
		tmplDir = filepath.Join(tmplDir, t.Theme)
	}
	renderer := New(t.Engine, tmplDir)
	if len(manager) > 0 && manager[0] != nil {
		renderer.SetManager(manager[0])
	}
	if t.RendererDo != nil {
		for _, rendererDo := range t.RendererDo {
			rendererDo(renderer)
		}
	}
	renderer.Init()
	renderer.SetContentProcessor(t.Parser())
	return renderer
}

func (t *Config) AddFuncSetter(set ...echo.HandlerFunc) *Config {
	if t.errorPageFuncSetter == nil {
		t.errorPageFuncSetter = make([]echo.HandlerFunc, len(DefaultOptions.SetFuncMap))
		copy(t.errorPageFuncSetter, DefaultOptions.SetFuncMap)
	}
	t.errorPageFuncSetter = append(t.errorPageFuncSetter, set...)
	return t
}

func (t *Config) SetFuncSetter(set ...echo.HandlerFunc) *Config {
	t.errorPageFuncSetter = set
	return t
}

func (t *Config) HTTPErrorHandler() echo.HTTPErrorHandler {
	opt := &Options{
		ErrorPages:           t.ErrorPages,
		ErrorProcessors:      t.ErrorProcessors,
		DefaultHTTPErrorCode: t.DefaultHTTPErrorCode,
	}
	opt.SetFuncSetter(t.errorPageFuncSetter...)
	return HTTPErrorHandler(opt)
}

func (t *Config) StaticMiddleware() interface{} {
	if t.StaticOptions != nil {
		return middleware.Static(t.StaticOptions)
	}
	return nil
}

func (t *Config) ApplyTo(e *echo.Echo, manager ...driver.Manager) *Config {
	if t.renderer != nil {
		t.renderer.Close()
	}
	e.SetHTTPErrorHandler(t.HTTPErrorHandler())
	staticMW := t.StaticMiddleware()
	if staticMW != nil {
		e.Use(staticMW)
	}
	renderer := t.MakeRenderer(manager...)
	e.SetRenderer(renderer)
	return t
}

func defaultTplFuncMap() map[string]interface{} {
	return tplfunc.TplFuncMap
}

func (t *Config) MakeRenderer(manager ...driver.Manager) driver.Driver {
	renderer := t.NewRenderer(manager...)
	if t.FuncMapGlobal == nil {
		renderer.SetFuncMap(defaultTplFuncMap)
	} else {
		renderer.SetFuncMap(func() map[string]interface{} { return t.FuncMapGlobal })
	}
	t.renderer = renderer
	return renderer
}

func (t *Config) Renderer() driver.Driver {
	return t.renderer
}

// ThemeDir 主题所在文件夹的路径
func (t *Config) ThemeDir(args ...string) string {
	if len(args) < 1 {
		return filepath.Join(t.TmplDir, t.Theme)
	}
	return filepath.Join(t.TmplDir, args[0])
}
