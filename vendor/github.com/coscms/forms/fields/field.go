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

//Package fields This package provides all the input fields logic and customization methods.
package fields

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"

	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/widgets"
)

// Field is a generic type containing all data associated to an input field.
type Field struct {
	Type       string
	Tmpl       string
	Widget     widgets.WidgetInterface // Public Widget field for widget customization
	CurrName   string
	OrigName   string
	Class      common.HTMLAttrValues
	ID         string
	Params     map[string]interface{}
	CSS        map[string]string
	Label      string
	LabelCols  int
	FieldCols  int
	LabelClass common.HTMLAttrValues
	Tag        common.HTMLAttrValues
	Value      string
	Helptext   string
	Errors     []string
	Additional map[string]interface{}
	Choices    interface{}
	ChoiceKeys map[string]ChoiceIndex
	AppendData map[string]interface{}
	TmplStyle  string
	Format     string
	Language   string
	data       map[string]interface{}
}

// FieldWithType creates an empty field of the given type and identified by name.
func FieldWithType(name, t string) *Field {
	return &Field{
		Type:       t,
		Widget:     nil,
		CurrName:   name,
		OrigName:   name,
		Class:      common.HTMLAttrValues{},
		ID:         "",
		Params:     map[string]interface{}{},
		CSS:        map[string]string{},
		Label:      "",
		LabelClass: common.HTMLAttrValues{},
		Tag:        common.HTMLAttrValues{},
		Value:      "",
		Helptext:   "",
		Errors:     []string{},
		Additional: map[string]interface{}{},
		Choices:    nil,
		ChoiceKeys: map[string]ChoiceIndex{},
		AppendData: map[string]interface{}{},
		TmplStyle:  "",
	}
}

func (f *Field) SetTmpl(tmpl string, style ...string) FieldInterface {
	f.Tmpl = tmpl
	if len(f.Tmpl) > 0 && f.Widget != nil && f.Tmpl != tmpl {
		var s string
		if len(style) > 0 {
			s = style[0]
		} else {
			s = f.TmplStyle
		}
		f.Widget = widgets.BaseWidget(s, f.Type, f.Tmpl)
	}
	return f
}

func (f *Field) SetName(name string) {
	f.CurrName = name
}

func (f *Field) OriginalName() string {
	return f.OrigName
}

func (f *Field) Clone() config.FormElement {
	fc := *f
	return &fc
}

func (f *Field) SetLang(lang string) {
	f.Language = lang
}

func (f *Field) Lang() string {
	return f.Language
}

// SetStyle sets the style (e.g.: BASE, BOOTSTRAP) of the field, correctly populating the Widget field.
func (f *Field) SetStyle(style string) FieldInterface {
	f.TmplStyle = style
	f.Widget = widgets.BaseWidget(style, f.Type, f.Tmpl)
	return f
}

func (f *Field) SetData(key string, value interface{}) {
	f.AppendData[key] = value
}

func (f *Field) SetLabelCols(cols int) {
	f.LabelCols = cols
}

func (f *Field) SetFieldCols(cols int) {
	f.FieldCols = cols
}

// Name returns the name of the field.
func (f *Field) Name() string {
	return strings.TrimSuffix(f.CurrName, "[]")
}

func (f *Field) Data() map[string]interface{} {
	if len(f.data) > 0 {
		return f.data
	}
	safeParams := make(map[template.HTMLAttr]interface{})
	for k, v := range f.Params {
		safeParams[template.HTMLAttr(k)] = v
	}
	f.data = map[string]interface{}{
		"classes":      f.Class,
		"id":           f.ID,
		"name":         f.CurrName,
		"params":       safeParams,
		"css":          f.CSS,
		"type":         f.Type,
		"label":        f.Label,
		"labelCols":    f.LabelCols,
		"fieldCols":    f.FieldCols,
		"labelClasses": f.LabelClass,
		"tags":         f.Tag,
		"value":        f.Value,
		"helptext":     f.Helptext,
		"errors":       f.Errors,
		"container":    "form",
		"choices":      f.Choices,
	}
	for k, v := range f.Additional {
		f.data[k] = v
	}
	for k, v := range f.AppendData {
		f.data[k] = v
	}
	return f.data
}

// Render packs all data and executes widget render method.
func (f *Field) Render() template.HTML {
	if f.Widget != nil {
		return template.HTML(f.Widget.Render(f.Data()))
	}
	return template.HTML("")
}

func (f *Field) String() string {
	if f.Widget != nil {
		return f.Widget.Render(f.Data())
	}
	return ""
}

// AddClass adds a class to the field.
func (f *Field) AddClass(class string) FieldInterface {
	f.Class.Add(class)
	return f
}

// RemoveClass removes a class from the field, if it was present.
func (f *Field) RemoveClass(class string) FieldInterface {
	f.Class.Remove(class)
	return f
}

// SetID associates the given id to the field, overwriting any previous id.
func (f *Field) SetID(id string) FieldInterface {
	f.ID = id
	return f
}

// SetLabel saves the label to be rendered along with the field.
func (f *Field) SetLabel(label string) FieldInterface {
	f.Label = label
	return f
}

// AddLabelClass allows to define custom classes for the label.
func (f *Field) AddLabelClass(class string) FieldInterface {
	f.LabelClass.Add(class)
	return f
}

// RemoveLabelClass removes the given class from the field label.
func (f *Field) RemoveLabelClass(class string) FieldInterface {
	f.LabelClass.Remove(class)
	return f
}

// SetParam adds a parameter (defined as key-value pair) in the field.
func (f *Field) SetParam(key string, value interface{}) FieldInterface {
	switch key {
	case `class`:
		f.AddClass(value.(string))
	default:
		f.Params[key] = value
	}
	return f
}

// DeleteParam removes a parameter identified by key from the field.
func (f *Field) DeleteParam(key string) FieldInterface {
	delete(f.Params, key)
	return f
}

// AddCSS adds a custom CSS style the field.
func (f *Field) AddCSS(key, value string) FieldInterface {
	f.CSS[key] = value
	return f
}

// RemoveCSS removes CSS options identified by key from the field.
func (f *Field) RemoveCSS(key string) FieldInterface {
	delete(f.CSS, key)
	return f
}

// Disabled add the "disabled" tag to the field, making it unresponsive in some environments (e.g. Bootstrap).
func (f *Field) Disabled() FieldInterface {
	f.AddTag("disabled")
	return f
}

// Enabled removes the "disabled" tag from the field, making it responsive.
func (f *Field) Enabled() FieldInterface {
	f.RemoveTag("disabled")
	return f
}

// AddTag adds a no-value parameter (e.g.: checked, disabled) to the field.
func (f *Field) AddTag(tag string) FieldInterface {
	f.Tag.Add(tag)
	return f
}

// RemoveTag removes a no-value parameter from the field.
func (f *Field) RemoveTag(tag string) FieldInterface {
	f.Tag.Remove(tag)
	return f
}

// SetValue sets the value parameter for the field.
func (f *Field) SetValue(value string) FieldInterface {
	f.Value = value
	f.SetSelected(f.Value)
	return f
}

// SetHelptext saves the field helptext.
func (f *Field) SetHelptext(text string) FieldInterface {
	f.Helptext = text
	return f
}

// AddError adds an error string to the field. It's valid only for Bootstrap forms.
func (f *Field) AddError(err string) FieldInterface {
	f.Errors = append(f.Errors, err)
	return f
}

// MultipleChoice configures the SelectField to accept and display multiple choices.
// It has no effect if type is not SELECT.
func (f *Field) MultipleChoice() FieldInterface {
	switch f.Type {
	case common.SELECT:
		f.AddTag("multiple")
		fallthrough
	case common.CHECKBOX:
		// fix name if necessary
		if !strings.HasSuffix(f.CurrName, "[]") {
			f.CurrName = f.CurrName + "[]"
		}
	}
	return f
}

// SingleChoice configures the Field to accept and display only one choice (valid for SelectFields only).
// It has no effect if type is not SELECT.
func (f *Field) SingleChoice() FieldInterface {
	switch f.Type {
	case common.SELECT:
		f.RemoveTag("multiple")
		fallthrough
	case common.CHECKBOX:
		if strings.HasSuffix(f.CurrName, "[]") {
			f.CurrName = strings.TrimSuffix(f.CurrName, "[]")
		}
	}
	return f
}

// AddSelected If the field is configured as "multiple", AddSelected adds a selected value to the field (valid for SelectFields only).
// It has no effect if type is not SELECT.
func (f *Field) AddSelected(opt ...string) FieldInterface {
	switch f.Type {
	case common.SELECT:
		for _, v := range opt {
			i, ok := f.ChoiceKeys[v]
			if !ok {
				continue
			}
			choice := f.Choices.(map[string][]InputChoice)
			if vc, ok := choice[i.Group]; ok {
				if len(vc) > i.Index {
					choice[i.Group][i.Index].Checked = true
				}
			}
		}
	case common.RADIO, common.CHECKBOX:
		choice := f.Choices.([]InputChoice)
		size := len(choice)
		for _, v := range opt {
			i, ok := f.ChoiceKeys[v]
			if !ok {
				continue
			}
			if size > i.Index {
				choice[i.Index].Checked = true
			}
		}
	}
	return f
}

func (f *Field) SetSelected(opt ...string) FieldInterface {
	switch f.Type {
	case common.SELECT:
		choice := f.Choices.(map[string][]InputChoice)
		for key, i := range f.ChoiceKeys {
			vc, ok := choice[i.Group]
			if !ok || len(vc) <= i.Index {
				continue
			}
			checked := false
			for _, v := range opt {
				if key == v {
					checked = true
					break
				}
			}
			choice[i.Group][i.Index].Checked = checked
		}
	case common.RADIO, common.CHECKBOX:
		choice := f.Choices.([]InputChoice)
		size := len(choice)
		for key, i := range f.ChoiceKeys {
			if size <= i.Index {
				continue
			}
			checked := false
			for _, v := range opt {
				if key == v {
					checked = true
					break
				}
			}
			choice[i.Index].Checked = checked
		}
	}
	return f
}

//RemoveSelected If the field is configured as "multiple", AddSelected removes the selected value from the field (valid for SelectFields only).
// It has no effect if type is not SELECT.
func (f *Field) RemoveSelected(opt string) FieldInterface {
	switch f.Type {
	case common.SELECT:
		i := f.ChoiceKeys[opt]
		if vc, ok := f.Choices.(map[string][]InputChoice)[i.Group]; ok {
			if len(vc) > i.Index {
				f.Choices.(map[string][]InputChoice)[i.Group][i.Index].Checked = false
			}
		}

	case common.RADIO, common.CHECKBOX:
		size := len(f.Choices.([]InputChoice))
		i := f.ChoiceKeys[opt]
		if size > i.Index {
			f.Choices.([]InputChoice)[i.Index].Checked = false
		}
	}
	return f
}

func (f *Field) AddChoice(key, value interface{}, checked ...bool) FieldInterface {
	var _checked bool
	if len(checked) > 0 && checked[0] {
		_checked = true
	}
	switch f.Type {
	case common.SELECT:
		if f.Choices == nil {
			f.Choices = map[string][]InputChoice{
				"": []InputChoice{
					{
						ID:      fmt.Sprint(key),
						Val:     fmt.Sprint(value),
						Checked: _checked,
					},
				},
			}
		} else {
			v, _ := f.Choices.(map[string][]InputChoice)
			v[""] = append(v[""], InputChoice{
				ID:      fmt.Sprint(key),
				Val:     fmt.Sprint(value),
				Checked: _checked,
			})
			f.Choices = v
		}

	case common.RADIO, common.CHECKBOX:
		if f.Choices == nil {
			f.Choices = []InputChoice{
				{
					ID:      fmt.Sprint(key),
					Val:     fmt.Sprint(value),
					Checked: _checked,
				},
			}
		} else {
			v, _ := f.Choices.([]InputChoice)
			v = append(v, InputChoice{
				ID:      fmt.Sprint(key),
				Val:     fmt.Sprint(value),
				Checked: _checked,
			})
			f.Choices = v
		}
	}
	return f
}

// SetChoices takes as input a dictionary whose key-value entries are defined as follows: key is the group name (the empty string
// is the default group that is not explicitly rendered) and value is the list of choices belonging to that group.
// Grouping is only useful for Select fields, while groups are ignored in Radio fields.
// It has no effect if type is not SELECT.
func (f *Field) SetChoices(choices interface{}, saveIndex ...bool) FieldInterface {
	if choices == nil {
		return f
	}
	switch f.Type {
	case common.SELECT:
		var ch map[string][]InputChoice
		if c, ok := choices.(map[string][]InputChoice); ok {
			ch = c
		} else {
			c, y := choices.([]InputChoice)
			if !y {
				if v, y := choices.([]string); y {
					c = []InputChoice{
						InputChoice{},
					}
					switch len(v) {
					case 3:
						c[0].Checked, _ = strconv.ParseBool(v[2])
						fallthrough
					case 2:
						c[0].Val = v[1]
						fallthrough
					case 1:
						c[0].ID = v[0]
					}
				}
			}
			ch = map[string][]InputChoice{"": c}
		}
		f.Choices = ch
		if len(saveIndex) < 1 || saveIndex[0] {
			for k, v := range ch {
				for idx, ipt := range v {
					f.ChoiceKeys[ipt.ID] = ChoiceIndex{Group: k, Index: idx}
				}
			}
		}

	case common.RADIO, common.CHECKBOX:
		c, y := choices.([]InputChoice)
		if !y {
			if v, y := choices.([]string); y {
				c = []InputChoice{
					InputChoice{},
				}
				switch len(v) {
				case 3:
					c[0].Checked, _ = strconv.ParseBool(v[2])
					fallthrough
				case 2:
					c[0].Val = v[1]
					fallthrough
				case 1:
					c[0].ID = v[0]
				}
			}
		}
		f.Choices = c
		if len(saveIndex) < 1 || saveIndex[0] {
			for idx, ipt := range c {
				f.ChoiceKeys[ipt.ID] = ChoiceIndex{Group: "", Index: idx}
			}
		}
	}
	return f
}

// SetText saves the provided text as content of the field, usually a TextAreaField.
func (f *Field) SetText(text string) FieldInterface {
	if f.Type == common.BUTTON ||
		f.Type == common.SUBMIT ||
		f.Type == common.RESET ||
		f.Type == common.STATIC ||
		f.Type == common.TEXTAREA {
		f.Additional["text"] = text
	}
	return f
}

func (f *Field) Element() *config.Element {
	elem := &config.Element{
		ID:         f.ID,
		Type:       f.Type,
		Name:       f.CurrName,
		Label:      f.Label,
		LabelCols:  f.LabelCols,
		FieldCols:  f.FieldCols,
		Value:      f.Value,
		HelpText:   f.Helptext,
		Template:   f.Tmpl,
		Valid:      ``,
		Attributes: make([][]string, 0),
		Choices:    make([]*config.Choice, 0),
		Elements:   make([]*config.Element, 0),
		Format:     f.Format,
	}
	if f.AppendData != nil && len(f.AppendData) > 0 {
		elem.Data = f.AppendData
	}
	var (
		temp string
		join string
	)
	for _, c := range f.Class {
		temp += join + c
		join = ` `
	}
	if len(temp) > 0 {
		elem.Attributes = append(elem.Attributes, []string{`class`, temp})
		temp = ``
		join = ``
	}
	for _, c := range f.Tag {
		elem.Attributes = append(elem.Attributes, []string{c})
	}
	for c, v := range f.Params {
		elem.Attributes = append(elem.Attributes, []string{c, fmt.Sprintf(`%v`, v)})
	}
	for _, c := range f.CSS {
		temp += join + c
		join = `;`
	}
	if len(temp) > 0 {
		elem.Attributes = append(elem.Attributes, []string{`style`, temp})
		temp = ``
		join = ``
	}
	switch choices := f.Choices.(type) {
	case map[string][]InputChoice:
		for k, items := range choices {
			for _, v := range items {
				elem.Choices = append(elem.Choices, &config.Choice{
					Group:   k,
					Option:  []string{v.ID, v.Val},
					Checked: v.Checked,
				})
			}
		}
	case []InputChoice:
		for _, v := range choices {
			elem.Choices = append(elem.Choices, &config.Choice{
				Group:   ``,
				Option:  []string{v.ID, v.Val},
				Checked: v.Checked,
			})
		}
	}
	return elem
}
