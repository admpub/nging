/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package echo

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/webx-top/echo/param"
)

var (
	splitHTTPMethod = regexp.MustCompile(`[^A-Z]+`)

	methods = []string{
		CONNECT,
		DELETE,
		GET,
		HEAD,
		OPTIONS,
		PATCH,
		POST,
		PUT,
		TRACE,
	}

	//--------
	// Errors
	//--------

	ErrUnsupportedMediaType        error = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrNotFound                    error = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                error = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden                   error = NewHTTPError(http.StatusForbidden)
	ErrStatusRequestEntityTooLarge error = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrMethodNotAllowed            error = NewHTTPError(http.StatusMethodNotAllowed)
	ErrRendererNotRegistered             = errors.New("renderer not registered")
	ErrInvalidRedirectCode               = errors.New("invalid redirect status code")
	ErrNotFoundFileInput                 = errors.New("The specified name file input was not found")

	//----------------
	// Error handlers
	//----------------

	NotFoundHandler = HandlerFunc(func(c Context) error {
		return ErrNotFound
	})

	MethodNotAllowedHandler = HandlerFunc(func(c Context) error {
		return ErrMethodNotAllowed
	})

	_ MiddlewareFuncd = func(h Handler) HandlerFunc {
		return func(c Context) error {
			return h.Handle(c)
		}
	}

	globalVars = sync.Map{} //Custom global variable
)

func Set(key, value interface{}) {
	globalVars.Store(key, value)
}

func Get(key interface{}, defaults ...interface{}) interface{} {
	value, ok := globalVars.Load(key)
	if !ok && len(defaults) > 0 {
		if fallback, ok := defaults[0].(func() interface{}); ok {
			return fallback()
		}
		return defaults[0]
	}
	return value
}

func GetOk(key interface{}) (interface{}, bool) {
	return globalVars.Load(key)
}

func Has(key interface{}) bool {
	_, ok := globalVars.Load(key)
	return ok
}

func Delete(key interface{}) {
	globalVars.Delete(key)
}

func Range(f func(key, value interface{}) bool) {
	globalVars.Range(f)
}

func GetOrSet(key, value interface{}) (actual interface{}, loaded bool) {
	return globalVars.LoadOrStore(key, value)
}

func String(key interface{}, defaults ...interface{}) string {
	return param.AsString(Get(key, defaults...))
}

func Split(key interface{}, sep string, limit ...int) param.StringSlice {
	return param.Split(Get(key), sep, limit...)
}

func Trim(key interface{}, defaults ...interface{}) param.String {
	return param.Trim(Get(key, defaults...))
}

func HTML(key interface{}, defaults ...interface{}) template.HTML {
	return param.AsHTML(Get(key, defaults...))
}

func HTMLAttr(key interface{}, defaults ...interface{}) template.HTMLAttr {
	return param.AsHTMLAttr(Get(key, defaults...))
}

func JS(key interface{}, defaults ...interface{}) template.JS {
	return param.AsJS(Get(key, defaults...))
}

func CSS(key interface{}, defaults ...interface{}) template.CSS {
	return param.AsCSS(Get(key, defaults...))
}

func Bool(key interface{}, defaults ...interface{}) bool {
	return param.AsBool(Get(key, defaults...))
}

func Float64(key interface{}, defaults ...interface{}) float64 {
	return param.AsFloat64(Get(key, defaults...))
}

func Float32(key interface{}, defaults ...interface{}) float32 {
	return param.AsFloat32(Get(key, defaults...))
}

func Int8(key interface{}, defaults ...interface{}) int8 {
	return param.AsInt8(Get(key, defaults...))
}

func Int16(key interface{}, defaults ...interface{}) int16 {
	return param.AsInt16(Get(key, defaults...))
}

func Int(key interface{}, defaults ...interface{}) int {
	return param.AsInt(Get(key, defaults...))
}

func Int32(key interface{}, defaults ...interface{}) int32 {
	return param.AsInt32(Get(key, defaults...))
}

func Int64(key interface{}, defaults ...interface{}) int64 {
	return param.AsInt64(Get(key, defaults...))
}

func Decr(key interface{}, n int64, defaults ...interface{}) int64 {
	v := param.Decr(Get(key, defaults...), n)
	Set(key, v)
	return v
}

func Incr(key interface{}, n int64, defaults ...interface{}) int64 {
	v := param.Incr(Get(key, defaults...), n)
	Set(key, v)
	return v
}

func Uint8(key interface{}, defaults ...interface{}) uint8 {
	return param.AsUint8(Get(key, defaults...))
}

func Uint16(key interface{}, defaults ...interface{}) uint16 {
	return param.AsUint16(Get(key, defaults...))
}

func Uint(key interface{}, defaults ...interface{}) uint {
	return param.AsUint(Get(key, defaults...))
}

func Uint32(key interface{}, defaults ...interface{}) uint32 {
	return param.AsUint32(Get(key, defaults...))
}

func Uint64(key interface{}, defaults ...interface{}) uint64 {
	return param.AsUint64(Get(key, defaults...))
}

func Timestamp(key interface{}, defaults ...interface{}) time.Time {
	return param.AsTimestamp(Get(key, defaults...))
}

func DateTime(key interface{}, layouts ...string) time.Time {
	return param.AsDateTime(Get(key), layouts...)
}

func Children(key interface{}, keys ...interface{}) Store {
	r := GetStore(key)
	for _, key := range keys {
		r = GetStore(key)
	}
	return r
}

func GetStore(key interface{}, defaults ...interface{}) Store {
	return AsStore(Get(key, defaults...))
}
