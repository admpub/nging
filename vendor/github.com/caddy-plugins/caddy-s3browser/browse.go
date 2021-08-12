package s3browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/admpub/caddy/caddyhttp/httpserver"
	"github.com/minio/minio-go/v7"
)

type Browse struct {
	Next     httpserver.Handler
	Config   Config
	Client   *minio.Client
	Fs       map[string]Directory
	Template *template.Template
}

const (
	AssetsPrefix    = `/---caddy-s3browser---/`
	MIMEOctetStream = "application/octet-stream"
)

func (b Browse) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	path := r.URL.Path
	if len(path) == 0 {
		path = "/"
	} else if len(b.Config.CSSCDN) == 0 && strings.HasPrefix(path, AssetsPrefix) {
		path = strings.TrimPrefix(path, AssetsPrefix)
		fp, err := assets.Open(path)
		if err != nil {
			return http.StatusNotFound, nil
		}
		defer fp.Close()
		fi, err := fp.Stat()
		if err != nil {
			return http.StatusNotFound, nil
		}
		if fi.IsDir() {
			return http.StatusNotFound, nil
		}
		if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && fi.ModTime().Before(t.Add(1*time.Second)) {
			w.Header().Del("Content-Type")
			w.Header().Del("Content-Length")
			w.WriteHeader(http.StatusNotModified)
			return http.StatusNotModified, nil
		}
		w.Header().Set("Content-Type", ContentTypeByExtension(fi.Name()))
		w.Header().Set("Last-Modified", fi.ModTime().UTC().Format(http.TimeFormat))
		n, err := io.Copy(w, fp)
		if err != nil {
			return http.StatusInternalServerError, err
		}
		w.Header().Set("Content-Length", fmt.Sprintf(`%d`, n))
		return http.StatusOK, err
	}
	if _, ok := b.Fs[path]; !ok {
		if !strings.HasSuffix(path, `/`) { // 访问的是文件
			if len(b.Config.CDNURL) > 0 { // 如果指定了CDN的网址，则跳转到CDN网址
				endpoint := b.Config.CDNURL + path
				http.Redirect(w, r, endpoint, http.StatusFound)
				return http.StatusFound, nil
			}
			return b.Next.ServeHTTP(w, r)
		}
		// 访问未登记的目录返回 not found
		return http.StatusNotFound, nil
	}
	switch r.Method {
	case http.MethodGet, http.MethodHead:
		// proceed, noop
	case "PROPFIND", http.MethodOptions:
		return http.StatusNotImplemented, nil
	default:
		return b.Next.ServeHTTP(w, r)
	}

	var buf *bytes.Buffer
	var err error
	acceptHeader := strings.ToLower(strings.Join(r.Header["Accept"], ","))
	switch {
	case strings.Contains(acceptHeader, "application/json"):
		if buf, err = b.formatAsJSON(b.Fs[path]); err != nil {
			return http.StatusInternalServerError, err
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	default:
		if buf, err = b.formatAsHTML(b.Fs[path]); err != nil {
			if b.Config.Debug {
				fmt.Println(err)
			}
			return http.StatusInternalServerError, err
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	buf.WriteTo(w)

	return http.StatusOK, nil
}

func (b Browse) formatAsJSON(listing Directory) (*bytes.Buffer, error) {
	data := TmplData{CDNURL: b.Config.CDNURL, Directory: listing}
	marsh, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.Write(marsh)
	return buf, err
}

func (b Browse) formatAsHTML(listing Directory) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	data := TmplData{CDNURL: b.Config.CDNURL, Directory: listing}
	err := b.Template.Execute(buf, data)
	return buf, err
}
