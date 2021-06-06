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

//Package widgets This package contains the base logic for the creation and rendering of field widgets. Base widgets are defined for most input fields,
// both in classic and Bootstrap3 style; custom widgets can be defined and associated to a field, provided that they implement the
// WidgetInterface interface.
package widgets

import (
	"bytes"
	"html/template"

	"github.com/coscms/forms/common"
)

func New(t *template.Template) *Widget {
	return &Widget{template: t}
}

// Widget Simple widget object that gets executed at render time.
type Widget struct {
	template *template.Template
}

// WidgetInterface defines the requirements for custom widgets.
type WidgetInterface interface {
	Render(data interface{}) string
}

// Render executes the internal template and returns the result as a template.HTML object.
func (w *Widget) Render(data interface{}) string {
	buf := bytes.NewBuffer(nil)
	err := w.template.ExecuteTemplate(buf, "main", data)
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

// BaseWidget creates a Widget based on style and inpuType parameters, both defined in the common package.
func BaseWidget(style, inputType, tmplName string) *Widget {
	cachedKey := style + ", " + inputType + ", " + tmplName
	tmpl, err := common.GetOrSetCachedTemplate(cachedKey, func() (*template.Template, error) {
		fpath := common.TmplDir(style) + "/" + style + "/"
		urls := []string{common.LookupPath(fpath + "generic.html")}
		tpath := widgetTmpl(inputType, tmplName)
		urls = append(urls, common.LookupPath(fpath+tpath+".html"))
		return common.ParseFiles(urls...)
	})
	if err != nil {
		panic(err)
	}
	tmpl.Funcs(common.TplFuncs())
	return New(tmpl)
}

func widgetTmpl(inputType, tmpl string) (tpath string) {
	switch inputType {
	case common.BUTTON:
		tpath = "button"
		if len(tmpl) > 0 {
			tpath = tmpl
		}
	case common.TEXTAREA:
		tpath = "text/textareainput"
		if len(tmpl) > 0 {
			tpath = "text/" + tmpl
		}
	case common.PASSWORD:
		tpath = "text/passwordinput"
		if len(tmpl) > 0 {
			tpath = "text/" + tmpl
		}
	case common.TEXT:
		tpath = "text/textinput"
		if len(tmpl) > 0 {
			tpath = "text/" + tmpl
		}
	case common.CHECKBOX:
		tpath = "options/checkbox"
		if len(tmpl) > 0 {
			tpath = "options/" + tmpl
		}
	case common.SELECT:
		tpath = "options/select"
		if len(tmpl) > 0 {
			tpath = "options/" + tmpl
		}
	case common.RADIO:
		tpath = "options/radiobutton"
		if len(tmpl) > 0 {
			tpath = "options/" + tmpl
		}
	case common.RANGE:
		tpath = "number/range"
		if len(tmpl) > 0 {
			tpath = "number/" + tmpl
		}
	case common.NUMBER:
		tpath = "number/number"
		if len(tmpl) > 0 {
			tpath = "number/" + tmpl
		}
	case common.RESET, common.SUBMIT:
		tpath = "button"
		if len(tmpl) > 0 {
			tpath = tmpl
		}
	case common.DATE:
		tpath = "datetime/date"
		if len(tmpl) > 0 {
			tpath = "datetime/" + tmpl
		}
	case common.DATETIME:
		tpath = "datetime/datetime"
		if len(tmpl) > 0 {
			tpath = "datetime/" + tmpl
		}
	case common.TIME:
		tpath = "datetime/time"
		if len(tmpl) > 0 {
			tpath = "datetime/" + tmpl
		}
	case common.DATETIME_LOCAL:
		tpath = "datetime/datetime"
		if len(tmpl) > 0 {
			tpath = "datetime/" + tmpl
		}
	case common.STATIC:
		tpath = "static"
		if len(tmpl) > 0 {
			tpath = tmpl
		}
	case common.SEARCH, common.TEL, common.URL, common.WEEK, common.COLOR, common.EMAIL, common.FILE, common.HIDDEN, common.IMAGE, common.MONTH:
		fallthrough
	default:
		tpath = "input"
		if len(tmpl) > 0 {
			tpath = tmpl
		}
	}
	return
}
