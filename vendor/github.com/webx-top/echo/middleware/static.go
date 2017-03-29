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
	"github.com/webx-top/echo/engine"
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

		Path     string `json:"path"` //UrlPath
		Root     string `json:"root"`
		Index    string `json:"index"`
		Browse   bool   `json:"browse"`
		Template string `json:"template"`
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
	length := len(opts.Path)
	if length > 0 && opts.Path[0] != '/' {
		opts.Path = `/` + opts.Path
		length++
	}
	var t *template.Template
	if opts.Browse {
		t = template.New(opts.Template)
		var e error
		if len(opts.Template) > 0 {
			t, e = t.ParseFiles(opts.Template)
		} else {
			t, e = t.Parse(ListDirTemplate)
		}
		if e != nil {
			panic(e)
		}
	}

	log.GetLogger("echo").Debugf("Static: %v\t-> %v", opts.Path, opts.Root)

	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if opts.Skipper(c) {
				return next.Handle(c)
			}
			file := c.Request().URL().Path()
			if len(file) < length || file[0:length] != opts.Path {
				return next.Handle(c)
			}
			file = file[length:]
			file = path.Clean(file)
			absFile := filepath.Join(opts.Root, file)
			if !strings.HasPrefix(absFile, opts.Root) {
				return next.Handle(c)
			}
			w := c.Response()
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
							return listDir(absFile, file, w, t)
						}
						return echo.ErrNotFound
					}
					absFile = indexFile
				} else {
					if opts.Browse {
						return listDir(absFile, file, w, t)
					}
					return echo.ErrNotFound
				}
			}
			return c.ServeContent(fp, fi.Name(), fi.ModTime())
		})
	}
}

func listDir(absFile string, file string, w engine.Response, t *template.Template) error {
	fs := http.Dir(filepath.Dir(absFile))
	d, err := fs.Open(filepath.Base(absFile))
	if err != nil {
		return echo.ErrNotFound
	}
	defer d.Close()
	dirs, err := d.Readdir(-1)
	if err != nil {
		return echo.ErrNotFound
	}

	w.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return t.Execute(w, map[string]interface{}{
		`file`: file,
		`dirs`: dirs,
	})
}
