package markdown

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	md2html "github.com/russross/blackfriday"
	"github.com/webx-top/echo"
)

type (
	Options struct {
		Path               string   `json:"path"` //UrlPath
		MarkdownExtensions []string `json:"markdownExtensions"`
		Index              string   `json:"index"`
		Root               string   `json:"root"`
		Browse             bool     `json:"browse"`
		Preprocessor       func(echo.Context, []byte) []byte
		Filter             func(string) bool // true: ok; false: ignore
	}
)

func Markdown(options ...*Options) echo.MiddlewareFunc {
	// Default options
	opts := new(Options)
	if len(options) > 0 {
		opts = options[0]
	}
	if opts.Index == "" {
		opts.Index = "SUMMARY.md"
	}
	if opts.MarkdownExtensions == nil {
		opts.MarkdownExtensions = []string{`.md`, `.mdown`, `.markdown`}
	}
	opts.Root, _ = filepath.Abs(opts.Root)

	if opts.Preprocessor == nil {
		opts.Preprocessor = func(c echo.Context, b []byte) []byte {
			return b
		}
	}
	if opts.Filter == nil {
		opts.Filter = func(string) bool {
			return true
		}
	}

	length := len(opts.Path)

	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			file := c.Request().URL().Path()
			if len(file) < length || file[0:length] != opts.Path {
				return next.Handle(c)
			}
			if !opts.Filter(file) {
				return next.Handle(c)
			}
			file = filepath.Clean(file[length:])
			absFile := filepath.Join(opts.Root, file)
			if !strings.HasPrefix(absFile, opts.Root) {
				return next.Handle(c)
			}
			fi, err := os.Stat(absFile)
			if err != nil {
				return echo.ErrNotFound
			}
			w := c.Response()
			if fi.IsDir() {
				// Index file
				indexFile := filepath.Join(absFile, opts.Index)
				fi, err = os.Stat(indexFile)
				if err != nil || fi.IsDir() {
					if opts.Browse {
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
							if !opts.Filter(name) {
								continue
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
					return echo.ErrNotFound
				}
				absFile = indexFile
			}
			ext := strings.ToLower(filepath.Ext(fi.Name()))
			isMarkdownDocument := false
			for _, vext := range opts.MarkdownExtensions {
				if ext == vext {
					isMarkdownDocument = true
					break
				}
			}
			if isMarkdownDocument {
				modtime := fi.ModTime()
				if t, err := time.Parse(http.TimeFormat, c.Request().Header().Get(echo.HeaderIfModifiedSince)); err == nil && modtime.Before(t.Add(1*time.Second)) {
					w.Header().Del(echo.HeaderContentType)
					w.Header().Del(echo.HeaderContentLength)
					return c.NoContent(http.StatusNotModified)
				}
				var b []byte
				b, err = ioutil.ReadFile(absFile)
				if err != nil {
					return echo.ErrNotFound
				}
				b = opts.Preprocessor(c, b)
				b = md2html.MarkdownCommon(b)

				w.Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
				w.Header().Set(echo.HeaderLastModified, modtime.UTC().Format(http.TimeFormat))
				w.WriteHeader(http.StatusOK)
				_, err = w.Write(b)
			} else {
				w.Header().Set(echo.HeaderContentType, echo.ContentTypeByExtension(ext))
				w.ServeFile(absFile)
			}
			return err
		})
	}
}
