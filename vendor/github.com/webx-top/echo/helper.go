package echo

import (
	"fmt"
	"mime"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/encoding/json"
)

var workDir string

func SetWorkDir(dir string) {
	if len(dir) == 0 {
		if len(workDir) == 0 {
			setWorkDir()
		}
		return
	}
	if !strings.HasSuffix(dir, FilePathSeparator) {
		dir += FilePathSeparator
	}
	workDir = dir
}

func setWorkDir() {
	workDir, _ = os.Getwd()
	workDir = workDir + FilePathSeparator
}

func init() {
	if len(workDir) == 0 {
		setWorkDir()
	}
}

func Wd() string {
	if len(workDir) == 0 {
		setWorkDir()
	}
	return workDir
}

// HandlerName returns the handler name
func HandlerName(h interface{}) string {
	if h == nil {
		return `<nil>`
	}
	v := reflect.ValueOf(h)
	t := v.Type()
	if t.Kind() == reflect.Func {
		return runtime.FuncForPC(v.Pointer()).Name()
	}
	return t.String()
}

// HandlerPath returns the handler path
func HandlerPath(h interface{}) string {
	v := reflect.ValueOf(h)
	t := v.Type()
	switch t.Kind() {
	case reflect.Func:
		return runtime.FuncForPC(v.Pointer()).Name()
	case reflect.Ptr:
		t = t.Elem()
		fallthrough
	case reflect.Struct:
		return t.PkgPath() + `.` + t.Name()
	}
	return ``
}

func HandlerTmpl(handlerPath string) string {
	name := path.Base(handlerPath)
	var r []string
	var u []rune
	for _, b := range name {
		switch b {
		case '*', '(', ')':
			continue
		case '-':
			goto END
		case '.':
			r = append(r, string(u))
			u = []rune{}
		default:
			u = append(u, b)
		}
	}

END:
	if len(u) > 0 {
		r = append(r, string(u))
		u = []rune{}
	}
	for i, s := range r {
		r[i] = com.SnakeCase(s)
	}
	return `/` + strings.Join(r, `/`)
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

func Clear(old []interface{}, clears ...interface{}) []interface{} {
	if len(clears) == 0 {
		return nil
	}
	if len(old) == 0 {
		return old
	}
	result := []interface{}{}
	for _, el := range old {
		var exists bool
		for _, d := range clears {
			if d == el {
				exists = true
				break
			}
		}
		if !exists {
			result = append(result, el)
		}
	}
	return result
}

// Dump 输出对象和数组的结构信息
func Dump(m interface{}, printOrNot ...bool) (r string) {
	v, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	r = string(v)
	l := len(printOrNot)
	if l < 1 || printOrNot[0] {
		fmt.Println(r)
	}
	return
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func LogIf(err error, types ...string) {
	if err == nil {
		return
	}
	var typ string
	if len(types) > 0 {
		typ = types[0]
	}
	typ = strings.Title(typ)
	switch typ {
	case `Fatal`:
		log.Fatal(err)
	case `Warn`:
		log.Debug(err)
	case `Debug`:
		log.Debug(err)
	case `Info`:
		log.Info(err)
	default:
		log.Error(err)
	}
}

func URLEncode(s string, rfc ...bool) string {
	encoded := url.QueryEscape(s)
	if len(rfc) > 0 && rfc[0] { // RFC 3986
		encoded = strings.Replace(encoded, `+`, `%20`, -1)
	}
	return encoded
}

func URLDecode(encoded string, rfc ...bool) (string, error) {
	if len(rfc) > 0 && rfc[0] {
		encoded = strings.Replace(encoded, `%20`, `+`, -1)
	}
	return url.QueryUnescape(encoded)
}

func InSliceFold(value string, items []string) bool {
	for _, item := range items {
		if strings.EqualFold(item, value) {
			return true
		}
	}
	return false
}

type HandlerFuncs map[string]func(Context) error

func (h *HandlerFuncs) Register(key string, fn func(Context) error) {
	(*h)[key] = fn
}

func (h *HandlerFuncs) Unregister(keys ...string) {
	for _, key := range keys {
		delete(*h, key)
	}
}

func (h HandlerFuncs) Call(c Context, key string) error {
	fn, ok := h[key]
	if !ok {
		return ErrNotFound
	}
	return fn(c)
}
