package middleware

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

type (
	StaticOptions struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper `json:"-"`

		Path   string `json:"path"` //UrlPath
		Root   string `json:"root"`
		Index  string `json:"index"`
		Browse bool   `json:"browse"`
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
	check := hasIndex || opts.Browse
	opts.Root, _ = filepath.Abs(opts.Root)
	length := len(opts.Path)
	if length > 0 && opts.Path[0] != '/' {
		opts.Path = `/` + opts.Path
		length++
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
			if !check {
				w.ServeFile(absFile)
				return nil
			}
			fi, err := os.Stat(absFile)
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
							return listDir(absFile, file, w)
						}
						return echo.ErrNotFound
					}
					absFile = indexFile
				} else {
					if opts.Browse {
						return listDir(absFile, file, w)
					}
					return echo.ErrNotFound
				}
			}
			w.ServeFile(absFile)
			return nil
		})
	}
}

func listDir(absFile string, file string, w engine.Response) error {
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

	// Create a directory index
	w.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	if _, err = fmt.Fprintf(w, `<!doctype html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>`+file+`</title>
        <meta content="IE=edge,chrome=1" http-equiv="X-UA-Compatible" />
        <meta content="width=device-width,initial-scale=1.0,minimum-scale=1.0,maximum-scale=1.0,user-scalable=no" name="viewport" />
        <link href="/favicon.ico" rel="shortcut icon">
    </head>
    <body>`); err != nil {
		return err
	}
	if _, err = fmt.Fprintf(w, "<ul id=\"fileList\">\n"); err != nil {
		return err
	}
	for _, d := range dirs {
		name := d.Name()
		color := "#212121"
		if d.IsDir() {
			color = "#e91e63"
			name += "/"
		}
		if _, err = fmt.Fprintf(w, "<li><a href=\"%s\" style=\"color: %s;\">%s</a></li>\n", name, color, name); err != nil {
			return err
		}
	}
	if _, err = fmt.Fprintf(w, "</ul>\n"); err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "</body>\n</html>")
	return err
}
