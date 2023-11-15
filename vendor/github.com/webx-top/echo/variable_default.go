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
	"strings"
	"time"

	"github.com/webx-top/com"
	"github.com/webx-top/echo/encoding/json"
	"github.com/webx-top/echo/param"
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
	DefaultBinderDecoders = map[string]func(interface{}, Context, BinderValueCustomDecoders, ...FormDataFilter) error{
		MIMEApplicationJSON: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
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
		MIMEApplicationXML: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
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
		MIMEApplicationForm: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
			return FormToStructWithDecoder(ctx.Echo(), i, ctx.Request().PostForm().All(), ``, valueDecoders, filter...)
		},
		MIMEMultipartForm: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
			_, err := ctx.Request().MultipartForm()
			if err != nil {
				return err
			}
			return FormToStructWithDecoder(ctx.Echo(), i, ctx.Request().Form().All(), ``, valueDecoders, filter...)
		},
		`*`: func(i interface{}, ctx Context, valueDecoders BinderValueCustomDecoders, filter ...FormDataFilter) error {
			return FormToStructWithDecoder(ctx.Echo(), i, ctx.Request().Form().All(), ``, valueDecoders, filter...)
		},
	}
	// DefaultHTMLFilter html filter (`form_filter:"html"`)
	DefaultHTMLFilter = func(v string) (r string) {
		return v
	}
	DefaultBinderValueEncoders = map[string]BinderValueEncoder{
		`joinKVRows`: binderValueEncoderJoinKVRows,
		`join`:       binderValueEncoderJoin,
		`unix2time`:  binderValueEncoderUnix2time,
	}
	DefaultBinderValueDecoders = map[string]BinderValueDecoder{
		`splitKVRows`: binderValueDecoderSplitKVRows,
		`split`:       binderValueDecoderSplit,
		`time2unix`:   binderValueDecoderTime2unix,
	}
)

func binderValueDecoderSplitKVRows(field string, values []string, seperator string) (interface{}, error) {
	return com.SplitKVRows(values[0], seperator), nil
}

func binderValueDecoderSplit(field string, values []string, seperator string) (interface{}, error) {
	return strings.Split(values[0], seperator), nil
}

func binderValueEncoderJoin(field string, value interface{}, seperator string) []string {
	return []string{strings.Join(param.AsStdStringSlice(value), seperator)}
}

func binderValueEncoderJoinKVRows(field string, value interface{}, seperator string) []string {
	result := com.JoinKVRows(value, seperator)
	if len(result) == 0 {
		return nil
	}
	return []string{result}
}

func binderValueEncoderUnix2time(field string, value interface{}, seperator string) []string {
	ts := param.AsInt64(value)
	if ts <= 0 {
		return []string{}
	}
	return []string{param.AsString(time.Unix(ts, 0))}
}

func binderValueDecoderTime2unix(field string, values []string, layout string) (interface{}, error) {
	return param.AsDateTime(values[0], layout).Unix(), nil
}
