package s3browser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"net/http"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"github.com/minio/minio-go/v6"
)

type Browse struct {
	Next     httpserver.Handler
	Config   Config
	Client   *minio.Client
	Fs       map[string]Directory
	Template *template.Template
}

func (b Browse) ServeHTTP(w http.ResponseWriter, r *http.Request) (int, error) {

	path := r.URL.Path
	if path == "" {
		path = "/"
	}
	if _, ok := b.Fs[path]; !ok {
		return b.Next.ServeHTTP(w, r)
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
			fmt.Println(err)
			return http.StatusInternalServerError, err
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	buf.WriteTo(w)

	return http.StatusOK, nil
}

func (b Browse) formatAsJSON(listing Directory) (*bytes.Buffer, error) {
	marsh, err := json.Marshal(listing)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.Write(marsh)
	return buf, err
}

func (b Browse) formatAsHTML(listing Directory) (*bytes.Buffer, error) {
	buf := new(bytes.Buffer)
	err := b.Template.Execute(buf, listing)
	return buf, err
}
