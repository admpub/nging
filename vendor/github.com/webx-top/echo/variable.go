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

	ErrUnsupportedMediaType         error = NewHTTPError(http.StatusUnsupportedMediaType)
	ErrBadRequest                   error = NewHTTPError(http.StatusBadRequest)
	ErrPaymentRequired              error = NewHTTPError(http.StatusPaymentRequired)
	ErrNotAcceptable                error = NewHTTPError(http.StatusNotAcceptable)
	ErrProxyAuthRequired            error = NewHTTPError(http.StatusProxyAuthRequired)
	ErrRequestTimeout               error = NewHTTPError(http.StatusRequestTimeout)
	ErrConflict                     error = NewHTTPError(http.StatusConflict)
	ErrGone                         error = NewHTTPError(http.StatusGone)
	ErrLengthRequired               error = NewHTTPError(http.StatusLengthRequired)
	ErrPreconditionFailed           error = NewHTTPError(http.StatusPreconditionFailed)
	ErrRequestEntityTooLarge        error = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrRequestURITooLong            error = NewHTTPError(http.StatusRequestURITooLong)
	ErrRequestedRangeNotSatisfiable error = NewHTTPError(http.StatusRequestedRangeNotSatisfiable)
	ErrExpectationFailed            error = NewHTTPError(http.StatusExpectationFailed)
	ErrUnprocessableEntity          error = NewHTTPError(http.StatusUnprocessableEntity)
	ErrLocked                       error = NewHTTPError(http.StatusLocked)
	ErrFailedDependency             error = NewHTTPError(http.StatusFailedDependency)
	ErrTooEarly                     error = NewHTTPError(http.StatusTooEarly)
	ErrUpgradeRequired              error = NewHTTPError(http.StatusUpgradeRequired)
	ErrPreconditionRequired         error = NewHTTPError(http.StatusPreconditionRequired)
	ErrTooManyRequests              error = NewHTTPError(http.StatusTooManyRequests)
	ErrRequestHeaderFieldsTooLarge  error = NewHTTPError(http.StatusRequestHeaderFieldsTooLarge)
	ErrUnavailableForLegalReasons   error = NewHTTPError(http.StatusUnavailableForLegalReasons)
	ErrNotImplemented               error = NewHTTPError(http.StatusNotImplemented)
	ErrNotFound                     error = NewHTTPError(http.StatusNotFound)
	ErrUnauthorized                 error = NewHTTPError(http.StatusUnauthorized)
	ErrForbidden                    error = NewHTTPError(http.StatusForbidden)
	ErrStatusRequestEntityTooLarge  error = NewHTTPError(http.StatusRequestEntityTooLarge)
	ErrMethodNotAllowed             error = NewHTTPError(http.StatusMethodNotAllowed)
	ErrRendererNotRegistered              = errors.New("renderer not registered")
	ErrInvalidRedirectCode                = errors.New("invalid redirect status code")
	ErrNotFoundFileInput                  = errors.New("the specified name file input was not found")

	//----------------
	// Error handlers
	//----------------

	NotFoundHandler = HandlerFunc(func(c Context) error {
		return ErrNotFound
	})

	ErrorHandler = func(err error) Handler {
		return HandlerFunc(func(c Context) error {
			return err
		})
	}

	MethodNotAllowedHandler = HandlerFunc(func(c Context) error {
		return ErrMethodNotAllowed
	})

	_ MiddlewareFuncd = func(h Handler) HandlerFunc {
		return func(c Context) error {
			return h.Handle(c)
		}
	}

	//----------------
	// Shortcut
	//----------------

	StringerMapStart = param.StringerMapStart
	StoreStart       = param.StoreStart
	HStart           = param.StoreStart

	//Custom global variable
	globalVars = param.NewMap()
)

func Set(key, value interface{}) {
	globalVars.Set(key, value)
}

func Get(key interface{}, defaults ...interface{}) interface{} {
	return globalVars.Get(key, defaults...)
}

func GetStoreByKeys(key interface{}, keys ...string) H {
	st, ok := Get(key).(H)
	if !ok {
		if st == nil {
			st = H{}
		}
		return st
	}
	return st.GetStoreByKeys(keys...)
}

func GetOk(key interface{}) (interface{}, bool) {
	return globalVars.GetOk(key)
}

func Has(key interface{}) bool {
	return globalVars.Has(key)
}

func Delete(key interface{}) {
	globalVars.Delete(key)
}

func Range(f func(key, value interface{}) bool) {
	globalVars.Range(f)
}

func GetOrSet(key, value interface{}) (actual interface{}, loaded bool) {
	return globalVars.GetOrSet(key, value)
}

func String(key interface{}, defaults ...interface{}) string {
	return globalVars.String(key, defaults...)
}

func Split(key interface{}, sep string, limit ...int) param.StringSlice {
	return globalVars.Split(key, sep, limit...)
}

func Trim(key interface{}, defaults ...interface{}) param.String {
	return globalVars.Trim(key, defaults...)
}

func HTML(key interface{}, defaults ...interface{}) template.HTML {
	return globalVars.HTML(key, defaults...)
}

func HTMLAttr(key interface{}, defaults ...interface{}) template.HTMLAttr {
	return globalVars.HTMLAttr(key, defaults...)
}

func JS(key interface{}, defaults ...interface{}) template.JS {
	return globalVars.JS(key, defaults...)
}

func CSS(key interface{}, defaults ...interface{}) template.CSS {
	return globalVars.CSS(key, defaults...)
}

func Bool(key interface{}, defaults ...interface{}) bool {
	return globalVars.Bool(key, defaults...)
}

func Float64(key interface{}, defaults ...interface{}) float64 {
	return globalVars.Float64(key, defaults...)
}

func Float32(key interface{}, defaults ...interface{}) float32 {
	return globalVars.Float32(key, defaults...)
}

func Int8(key interface{}, defaults ...interface{}) int8 {
	return globalVars.Int8(key, defaults...)
}

func Int16(key interface{}, defaults ...interface{}) int16 {
	return globalVars.Int16(key, defaults...)
}

func Int(key interface{}, defaults ...interface{}) int {
	return globalVars.Int(key, defaults...)
}

func Int32(key interface{}, defaults ...interface{}) int32 {
	return globalVars.Int32(key, defaults...)
}

func Int64(key interface{}, defaults ...interface{}) int64 {
	return globalVars.Int64(key, defaults...)
}

func Decr(key interface{}, n int64, defaults ...interface{}) int64 {
	return globalVars.Decr(key, n, defaults...)
}

func Incr(key interface{}, n int64, defaults ...interface{}) int64 {
	return globalVars.Incr(key, n, defaults...)
}

func Uint8(key interface{}, defaults ...interface{}) uint8 {
	return globalVars.Uint8(key, defaults...)
}

func Uint16(key interface{}, defaults ...interface{}) uint16 {
	return globalVars.Uint16(key, defaults...)
}

func Uint(key interface{}, defaults ...interface{}) uint {
	return globalVars.Uint(key, defaults...)
}

func Uint32(key interface{}, defaults ...interface{}) uint32 {
	return globalVars.Uint32(key, defaults...)
}

func Uint64(key interface{}, defaults ...interface{}) uint64 {
	return globalVars.Uint64(key, defaults...)
}

func Timestamp(key interface{}, defaults ...interface{}) time.Time {
	return globalVars.Timestamp(key, defaults...)
}

func DateTime(key interface{}, layouts ...string) time.Time {
	return globalVars.DateTime(key, layouts...)
}

func Children(key interface{}, keys ...interface{}) Store {
	r := GetStore(key)
	for _, key := range keys {
		r = GetStore(key)
	}
	return r
}

func GetStore(key interface{}, defaults ...interface{}) Store {
	return AsStore(globalVars.Get(key, defaults...))
}
