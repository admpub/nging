package middleware

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/admpub/log"

	"github.com/webx-top/echo"
)

var ListDirTemplate = `<!doctype html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>{{.file}}</title>
        <meta content="IE=edge,chrome=1" http-equiv="X-UA-Compatible" />
        <meta content="width=device-width,initial-scale=1.0,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no" name="viewport" />
        <link href="/favicon.ico" rel="shortcut icon">
    </head>
    <body>
		<ul id="fileList">
		{{range $k, $d := .dirs}}
		<li><a href="{{$d.Name}}{{if $d.IsDir}}/{{end}}" style="color: {{if $d.IsDir}}#e91e63{{else}}#212121{{end}};">{{$d.Name}}{{if $d.IsDir}}/{{end}}</a></li>
		{{end}}
		</ul>
	</body>
</html>`

type (
	StaticOptions struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper `json:"-"`

		Path     string          `json:"path"` //UrlPath
		Root     string          `json:"root"`
		Fallback []string        `json:"fallback"`
		Index    string          `json:"index"`
		Browse   bool            `json:"browse"`
		Template string          `json:"template"`
		Debug    bool            `json:"debug"`
		FS       http.FileSystem `json:"-"`

		open   func(string) (http.File, error)
		render func(echo.Context, interface{}) error
	}
)

func Static(options ...*StaticOptions) echo.MiddlewareFunc {
	// Default options
	opts := new(StaticOptions)
	if len(options) > 0 {
		opts = options[0]
	}
	return opts.Init().Middleware()
}

func (s *StaticOptions) Init() *StaticOptions {
	if s.Skipper == nil {
		s.Skipper = echo.DefaultSkipper
	}
	var err error
	s.Root, err = filepath.Abs(s.Root)
	if err != nil {
		panic(err)
	}
	for index, fallback := range s.Fallback {
		s.Fallback[index], err = filepath.Abs(fallback)
		if err != nil {
			panic(err)
		}
		if s.Debug {
			log.GetLogger("echo").Debug(`[middleware][static] `, `Register assets directory: `, fallback)
		}
	}
	if len(s.Path) > 0 && s.Path[0] != '/' {
		s.Path = `/` + s.Path
	}
	if s.Debug {
		log.GetLogger("echo").Debug(`[middleware][static] `, `Static: `, s.Path, "\t-> ", s.Root)
	}
	return s
}

func (s *StaticOptions) AddFallback(fallback string) *StaticOptions {
	var err error
	fallback, err = filepath.Abs(fallback)
	if err != nil {
		panic(err)
	}
	s.Fallback = append(s.Fallback, fallback)
	if s.Debug {
		log.GetLogger("echo").Debug(`[middleware][static] `, `Register assets directory: `, fallback)
	}
	return s
}

func (s *StaticOptions) getOpener() func(file string) (http.File, error) {
	if s.open != nil {
		return s.open
	}
	if s.FS != nil {
		s.open = s.FS.Open
	} else {
		s.open = func(name string) (http.File, error) {
			fp, err := os.Open(name)
			return fp, err
		}
	}
	return s.open
}

func (s *StaticOptions) getRender() func(c echo.Context, data interface{}) error {
	if s.render != nil {
		return s.render
	}
	if !s.Browse {
		return nil
	}
	if len(s.Template) > 0 {
		s.render = func(c echo.Context, data interface{}) error {
			return c.Render(s.Template, data)
		}
	} else {
		t := template.New(s.Template)
		_, err := t.Parse(ListDirTemplate)
		if err != nil {
			panic(err)
		}
		s.render = func(c echo.Context, data interface{}) error {
			w := c.Response()
			w.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
			return t.Execute(w, data)
		}
	}
	return s.render
}

func (s *StaticOptions) findFile(c echo.Context, root string, hasIndex bool, file string, render func(echo.Context, interface{}) error, opener func(string) (http.File, error)) error {
	absFile := filepath.Join(root, file)
	fp, err := opener(absFile)
	if err != nil {
		return echo.ErrNotFound
	}
	defer fp.Close()
	fi, err := fp.Stat()
	if err != nil {
		return echo.ErrNotFound
	}
	if fi.IsDir() {
		if hasIndex {
			fp.Close()
			// Index file
			indexFile := filepath.Join(absFile, s.Index)
			fp, err = opener(indexFile)
			if err != nil {
				return echo.ErrNotFound
			}
			fi, err = fp.Stat()
			if err != nil || fi.IsDir() {
				if s.Browse {
					return listDirByCustomFS(absFile, file, c, render, opener)
				}
				return echo.ErrNotFound
			}
			absFile = indexFile
		} else {
			if s.Browse {
				return listDirByCustomFS(absFile, file, c, render, opener)
			}
			return echo.ErrNotFound
		}
	}
	return c.ServeContent(fp, fi.Name(), fi.ModTime())
}

func (s *StaticOptions) Middleware() echo.MiddlewareFunc {
	render := s.getRender()
	opener := s.getOpener()
	hasIndex := len(s.Index) > 0
	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if s.Skipper(c) {
				return next.Handle(c)
			}
			file := c.Request().URL().Path()
			length := len(s.Path)
			if len(file) < length || file[0:length] != s.Path {
				return next.Handle(c)
			}
			file = file[length:]
			file = path.Clean(file)
			err := s.findFile(c, s.Root, hasIndex, file, render, opener)
			if err == nil {
				return err
			}
			if err == echo.ErrNotFound {
				if len(s.Fallback) == 0 {
					return err
				}
				for _, fallback := range s.Fallback {
					if s.Debug {
						log.GetLogger("echo").Debug(`[middleware][static] `, `fallback ->  `, filepath.Join(fallback, file))
					}
					err = s.findFile(c, fallback, hasIndex, file, render, opener)
					if err == nil {
						return err
					}
				}
			}
			return err
		})
	}
}

func listDirByCustomFS(absFile string, file string, c echo.Context, render func(echo.Context, interface{}) error, opener func(string) (http.File, error)) error {
	d, err := opener(absFile)
	if err != nil {
		return echo.ErrNotFound
	}
	defer d.Close()
	dirs, err := d.Readdir(-1)
	if err != nil {
		return echo.ErrNotFound
	}

	return render(c, map[string]interface{}{
		`file`: file,
		`dirs`: dirs,
	})
}
