/*

   Copyright 2016-present Wenhui Shen <www.webx.top>

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
	stdJSON "encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/webx-top/echo/encoding/json"
)

var (
	DefaultAcceptFormats = map[string]string{
		//json
		MIMEApplicationJSON:       ContentTypeJSON,
		`text/javascript`:         ContentTypeJSON,
		MIMEApplicationJavaScript: ContentTypeJSON,

		//xml
		MIMEApplicationXML: ContentTypeXML,
		`text/xml`:         ContentTypeXML,

		//text
		MIMETextPlain: ContentTypeText,

		//html
		`*/*`:               ContentTypeHTML,
		`application/xhtml`: ContentTypeHTML,
		MIMETextHTML:        ContentTypeHTML,

		//default
		`*`: ContentTypeHTML,
	}
	DefaultFormatRenderers = map[string]FormatRender{
		ContentTypeJSON: func(c Context, data interface{}) error {
			return c.JSON(c.Data())
		},
		ContentTypeJSONP: func(c Context, data interface{}) error {
			return c.JSONP(c.Query(c.Echo().JSONPVarName), c.Data())
		},
		ContentTypeXML: func(c Context, data interface{}) error {
			return c.XML(c.Data())
		},
		ContentTypeText: func(c Context, data interface{}) error {
			return c.String(fmt.Sprint(data))
		},
	}
	DefaultBinderDecoders = map[string]func(interface{}, Context, ...FormDataFilter) error{
		MIMEApplicationJSON: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
			body := ctx.Request().Body()
			if body == nil {
				return NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
			}
			defer body.Close()
			err := json.NewDecoder(body).Decode(i)
			if err != nil {
				if ute, ok := err.(*stdJSON.UnmarshalTypeError); ok {
					return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetRaw(err)
				}
				if se, ok := err.(*stdJSON.SyntaxError); ok {
					return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetRaw(err)
				}
				return NewHTTPError(http.StatusBadRequest, err.Error()).SetRaw(err)
			}
			return err
		},
		MIMEApplicationXML: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
			body := ctx.Request().Body()
			if body == nil {
				return NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
			}
			defer body.Close()
			err := xml.NewDecoder(body).Decode(i)
			if err != nil {
				if ute, ok := err.(*xml.UnsupportedTypeError); ok {
					return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unsupported type error: type=%v, error=%v", ute.Type, ute.Error())).SetRaw(err)
				}
				if se, ok := err.(*xml.SyntaxError); ok {
					return NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: line=%v, error=%v", se.Line, se.Error())).SetRaw(err)
				}
				return NewHTTPError(http.StatusBadRequest, err.Error()).SetRaw(err)
			}
			return err
		},
		MIMEApplicationForm: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
			return NamedStructMap(ctx.Echo(), i, ctx.Request().PostForm().All(), ``, filter...)
		},
		MIMEMultipartForm: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
			return NamedStructMap(ctx.Echo(), i, ctx.Request().Form().All(), ``, filter...)
		},
		`*`: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
			return NamedStructMap(ctx.Echo(), i, ctx.Request().Form().All(), ``, filter...)
		},
	}
	// DefaultHTMLFilter html filter (`form_filter:"html"`)
	DefaultHTMLFilter = func(v string) (r string) {
		return v
	}
)
