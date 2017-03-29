package render

import (
	"path/filepath"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/mvc/static/resource"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/echo/middleware/tplfunc"
)

type Config struct {
	TmplDir       string
	Theme         string
	Engine        string
	Style         string
	Reload        bool
	ParseStrings  map[string]string
	ErrorPages    map[int]string
	StaticOptions *middleware.StaticOptions
	Debug         bool
	renderer      driver.Driver
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
	renderer.Init(true, t.Reload)
	if len(manager) > 0 {
		renderer.SetManager(manager[0])
	}
	renderer.SetContentProcessor(t.Parser())
	if t.StaticOptions != nil {
		st := t.NewStatic()
		renderer.SetFuncMap(func() map[string]interface{} {
			return st.Register(nil)
		})
		absTmplPath, err := filepath.Abs(tmplDir)
		var absFilePath string
		if err == nil {
			absFilePath, err = filepath.Abs(t.StaticOptions.Root)
		}
		if err == nil {
			if strings.HasPrefix(absFilePath, absTmplPath) {
				//如果静态文件在模板的子文件夹时，监控模板时判断静态文件更改
				renderer.MonitorEvent(st.OnUpdate(tmplDir))
			}
		}
	}
	return renderer
}

func (t *Config) ApplyTo(e *echo.Echo, manager ...driver.Manager) *Config {
	if t.renderer != nil {
		t.renderer.Close()
	}
	e.SetHTTPErrorHandler(HTTPErrorHandler(t.ErrorPages))
	e.Use(middleware.FuncMap(tplfunc.New(), func(c echo.Context) bool {
		return c.Format() != `html`
	}))
	renderer := t.NewRenderer(manager...)
	if t.StaticOptions != nil {
		e.Use(middleware.Static(t.StaticOptions))
	}
	e.SetRenderer(renderer)
	t.renderer = renderer
	return t
}

func (t *Config) Renderer() driver.Driver {
	return t.renderer
}

func (t *Config) NewStatic() *resource.Static {
	return resource.NewStatic(t.StaticOptions.Path, t.StaticOptions.Root)
}

// ThemeDir 主题所在文件夹的路径
func (t *Config) ThemeDir(args ...string) string {
	if len(args) < 1 {
		return filepath.Join(t.TmplDir, t.Theme)
	}
	return filepath.Join(t.TmplDir, args[0])
}
