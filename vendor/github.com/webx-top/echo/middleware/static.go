package middleware

import (
	"html/template"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

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
		Index    string          `json:"index"`
		Browse   bool            `json:"browse"`
		Template string          `json:"template"`
		Debug    bool            `json:"debug"`
		FS       http.FileSystem `json:"-"`
	}
)

func Static(options ...*StaticOptions) echo.MiddlewareFunc {
	// Default options
	opts := new(StaticOptions)
	if len(options) > 0 {
		opts = options[0]
	}
	hasIndex := len(opts.Index) > 0
	if opts.Skipper == nil {
		opts.Skipper = echo.DefaultSkipper
	}
	opts.Root, _ = filepath.Abs(opts.Root)
	if len(opts.Path) > 0 && opts.Path[0] != '/' {
		opts.Path = `/` + opts.Path
	}
	var render func(echo.Context, interface{}) error
	if opts.Browse {
		if len(opts.Template) > 0 {
			render = func(c echo.Context, data interface{}) error {
				return c.Render(opts.Template, data)
			}
		} else {
			t := template.New(opts.Template)
			_, e := t.Parse(ListDirTemplate)
			if e != nil {
				panic(e)
			}
			render = func(c echo.Context, data interface{}) error {
				w := c.Response()
				w.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
				return t.Execute(w, data)
			}
		}
	}

	if opts.Debug {
		log.GetLogger("echo").Debugf("Static: %v\t-> %v", opts.Path, opts.Root)
	}
	if opts.FS != nil {
		return customFS(opts, hasIndex, render)
	}
	return defaultFS(opts, hasIndex, render)
}

func defaultFS(opts *StaticOptions, hasIndex bool, render func(echo.Context, interface{}) error) echo.MiddlewareFunc {
	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if opts.Skipper(c) {
				return next.Handle(c)
			}
			file := c.Request().URL().Path()
			length := len(opts.Path)
			if len(file) < length || file[0:length] != opts.Path {
				return next.Handle(c)
			}
			file = file[length:]
			file = path.Clean(file)
			absFile := filepath.Join(opts.Root, file)
			if !strings.HasPrefix(absFile, opts.Root) {
				return next.Handle(c)
			}
			fp, err := os.Open(absFile)
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
					// Index file
					indexFile := filepath.Join(absFile, opts.Index)
					fi, err = os.Stat(indexFile)
					if err != nil || fi.IsDir() {
						if opts.Browse {
							return listDir(absFile, file, c, render)
						}
						return echo.ErrNotFound
					}
					absFile = indexFile
				} else {
					if opts.Browse {
						return listDir(absFile, file, c, render)
					}
					return echo.ErrNotFound
				}
			}
			return c.ServeContent(fp, fi.Name(), fi.ModTime())
		})
	}
}

func customFS(opts *StaticOptions, hasIndex bool, render func(echo.Context, interface{}) error) echo.MiddlewareFunc {
	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if opts.Skipper(c) {
				return next.Handle(c)
			}
			file := c.Request().URL().Path()
			length := len(opts.Path)
			if len(file) < length || file[0:length] != opts.Path {
				return next.Handle(c)
			}
			file = file[length:]
			file = path.Clean(file)
			absFile := filepath.Join(opts.Root, file)
			if !strings.HasPrefix(absFile, opts.Root) {
				return next.Handle(c)
			}
			fp, err := opts.FS.Open(file)
			if err != nil {
				return echo.ErrNotFound
			}
			fi, err := fp.Stat()
			if err != nil {
				fp.Close()
				return echo.ErrNotFound
			}
			if fi.IsDir() {
				fp.Close()
				if hasIndex {
					// Index file
					indexFile := filepath.Join(file, opts.Index)
					fp, err = opts.FS.Open(indexFile)
					if err != nil {
						return echo.ErrNotFound
					}
					fi, err = fp.Stat()
					if err != nil || fi.IsDir() {
						fp.Close()
						if opts.Browse {
							return listDirByCustomFS(absFile, file, c, render, opts.FS)
						}
						return echo.ErrNotFound
					}
				} else {
					if opts.Browse {
						return listDirByCustomFS(absFile, file, c, render, opts.FS)
					}
					return echo.ErrNotFound
				}
			}
			defer fp.Close()
			return c.ServeContent(fp, fi.Name(), fi.ModTime())
		})
	}
}

func listDir(absFile string, file string, c echo.Context, render func(echo.Context, interface{}) error) error {
	fs := http.Dir(filepath.Dir(absFile))
	return listDirByCustomFS(absFile, filepath.Base(absFile), c, render, fs)
}

func listDirByCustomFS(absFile string, file string, c echo.Context, render func(echo.Context, interface{}) error, fs http.FileSystem) error {
	d, err := fs.Open(file)
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
