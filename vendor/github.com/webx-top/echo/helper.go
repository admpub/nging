package echo

import (
	"mime"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// HandlerName returns the handler name
func HandlerName(h interface{}) string {
	v := reflect.ValueOf(h)
	t := v.Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(v.Pointer()).Name()
	}
	return t.String()
}

// Methods returns methods
func Methods() []string {
	return methods
}

// ContentTypeByExtension returns the MIME type associated with the file based on
// its extension. It returns `application/octet-stream` incase MIME type is not
// found.
func ContentTypeByExtension(name string) (t string) {
	if t = mime.TypeByExtension(filepath.Ext(name)); len(t) == 0 {
		t = MIMEOctetStream
	}
	return
}

func static(r RouteRegister, prefix, root string) {
	var err error
	root, err = filepath.Abs(root)
	if err != nil {
		panic(err)
	}
	h := func(c Context) error {
		name := filepath.Join(root, c.Param("*"))
		if !strings.HasPrefix(name, root) {
			return ErrNotFound
		}
		return c.File(name)
	}
	if prefix == "/" {
		r.Get(prefix+"*", h)
	} else {
		r.Get(prefix+"/*", h)
	}
}
