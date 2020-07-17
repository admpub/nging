package goseaweedfs

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	workerpool "github.com/linxGnu/gumble/worker-pool"
)

func createWorkerPool() *workerpool.Pool {
	return workerpool.NewPool(context.Background(), workerpool.Option{
		NumberWorker: runtime.NumCPU() << 1,
	})
}

func parseURI(uri string) (u *url.URL, err error) {
	u, err = url.Parse(uri)
	if err == nil && u.Scheme == "" {
		u.Scheme = "http"
	}
	return
}

func encodeURI(base url.URL, path string, args url.Values) string {
	base.Path = path
	query := base.Query()
	args = normalize(args, "", "")
	for k, vs := range args {
		for _, v := range vs {
			query.Add(k, v)
		}
	}
	base.RawQuery = query.Encode()
	return base.String()
}

func valid(c rune) bool {
	return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') || '.' == c || '-' == c || '_' == c
}

func normalizeName(st string) string {
	for _, _c := range st {
		if !valid(_c) {
			var sb strings.Builder
			sb.Grow(len(st))

			for _, c := range st {
				if valid(c) {
					_, _ = sb.WriteRune(c)
				}
			}

			return sb.String()
		}
	}
	return st
}

func drainAndClose(body io.ReadCloser) {
	_, _ = io.Copy(ioutil.Discard, body)
	_ = body.Close()
}

func normalize(values url.Values, collection, ttl string) url.Values {
	if values == nil {
		values = make(url.Values)
	}

	if len(collection) > 0 {
		values.Set(ParamCollection, collection)
	}

	if len(ttl) > 0 {
		values.Set(ParamTTL, ttl)
	}

	return values
}

func readAll(r *http.Response) (body []byte, statusCode int, err error) {
	statusCode = r.StatusCode
	body, err = ioutil.ReadAll(r.Body)
	r.Body.Close()
	return
}
