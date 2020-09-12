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
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/webx-top/echo/encoding/json"
)

var (
	DefaultAcceptFormats = map[string]string{
		//json
		`application/json`:       `json`,
		`text/javascript`:        `json`,
		`application/javascript`: `json`,

		//xml
		`application/xml`: `xml`,
		`text/xml`:        `xml`,

		//text
		`text/plain`: `text`,

		//html
		`*/*`:               `html`,
		`application/xhtml`: `html`,
		`text/html`:         `html`,

		//default
		`*`: `html`,
	}
	DefaultFormatRenderers = map[string]func(c Context, data interface{}) error{
		`json`: func(c Context, data interface{}) error {
			return c.JSON(c.Data())
		},
		`jsonp`: func(c Context, data interface{}) error {
			return c.JSONP(c.Query(c.Echo().JSONPVarName), c.Data())
		},
		`xml`: func(c Context, data interface{}) error {
			return c.XML(c.Data())
		},
		`text`: func(c Context, data interface{}) error {
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
			return json.NewDecoder(body).Decode(i)
		},
		MIMEApplicationXML: func(i interface{}, ctx Context, filter ...FormDataFilter) error {
			body := ctx.Request().Body()
			if body == nil {
				return NewHTTPError(http.StatusBadRequest, "Request body can't be nil")
			}
			defer body.Close()
			return xml.NewDecoder(body).Decode(i)
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
