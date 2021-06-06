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

// Package forms This package provides form creation and rendering functionalities, as well as FieldSet definition.
// Two kind of forms can be created: base forms and Bootstrap3 compatible forms; even though the latters are automatically provided
// the required classes to make them render correctly in a Bootstrap environment, every form can be given custom parameters such as
// classes, id, generic parameters (in key-value form) and stylesheet options.
package forms

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/coscms/forms/common"
	"github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"
	"github.com/webx-top/validation"
)

// Form methods: POST or GET.
const (
	POST = "POST"
	GET  = "GET"
)

func NewWithConfig(c *config.Config, args ...interface{}) *Form {
	form := New()
	form.Init(c, args...)
	return form
}

func NewWithConfigFile(m interface{}, configJSONFile string) *Form {
	config, err := UnmarshalFile(configJSONFile)
	if err != nil {
		panic(err)
	}
	return NewWithModelConfig(m, config)
}

func New() *Form {
	return &Form{
		FieldList:             []config.FormElement{},
		fieldMap:              make(map[string]int),
		containerMap:          make(map[string]string),
		Style:                 common.BASE,
		Class:                 common.HTMLAttrValues{},
		ID:                    "",
		Params:                map[string]string{},
		CSS:                   map[string]string{},
		Method:                POST,
		AppendData:            map[string]interface{}{},
		labelFn:               func(s string) string { return s },
		validTagFn:            Html5Validate,
		beforeRender:          []func(){},
		omitOrMustFieldsValue: map[string]bool{},
	}
}

func NewFromModel(m interface{}, c *config.Config) *Form {
	form := NewWithConfig(c, m)
	form.SetModel(m)
	form.ParseModel(m)
	return form
}

// Form structure.
type Form struct {
	AppendData map[string]interface{} `json:"appendData,omitempty" xml:"appendData,omitempty"`
	Model      interface{}            `json:"model" xml:"model"`
	FieldList  []config.FormElement   `json:"fieldList" xml:"fieldList"`
	Style      string                 `json:"style" xml:"style"`
	Class      common.HTMLAttrValues  `json:"class" xml:"class"`
	ID         string                 `json:"id" xml:"id"`
	Params     map[string]string      `json:"params" xml:"params"`
	CSS        map[string]string      `json:"css" xml:"css"`
	Method     string                 `json:"method" xml:"method"`
	Action     template.HTML          `json:"action" xml:"action"`

	omitOrMustFieldsValue map[string]bool //true:omit; false:must
	omitAllFieldsValue    bool
	fieldMap              map[string]int
	containerMap          map[string]string
	ignoreValid           []string
	template              *template.Template
	valid                 *validation.Validation
	labelFn               func(string) string
	validTagFn            func(string, fields.FieldInterface)
	config                *config.Config
	beforeRender          []func()
	debug                 bool
	data                  map[string]interface{}
}

func (f *Form) Debug(args ...bool) *Form {
	if len(args) > 0 {
		f.debug = args[0]
	} else {
		f.debug = true
	}
	return f
}

func (f *Form) Config() *config.Config {
	return f.config
}

func (f *Form) IsDebug() bool {
	return f.debug
}

// Must 使用结构体实例中某些字段的值（和OmitAll配合起来使用）
func (f *Form) Must(fields ...string) *Form {
	for _, field := range fields {
		f.omitOrMustFieldsValue[field] = true
	}
	return f
}

func (f *Form) ResetOmitOrMust() *Form {
	f.omitOrMustFieldsValue = map[string]bool{}
	return f
}

// Omit 忽略结构体实例中某些字段的值
func (f *Form) Omit(fields ...string) *Form {
	for _, field := range fields {
		f.omitOrMustFieldsValue[field] = false
	}
	return f
}

// OmitAll 忽略结构体实例中所有字段的值
func (f *Form) OmitAll(on ...bool) *Form {
	if len(on) > 0 {
		f.omitAllFieldsValue = on[0]
	} else {
		f.omitAllFieldsValue = true
	}
	return f
}

// IsOmit 是否忽略结构体实例中指定字段的值
func (f *Form) IsOmit(fieldName string) (omitFieldValue bool) {
	if f.omitAllFieldsValue {
		//默认为忽略，设置为不忽略时才不忽略
		omit, ok := f.omitOrMustFieldsValue[fieldName]
		if !ok || omit {
			omitFieldValue = true
		}
	} else {
		//设置为忽略时，才忽略
		omit, ok := f.omitOrMustFieldsValue[fieldName]
		if ok && omit {
			omitFieldValue = true
		}
	}
	return
}

func (f *Form) SetLabelFunc(fn func(string) string) *Form {
	f.labelFn = fn
	return f
}

func (f *Form) SetValidTagFunc(fn func(string, fields.FieldInterface)) *Form {
	f.validTagFn = fn
	return f
}

func (f *Form) LabelFunc() func(string) string {
	return f.labelFn
}

func (f *Form) ValidTagFunc() func(string, fields.FieldInterface) {
	return f.validTagFn
}

func (f *Form) Init(c *config.Config, model ...interface{}) *Form {
	if c == nil {
		c = &config.Config{}
	}
	f.config = c
	if len(c.Theme) == 0 {
		c.Theme = common.BASE
	}
	f.Style = c.Theme
	if len(c.Template) == 0 {
		c.Template = common.TmplDir(f.Style) + `/baseform.html`
		//c.Template = common.TmplDir(f.Style) + `/allfields.html`
	}
	tpf := c.Template
	tmpl, err := common.GetOrSetCachedTemplate(tpf, func() (*template.Template, error) {
		return common.ParseFiles(common.LookupPath(c.Template))
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	f.template = tmpl
	f.Method = c.Method
	f.Action = template.HTML(c.Action)
	if len(model) > 0 {
		f.Model = model[0]
	}
	return f
}

func (f *Form) Valid(onlyCheckFields ...string) error {
	if f.Model == nil {
		return nil
	}
	passed, err := f.Validate().Valid(f.Model, onlyCheckFields...)
	if err != nil {
		return err
	}
	if !passed { // validation does not pass
		for index, validErr := range f.Validate().Errors {
			if index == 0 {
				err = validErr
			}
			f.Field(validErr.Field).AddError(f.labelFn(validErr.Message))
		}
	}
	return err
}

func (f *Form) InsertErrors() *Form {
	if f.valid != nil && f.valid.HasError() {
		for _, err := range f.valid.Errors {
			f.Field(err.Field).AddError(f.labelFn(err.Message))
		}
	}
	return f
}

func (f *Form) Error() (err *validation.ValidationError) {
	if f.valid != nil && f.valid.HasError() {
		err = f.valid.Errors[0]
	}
	return
}

func (f *Form) Errors() (errs []*validation.ValidationError) {
	if f.valid != nil && f.valid.HasError() {
		errs = f.valid.Errors
	}
	return
}

func (f *Form) HasError() bool {
	return f.valid != nil && f.valid.HasError()
}

func (f *Form) HasErrors() bool {
	return f.valid != nil && f.valid.HasErrors()
}

func (f *Form) AddBeforeRender(fn func()) *Form {
	if fn == nil {
		return f
	}
	f.beforeRender = append(f.beforeRender, fn)
	return f
}

func (f *Form) Validate() *validation.Validation {
	if f.valid == nil {
		f.valid = &validation.Validation{}
	}
	return f.valid
}

func (f *Form) SetModel(m interface{}) *Form {
	f.Model = m
	return f
}

func (f *Form) ParseModel(model ...interface{}) *Form {
	var m interface{}
	if len(model) > 0 {
		m = model[0]
	}
	if m == nil {
		m = f.Model
	}
	flist, fsort := f.unWindStructure(m, ``)
	for _, v := range flist {
		f.Elements(v.(config.FormElement))
	}
	if len(fsort) > 0 {
		f.Sort(fsort)
	}
	return f
}

func (f *Form) AddButton(tmpl string, args ...string) *Form {
	btnFields := make([]fields.FieldInterface, 0)
	if len(args) < 1 {
		btnFields = append(btnFields, fields.SubmitButton("submit", f.labelFn("Submit")))
		btnFields = append(btnFields, fields.ResetButton("reset", f.labelFn("Reset")))
	} else {
		for _, field := range args {
			switch field {
			case `submit`:
				btnFields = append(btnFields, fields.SubmitButton("submit", f.labelFn("Submit")))
			case `reset`:
				btnFields = append(btnFields, fields.ResetButton("reset", f.labelFn("Reset")))
			default:
				btnFields = append(btnFields, fields.Button(field, f.labelFn(strings.ToTitle(field))))
			}
		}
	}
	if len(tmpl) == 0 {
		tmpl = `fieldset_buttons`
	}
	f.Elements(f.NewFieldSet("_button_group", "", btnFields...).SetTemplate(tmpl))
	return f
}

func (f *Form) GenChoicesForField(name string, lenType interface{}, fnType interface{}) *Form {
	f.Field(name).SetChoices(f.GenChoices(lenType, fnType))
	return f
}

// GenChoices generate choices
// type Data struct{
// 	ID string
// 	Name string
// }
// data:=[]*Data{
// 	&Data{ID:"a",Name:"One"},
// 	&Data{ID:"b",Name:"Two"},
// }
// GenChoices(len(data), func(index int) (string, string, bool){
// 	return data[index].ID,data[index].Name,false
// })
// or
// GenChoices(map[string]int{
// 	"":len(data),
// }, func(group string,index int) (string, string, bool){
// 	return data[index].ID,data[index].Name,false
// })
func (f *Form) GenChoices(lenType interface{}, fnType interface{}) interface{} {
	return GenChoices(lenType, fnType)
}

func (form *Form) unWindStructure(m interface{}, baseName string) ([]interface{}, string) {
	t := reflect.TypeOf(m)
	v := reflect.ValueOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}
	var (
		fieldList []interface{}
		fieldSort string
	)
	fieldSetList := map[string]*FieldSetType{}
	fieldSetSort := map[string]string{}
	for i := 0; i < t.NumField(); i++ {
		if !v.Field(i).CanInterface() {
			continue
		}
		options := make(map[string]struct{})
		tag, tagf := common.Tag(t, t.Field(i), "form_options")
		if len(tag) > 0 {
			var optionsArr []string
			if tagf != nil {
				cached := tagf.Parsed("form_options", func() interface{} {
					return strings.Split(common.TagVal(t, i, "form_options"), ";")
				})
				optionsArr = cached.([]string)
			}
			for _, opt := range optionsArr {
				if len(opt) > 0 {
					options[opt] = struct{}{}
				}
			}
		}
		if _, ok := options["-"]; ok {
			continue
		}
		widget := common.TagVal(t, i, "form_widget")
		var f fields.FieldInterface
		var fName string
		if len(baseName) == 0 {
			fName = t.Field(i).Name
		} else {
			fName = strings.Join([]string{baseName, t.Field(i).Name}, ".")
		}
		useFieldValue := !form.IsOmit(fName)
		//fmt.Println(fName, t.Field(i).Type.String(), t.Field(i).Type.Kind())
		switch widget {
		case "color", "email", "file", "image", "month", "search", "tel", "url", "week":
			f = fields.TextFieldFromInstance(v, t, i, fName, useFieldValue, widget)
		case "text":
			f = fields.TextFieldFromInstance(v, t, i, fName, useFieldValue)
		case "hidden":
			f = fields.HiddenFieldFromInstance(v, t, i, fName, useFieldValue)
		case "textarea":
			f = fields.TextAreaFieldFromInstance(v, t, i, fName, useFieldValue)
		case "password":
			f = fields.PasswordFieldFromInstance(v, t, i, fName, useFieldValue)
		case "select":
			f = fields.SelectFieldFromInstance(v, t, i, fName, useFieldValue, options, form.labelFn)
		case "date":
			f = fields.DateFieldFromInstance(v, t, i, fName, useFieldValue)
		case "datetime":
			f = fields.DatetimeFieldFromInstance(v, t, i, fName, useFieldValue)
		case "time":
			f = fields.TimeFieldFromInstance(v, t, i, fName, useFieldValue)
		case "number":
			f = fields.NumberFieldFromInstance(v, t, i, fName, useFieldValue)
		case "range":
			f = fields.RangeFieldFromInstance(v, t, i, fName, useFieldValue)
		case "radio":
			f = fields.RadioFieldFromInstance(v, t, i, fName, useFieldValue, form.labelFn)
		case "checkbox":
			f = fields.CheckboxFieldFromInstance(v, t, i, fName, useFieldValue, form.labelFn)
		case "static":
			f = fields.StaticFieldFromInstance(v, t, i, fName, useFieldValue)
		default:
			switch t.Field(i).Type.String() {
			case "string":
				f = fields.TextFieldFromInstance(v, t, i, fName, useFieldValue)
			case "bool":
				f = fields.CheckboxFromInstance(v, t, i, fName, useFieldValue, options)
			case "time.Time":
				f = fields.DatetimeFieldFromInstance(v, t, i, fName, useFieldValue)
			case "int", "int64", "float", "float32", "float64":
				f = fields.NumberFieldFromInstance(v, t, i, fName, useFieldValue)
			case "struct":
				fl, fs := form.unWindStructure(v.Field(i).Interface(), fName)
				if len(fs) > 0 {
					if len(fieldSort) == 0 {
						fieldSort = fs
					} else {
						fieldSort += "," + fs
					}
				}
				fieldList = append(fieldList, fl...)
				f = nil
			default:
				if t.Field(i).Type.Kind() == reflect.Struct ||
					(t.Field(i).Type.Kind() == reflect.Ptr &&
						t.Field(i).Type.Elem().Kind() == reflect.Struct) {
					fl, fs := form.unWindStructure(v.Field(i).Interface(), fName)
					if len(fs) > 0 {
						if len(fieldSort) == 0 {
							fieldSort = fs
						} else {
							fieldSort += "," + fs
						}
					}
					fieldList = append(fieldList, fl...)
					f = nil
				} else {
					f = fields.TextFieldFromInstance(v, t, i, fName, useFieldValue)
				}
			}
		}
		if f != nil {
			label := common.TagVal(t, i, "form_label")
			if len(label) == 0 {
				label = strings.Title(t.Field(i).Name)
			} else if label != `-` {
				label = form.labelFn(label)
			} else {
				label = ``
			}
			f.SetLabel(label)

			params := common.TagVal(t, i, "form_params")
			if len(params) > 0 {
				if paramsMap, err := url.ParseQuery(params); err == nil {
					for k, v := range paramsMap {
						if k == "placeholder" || k == "title" {
							v[0] = form.labelFn(v[0])
						}
						f.SetParam(k, v[0])
					}
				} else {
					fmt.Println(err)
				}
			}
			valid := common.TagVal(t, i, "valid")
			if len(valid) > 0 {
				form.validTagFn(valid, f)
			}
			fieldsetLabel := common.TagVal(t, i, "form_fieldset") // label;name or name
			fieldsort := common.TagVal(t, i, "form_sort")         // 1 ( or other number ) or "last"
			if len(fieldsetLabel) > 0 {
				fieldsets := strings.SplitN(fieldsetLabel, ";", 2)
				var fieldsetName string
				switch len(fieldsets) {
				case 1:
					fieldsetName = fieldsets[0]
				case 2:
					fieldsetLabel = fieldsets[0]
					fieldsetName = fieldsets[1]
				}
				fieldsetLabel = form.labelFn(fieldsetLabel)
				f.SetData("container", "fieldset")
				if _, ok := fieldSetList[fieldsetName]; !ok {
					fieldSetList[fieldsetName] = form.NewFieldSet(fieldsetName, fieldsetLabel, f)
				} else {
					fieldSetList[fieldsetName].Elements(f)
				}
				if len(fieldsort) > 0 {
					if _, ok := fieldSetSort[fieldsetName]; !ok {
						fieldSetSort[fieldsetName] = fName + ":" + fieldsort
					} else {
						fieldSetSort[fieldsetName] += "," + fName + ":" + fieldsort
					}
				}
			} else {
				fieldList = append(fieldList, f)
				if len(fieldsort) > 0 {
					if len(fieldSort) == 0 {
						fieldSort = fName + ":" + fieldsort
					} else {
						fieldSort += "," + fName + ":" + fieldsort
					}
				}
			}
		}
	}
	for _, v := range fieldSetList {
		if s, ok := fieldSetSort[v.OriginalName()]; ok {
			v.Sort(s)
		}
		fieldList = append(fieldList, v)
	}
	return fieldList, fieldSort
}

func (f *Form) SetData(key string, value interface{}) {
	f.AppendData[key] = value
}

func (f *Form) Data() map[string]interface{} {
	if len(f.data) > 0 {
		return f.data
	}
	safeParams := make(common.HTMLAttributes)
	safeParams.FillFromStringMap(f.Params)
	f.data = map[string]interface{}{
		"container": "",
		"fields":    f.FieldList,
		"classes":   f.Class,
		"id":        f.ID,
		"params":    safeParams,
		"css":       f.CSS,
		"method":    f.Method,
		"action":    f.Action,
	}
	for k, v := range f.AppendData {
		f.data[k] = v
	}
	return f.data
}

func (f *Form) runBefore() {
	for _, fn := range f.beforeRender {
		fn()
	}
}

func (f *Form) render() string {
	f.runBefore()
	buf := bytes.NewBuffer(nil)
	err := f.template.Execute(buf, f.Data())
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

// Render executes the internal template and renders the form, returning the result as a template.HTML object embeddable
// in any other template.
func (f *Form) Render() template.HTML {
	return template.HTML(f.render())
}

func (f *Form) Html(value interface{}) template.HTML {
	return template.HTML(fmt.Sprintf("%v", value))
}

func (f *Form) String() string {
	return f.render()
}

// Elements adds the provided elements to the form.
func (f *Form) Elements(elems ...config.FormElement) {
	for _, e := range elems {
		switch v := e.(type) {
		case fields.FieldInterface:
			f.addField(v)
		case *FieldSetType:
			f.addFieldSet(v)
		case *LangSetType:
			f.addLangSet(v)
		}
	}
}

func (f *Form) addField(field fields.FieldInterface) *Form {
	field.SetStyle(f.Style)
	f.FieldList = append(f.FieldList, field)
	f.fieldMap[field.OriginalName()] = len(f.FieldList) - 1
	return f
}

func (f *Form) addFieldSet(fs *FieldSetType) *Form {
	for _, v := range fs.FieldList {
		v.SetStyle(f.Style)
		v.SetData("container", "fieldset")
		f.containerMap[v.OriginalName()] = fs.OriginalName()
	}
	f.FieldList = append(f.FieldList, fs)
	f.fieldMap[fs.OriginalName()] = len(f.FieldList) - 1
	return f
}

func (f *Form) addLangSet(fs *LangSetType) *Form {
	for _, v := range fs.fieldMap {
		v.SetData("container", "langset")
		f.containerMap[v.OriginalName()] = fs.OriginalName()
	}
	f.FieldList = append(f.FieldList, fs)
	f.fieldMap[fs.OriginalName()] = len(f.FieldList) - 1
	return f
}

// RemoveElement removes an element (identified by name) from the Form.
func (f *Form) RemoveElement(name string) *Form {
	ind, ok := f.fieldMap[name]
	if !ok {
		return f
	}
	delete(f.fieldMap, name)
	f.FieldList = append(f.FieldList[:ind], f.FieldList[ind+1:]...)
	return f
}

// AddClass associates the provided class to the Form.
func (f *Form) AddClass(class string) *Form {
	f.Class.Add(class)
	return f
}

// RemoveClass removes the given class (if present) from the Form.
func (f *Form) RemoveClass(class string) *Form {
	f.Class.Remove(class)
	return f
}

// SetID set the given id to the form.
func (f *Form) SetID(id string) *Form {
	f.ID = id
	return f
}

// SetParam adds the given key-value pair to form parameters list.
func (f *Form) SetParam(key, value string) *Form {
	switch key {
	case `class`:
		f.AddClass(value)
	case `action`:
		f.Action = template.HTML(value)
	default:
		f.Params[key] = value
	}
	return f
}

// DeleteParam removes the parameter identified by key from form parameters list.
func (f *Form) DeleteParam(key string) *Form {
	delete(f.Params, key)
	return f
}

// AddCSS add a CSS value (in the form of option-value - e.g.: border - auto) to the form.
func (f *Form) AddCSS(key, value string) *Form {
	f.CSS[key] = value
	return f
}

// RemoveCSS removes CSS style from the form.
func (f *Form) RemoveCSS(key string) *Form {
	delete(f.CSS, key)
	return f
}

// Field returns the field identified by name. It returns an empty field if it is missing.
func (f *Form) Field(name string) fields.FieldInterface {
	ind, ok := f.fieldMap[name]
	if !ok {
		return &fields.Field{}
	}
	switch v := f.FieldList[ind].(type) {
	case fields.FieldInterface:
		return v
	case *FieldSetType:
		if v, ok := f.containerMap[name]; ok {
			return f.FieldSet(v).Field(name)
		}
	case *LangSetType:
		if v, ok := f.containerMap[name]; ok {
			return f.LangSet(v).Field(name)
		}
	}
	if f.debug {
		fmt.Printf("[Form] Not found field: %s , but has: %#v\n", name, f.fieldMap)
	}
	return &fields.Field{}
}

// Fields returns all field
func (f *Form) Fields() []config.FormElement {
	return f.FieldList
}

// LangSet returns the fieldset identified by name. It returns an empty field if it is missing.
func (f *Form) LangSet(name string) *LangSetType {
	ind, ok := f.fieldMap[name]
	if !ok {
		return &LangSetType{}
	}
	switch v := f.FieldList[ind].(type) {
	case *LangSetType:
		return v
	default:
		return &LangSetType{}
	}
}

func (f *Form) NewLangSet(name string, langs []*config.Language) *LangSetType {
	return LangSet(name, f.Style, langs...)
}

// FieldSet returns the fieldset identified by name. It returns an empty field if it is missing.
func (f *Form) FieldSet(name string) *FieldSetType {
	ind, ok := f.fieldMap[name]
	if !ok {
		return &FieldSetType{}
	}
	switch v := f.FieldList[ind].(type) {
	case *FieldSetType:
		return v
	default:
		return &FieldSetType{}
	}
}

// NewFieldSet creates and returns a new FieldSetType with the given name and list of fields.
// Every method for FieldSetType objects returns the object itself, so that call can be chained.
func (f *Form) NewFieldSet(name string, label string, elems ...fields.FieldInterface) *FieldSetType {
	return FieldSet(name, label, f.Style, elems...)
}

// SortAll SortAll("field1,field2") or SortAll("field1","field2")
func (f *Form) SortAll(sortList ...string) *Form {
	elem := f.FieldList
	size := len(elem)
	f.FieldList = make([]config.FormElement, size)
	var sortSlice []string
	if len(sortList) == 1 {
		sortSlice = strings.Split(sortList[0], ",")
	} else {
		sortSlice = sortList
	}
	for k, fieldName := range sortSlice {
		if oldIndex, ok := f.fieldMap[fieldName]; ok {
			f.FieldList[k] = elem[oldIndex]
			f.fieldMap[fieldName] = k
		}
	}
	return f
}

// Sort Sort("field1:1,field2:2") or Sort("field1:1","field2:2")
func (f *Form) Sort(sortList ...string) *Form {
	size := len(f.FieldList)
	var sortSlice []string
	if len(sortList) == 1 {
		sortSlice = strings.Split(sortList[0], ",")
	} else {
		sortSlice = sortList
	}
	var index int
	endIdx := size - 1

	for _, nameIndex := range sortSlice {
		ni := strings.Split(nameIndex, ":")
		fieldName := ni[0]
		if len(ni) > 1 {
			if ni[1] == "last" {
				index = endIdx
			} else if idx, err := strconv.Atoi(ni[1]); err != nil {
				continue
			} else {
				if idx >= 0 {
					index = idx
				} else {
					index = size + idx
				}
			}
		}
		if oldIndex, ok := f.fieldMap[fieldName]; ok {
			if oldIndex != index && size > index {
				f.sortFields(index, oldIndex, endIdx, size)
			}
		}
		index++
	}
	return f
}

func (f *Form) Sort2Last(fieldsName ...string) *Form {
	size := len(f.FieldList)
	endIdx := size - 1
	index := endIdx
	for n := len(fieldsName) - 1; n >= 0; n-- {
		fieldName := fieldsName[n]
		if oldIndex, ok := f.fieldMap[fieldName]; ok {
			if oldIndex != index && index >= 0 {
				f.sortFields(index, oldIndex, endIdx, size)
			}
		}
		index--
	}
	return f
}

func (f *Form) sortFields(index, oldIndex, endIdx, size int) {

	var newFields []config.FormElement
	oldFields := make([]config.FormElement, size)
	copy(oldFields, f.FieldList)
	var min, max int
	if index > oldIndex {
		//[ ][I][ ][ ][ ][ ] I:oldIndex=1
		//[ ][ ][ ][ ][I][ ] I:index=4
		if oldIndex > 0 {
			newFields = oldFields[0:oldIndex]
		}
		newFields = append(newFields, oldFields[oldIndex+1:index+1]...)
		newFields = append(newFields, f.FieldList[oldIndex])
		if index+1 <= endIdx {
			newFields = append(newFields, f.FieldList[index+1:]...)
		}
		min = oldIndex
		max = index
	} else {
		//[ ][ ][ ][ ][I][ ] I:oldIndex=4
		//[ ][I][ ][ ][ ][ ] I:index=1
		if index > 0 {
			newFields = oldFields[0:index]
		}
		newFields = append(newFields, oldFields[oldIndex])
		newFields = append(newFields, f.FieldList[index:oldIndex]...)
		if oldIndex+1 <= endIdx {
			newFields = append(newFields, f.FieldList[oldIndex+1:]...)
		}
		min = index
		max = oldIndex
	}
	for i := min; i <= max; i++ {
		f.fieldMap[newFields[i].OriginalName()] = i
	}
	f.FieldList = newFields
}
