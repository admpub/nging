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
package forms

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/admpub/json5"

	"github.com/webx-top/tagfast"

	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"
	"github.com/webx-top/validation"
)

func UnmarshalFile(filename string) (r *config.Config, err error) {
	filename, err = filepath.Abs(filename)
	if err != nil {
		return
	}
	return common.GetOrSetCachedConfig(filename, func() (*config.Config, error) {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		r = &config.Config{}
		err = json5.Unmarshal(b, r)
		if err != nil {
			return nil, err
		}
		fmt.Println(`cache form config:`, filename)
		return r, nil
	})
}

func Unmarshal(b []byte, key string) (r *config.Config, err error) {
	return common.GetOrSetCachedConfig(key, func() (*config.Config, error) {
		r := &config.Config{}
		err = json5.Unmarshal(b, r)
		if err != nil {
			return nil, err
		}
		fmt.Println(`cache form config:`, key)
		return r, nil
	})
}

func NewWithModelConfig(m interface{}, r *config.Config) *Form {
	form := NewWithConfig(r)
	form.SetModel(m).ParseFromConfig()
	return form
}

func (form *Form) Generate(m interface{}, jsonFile string) error {
	r, err := UnmarshalFile(jsonFile)
	if err != nil {
		return err
	}
	form.Init(r).SetModel(m)
	form.ParseFromConfig()
	return nil
}

func (form *Form) ParseFromJSONFile(jsonFile string) error {
	r, err := UnmarshalFile(jsonFile)
	if err != nil {
		return err
	}
	form.Init(r)
	form.ParseFromConfig()
	return nil
}

func (form *Form) ParseFromJSON(b []byte, key string) error {
	r, err := Unmarshal(b, key)
	if err != nil {
		return err
	}
	form.Init(r)
	form.ParseFromConfig()
	return nil
}

func (form *Form) ValidFromJSONFile(jsonFile string) error {
	r, err := UnmarshalFile(jsonFile)
	if err != nil {
		return err
	}
	form.Init(r)
	form.ValidFromConfig()
	return nil
}

func (form *Form) ValidFromJSON(b []byte, key string) error {
	r, err := Unmarshal(b, key)
	if err != nil {
		return err
	}
	form.Init(r)
	form.ValidFromConfig()
	return nil
}

func (form *Form) ValidFromConfig() *Form {
	form.Validate()
	if form.Model == nil {
		return form
	}
	t := reflect.TypeOf(form.Model)
	v := reflect.ValueOf(form.Model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	r := form.config
	form.ValidElements(r.Elements, t, v)
	return form
}

// Filter 过滤客户端提交的数据
func (form *Form) Filter(values url.Values) (url.Values, *validation.ValidationError) {
	form.Validate()
	r := url.Values{}
	var err *validation.ValidationError
	for _, ele := range form.config.Elements {
		switch ele.Type {
		case `langset`, `fieldset`:
			for _, e := range ele.Elements {
				r, err = form.FilterByElement(values, r, e)
				if err != nil {
					return r, err
				}
			}
		default:
			r, err = form.FilterByElement(values, r, ele)
			if err != nil {
				return r, err
			}
		}
	}
	return r, err
}

// FilterByElement 过滤单个元素
func (form *Form) FilterByElement(input url.Values, output url.Values, ele *config.Element) (url.Values, *validation.ValidationError) {
	if len(ele.Valid) == 0 {
		if vals, ok := input[ele.Name]; ok {
			output[ele.Name] = vals
		}
	} else {
		if vals, ok := input[ele.Name]; ok {
			for _, val := range vals {
				if !form.valid.ValidField(ele.Name, val, ele.Valid) {
					return output, form.Error()
				}
			}
			output[ele.Name] = vals
		}
	}
	return output, form.Error()
}

func (form *Form) ValidElements(elements []*config.Element, t reflect.Type, v reflect.Value) {
	for _, ele := range elements {
		switch ele.Type {
		case `langset`:
			//form.ValidElements(ele.Elements, t, v)
		case `fieldset`:
			for _, e := range ele.Elements {
				if !form.IsIgnored(e.Name) {
					form.validElement(e, t, v)
				}
			}
		default:
			if !form.IsIgnored(ele.Name) {
				form.validElement(ele, t, v)
			}
		}
	}
}

func (form *Form) IsIgnored(fieldName string) bool {
	for _, name := range form.ignoreValid {
		if fieldName == name {
			return true
		}
	}
	return false
}

func (form *Form) CloseValid(fieldName ...string) *Form {
	if form.ignoreValid == nil {
		form.ignoreValid = []string{}
	}
	form.ignoreValid = append(form.ignoreValid, fieldName...)
	return form
}

func (form *Form) ParseFromConfig(insertErrors ...bool) *Form {
	t := reflect.TypeOf(form.Model)
	v := reflect.ValueOf(form.Model)
	if t != nil && t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	r := form.config
	form.ParseElements(form, r.Elements, r.Languages, t, v, ``)
	if len(insertErrors) < 1 || insertErrors[0] {
		form.InsertErrors()
	}
	for _, attr := range r.Attributes {
		var k, v string
		switch len(attr) {
		case 2:
			v = attr[1]
			fallthrough
		case 1:
			k = attr[0]
			form.SetParam(k, v)
		}
	}
	if len(r.ID) > 0 {
		form.SetID(r.ID)
	}
	if r.WithButtons {
		if r.Buttons == nil {
			r.Buttons = []string{}
		}
		form.AddButton(r.BtnsTemplate, r.Buttons...)
	}
	for key, val := range r.Data {
		form.SetData(key, val)
	}
	return form
}

func (form *Form) ParseElements(es ElementSetter, elements []*config.Element, langs []*config.Language, t reflect.Type, v reflect.Value, lang string) {
	for _, ele := range elements {
		switch ele.Type {
		case `langset`:
			if ele.Languages == nil {
				ele.Languages = langs
			}
			f := form.NewLangSet(ele.Name, ele.Languages)
			if len(ele.Template) > 0 {
				f.SetTemplate(ele.Template)
			}
			f.SetData("container", "langset")
			for key, val := range ele.Data {
				f.SetData(key, val)
			}
			form.ParseElements(f, ele.Elements, ele.Languages, t, v, ``)
			for _, v := range ele.Attributes {
				switch len(v) {
				case 2:
					f.SetParam(v[0], v[1])
				case 1:
					f.AddTag(v[0])
				}
			}
			es.Elements(f)
		case `fieldset`:
			elems := []fields.FieldInterface{}
			for _, e := range ele.Elements {
				elem := form.parseElement(e, t, v)
				if elem != nil {
					elems = append(elems, elem)
				}
			}
			f := form.NewFieldSet(ele.Name, form.labelFn(ele.Label), elems...)
			if len(ele.Template) > 0 {
				f.SetTemplate(ele.Template)
			}
			f.SetData("container", "fieldset")
			for key, val := range ele.Data {
				f.SetData(key, val)
			}
			f.SetLabelCols(ele.LabelCols)
			f.SetFieldCols(ele.FieldCols)
			f.SetLang(lang)
			es.Elements(f)
		default:
			f := form.parseElement(ele, t, v)
			if f != nil {
				f.SetLang(lang)
				es.Elements(f)
			}
		}
	}
}

func (form *Form) parseElement(ele *config.Element, typ reflect.Type, val reflect.Value) (f *fields.Field) {
	var sv string
	value := val
	if form.Model != nil && !form.IsOmit(ele.Name) {
		parts := strings.Split(ele.Name, `.`)
		isValid := true
		for _, field := range parts {
			if value.Kind() == reflect.Ptr {
				if value.IsNil() {
					isValid = false
					break
				}
				value = value.Elem()
			}
			switch typ.Kind() {
			case reflect.Map:
				index := reflect.ValueOf(field)
				value = value.MapIndex(index)
			case reflect.Slice:
				index, _ := strconv.Atoi(field)
				if index >= value.Len() {
					isValid = false
					goto OUTLOOP
				}
				value = value.Index(index)
			case reflect.Struct:
				field = strings.Title(field)
				value = value.FieldByName(field)
			default:
				isValid = false
				goto OUTLOOP
			}
			if !value.IsValid() {
				isValid = false
				break
			}
			if value.Kind() == reflect.Interface {
				value = reflect.ValueOf(value.Interface())
			}
			value = reflect.Indirect(value)
			kind := value.Kind()
			if kind != reflect.Struct && kind != reflect.Map && kind != reflect.Slice {
				break
			}
			typ = value.Type()
		}

	OUTLOOP:
		if isValid {
			sv = fmt.Sprintf("%v", value.Interface())
		}
	}
	isStruct := typ != nil && typ.Kind() == reflect.Struct
	switch ele.Type {
	case common.DATE:
		dateFormat := fields.DATE_FORMAT
		if len(ele.Format) > 0 {
			dateFormat = ele.Format
		} else if isStruct {
			if structField, ok := typ.FieldByName(strings.Title(ele.Name)); ok {
				if format := tagfast.Value(typ, structField, `form_format`); len(format) > 0 {
					dateFormat = format
				}
			}
		}
		f = fields.TextField(ele.Name, ele.Type)
		if v, isEmpty := fields.ConvertTime(value.Interface()); !v.IsZero() {
			f.SetValue(v.Format(dateFormat))
		} else if isEmpty {
			f.SetValue(``)
		} else {
			f.SetValue(ele.Value)
		}

	case common.DATETIME:
		dateFormat := fields.DATETIME_FORMAT
		if len(ele.Format) > 0 {
			dateFormat = ele.Format
		} else if isStruct {
			if structField, ok := typ.FieldByName(strings.Title(ele.Name)); ok {
				if format := tagfast.Value(typ, structField, `form_format`); len(format) > 0 {
					dateFormat = format
				}
			}
		}
		f = fields.TextField(ele.Name, ele.Type)
		if v, isEmpty := fields.ConvertTime(value.Interface()); !v.IsZero() {
			f.SetValue(v.Format(dateFormat))
		} else if isEmpty {
			f.SetValue(``)
		} else {
			f.SetValue(ele.Value)
		}

	case common.DATETIME_LOCAL:
		dateFormat := fields.DATETIME_FORMAT
		if len(ele.Format) > 0 {
			dateFormat = ele.Format
		} else if isStruct {
			if structField, ok := typ.FieldByName(strings.Title(ele.Name)); ok {
				if format := tagfast.Value(typ, structField, `form_format`); len(format) > 0 {
					dateFormat = format
				}
			}
		}
		f = fields.TextField(ele.Name, ele.Type)
		if v, isEmpty := fields.ConvertTime(value.Interface()); !v.IsZero() {
			f.SetValue(v.Local().Format(dateFormat))
		} else if isEmpty {
			f.SetValue(``)
		} else {
			f.SetValue(ele.Value)
		}

	case common.TIME:
		dateFormat := fields.TIME_FORMAT
		if len(ele.Format) > 0 {
			dateFormat = ele.Format
		} else if isStruct {
			if structField, ok := typ.FieldByName(strings.Title(ele.Name)); ok {
				if format := tagfast.Value(typ, structField, `form_format`); len(format) > 0 {
					dateFormat = format
				}
			}
		}
		f = fields.TextField(ele.Name, ele.Type)
		if v, isEmpty := fields.ConvertTime(value.Interface()); !v.IsZero() {
			f.SetValue(v.Format(dateFormat))
		} else if isEmpty {
			f.SetValue(``)
		} else {
			f.SetValue(ele.Value)
		}

	case common.TEXT:
		f = fields.TextField(ele.Name, ele.Type)
		format := ele.Format
		if len(format) == 0 && isStruct {
			if structField, ok := typ.FieldByName(strings.Title(ele.Name)); ok {
				format = tagfast.Value(typ, structField, `form_format`)
			}
		}
		if len(format) > 0 { //时间格式
			if vt, isEmpty := fields.ConvertTime(sv); !vt.IsZero() {
				f.SetValue(vt.Format(format))
			} else if isEmpty {
				f.SetValue(``)
			}
		} else {
			if len(sv) == 0 {
				f.SetValue(ele.Value)
			} else {
				f.SetValue(sv)
			}
		}

	case common.COLOR, common.EMAIL, common.FILE, common.HIDDEN, common.IMAGE, common.MONTH, common.SEARCH, common.URL, common.TEL, common.WEEK, common.NUMBER, common.PASSWORD:
		f = fields.TextField(ele.Name, ele.Type)
		if len(sv) == 0 {
			f.SetValue(ele.Value)
		} else {
			f.SetValue(sv)
		}

	case common.CHECKBOX, common.RADIO:
		choices := []fields.InputChoice{}
		hasSet := len(sv) > 0
		for _, v := range ele.Choices {
			if v.Checked {
				if hasSet && sv != v.Option[0] {
					v.Checked = false
				}
			} else {
				if hasSet {
					v.Checked = sv == v.Option[0]
				}
			}
			ic := fields.InputChoice{
				ID:      v.Option[0],
				Val:     form.labelFn(v.Option[1]),
				Checked: v.Checked,
			}
			choices = append(choices, ic)
		}
		if ele.Type == common.CHECKBOX {
			f = fields.CheckboxField(ele.Name, choices)
		} else {
			f = fields.RadioField(ele.Name, choices)
		}
		if !hasSet {
			f.SetValue(ele.Value)
		} else {
			f.SetValue(sv)
		}

	case common.RANGE:
		f = fields.FieldWithType(ele.Name, ele.Type)
		if len(sv) == 0 {
			f.SetValue(ele.Value)
		} else {
			f.SetValue(sv)
		}

	case common.BUTTON, common.RESET, common.SUBMIT, common.STATIC, common.TEXTAREA:
		f = fields.FieldWithType(ele.Name, ele.Type)
		if len(sv) == 0 {
			f.SetText(ele.Value)
		} else {
			f.SetText(sv)
		}

	case common.SELECT:
		choices := map[string][]fields.InputChoice{}
		hasSet := len(sv) > 0
		for _, v := range ele.Choices {
			if _, ok := choices[v.Group]; !ok {
				choices[v.Group] = []fields.InputChoice{}
			}
			if v.Checked {
				if hasSet && sv != v.Option[0] {
					v.Checked = false
				}
			} else {
				if hasSet {
					v.Checked = sv == v.Option[0]
				}
			}
			ic := fields.InputChoice{
				ID:      v.Option[0],
				Val:     form.labelFn(v.Option[1]),
				Checked: v.Checked,
			}
			choices[v.Group] = append(choices[v.Group], ic)
		}
		f = fields.SelectField(ele.Name, choices)
		if !hasSet {
			f.SetValue(ele.Value)
		} else {
			f.SetValue(sv)
		}

	default:
		return nil
	}
	for _, v := range ele.Attributes {
		switch len(v) {
		case 2:
			f.SetParam(v[0], v[1])
		case 1:
			f.AddTag(v[0])
		}
	}
	f.SetHelptext(form.labelFn(ele.HelpText))
	f.SetLabel(form.labelFn(ele.Label))
	f.SetTemplate(ele.Template)
	f.SetID(ele.ID)
	if len(ele.Valid) > 0 {
		form.validTagFn(ele.Valid, f)
	}
	for key, val := range ele.Data {
		f.SetData(key, val)
	}
	f.SetLabelCols(ele.LabelCols)
	f.SetFieldCols(ele.FieldCols)
	return f
}

func (form *Form) validElement(ele *config.Element, typ reflect.Type, val reflect.Value) bool {
	if len(ele.Valid) == 0 {
		return true
	}
	parts := strings.Split(ele.Name, `.`)
	value := val
	isValid := true
	for _, field := range parts {
		field = strings.Title(field)
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				value.Set(reflect.New(value.Type().Elem()))
			}
			value = value.Elem()
		}
		value = value.FieldByName(field)
		if !value.IsValid() {
			isValid = false
			break
		}
	}
	if isValid {
		sv := fmt.Sprintf("%v", value.Interface())
		isValid = form.valid.ValidField(ele.Name, sv, ele.Valid)
	}
	return isValid
}

func (form *Form) ToJSONBlob(args ...*config.Config) (r []byte, err error) {
	var config *config.Config
	if len(args) > 0 {
		config = args[0]
	}
	if config == nil {
		config = form.ToConfig()
	}
	r, err = json.MarshalIndent(config, ``, `  `)
	return
}

func (form *Form) NewConfig() *config.Config {
	return NewConfig()
}

func (form *Form) ToConfig() *config.Config {
	conf := form.NewConfig()
	form.ParseModel()
	for _, v := range form.FieldList {
		var element *config.Element
		switch f := v.(type) {
		case *FieldSetType:
			element = &config.Element{
				ID:         ``,
				Type:       `fieldset`,
				Name:       ``,
				Label:      f.Name(),
				Value:      ``,
				HelpText:   ``,
				Template:   ``,
				Valid:      ``,
				Attributes: make([][]string, 0),
				Choices:    make([]*config.Choice, 0),
				Elements:   make([]*config.Element, 0),
			}
			var temp string
			var join string
			for _, c := range f.Classes {
				temp += join + c
				join = ` `
			}
			if len(temp) > 0 {
				element.Attributes = append(element.Attributes, []string{`class`, temp})
				temp = ``
				join = ``
			}
			for _, c := range f.Tags {
				element.Attributes = append(element.Attributes, []string{c})
			}
			for _, ff := range f.FieldList {
				element.Elements = append(element.Elements, ff.Element())
			}
		case fields.FieldInterface:
			element = f.Element()
		}
		if element != nil {
			conf.Elements = append(conf.Elements, element)
		}
	}
	return conf
}
