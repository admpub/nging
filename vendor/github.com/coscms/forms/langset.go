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
	"bytes"
	"strconv"
	"strings"

	"html/template"

	"github.com/coscms/forms/common"
	conf "github.com/coscms/forms/config"
	"github.com/coscms/forms/fields"
)

// LangSetType is a collection of fields grouped within a form.
type LangSetType struct {
	Languages    []*conf.Language
	LangMap      map[string]int              //{"zh-CN":1}
	FieldMap     map[string]conf.FormElement //{"zh-CN:title":0x344555}
	ContainerMap map[string]string           //{"name":"fieldset's name"}
	CurrName     string
	OrigName     string
	Tmpl         string
	Params       map[string]interface{}
	Tags         common.HTMLAttrValues
	AppendData   map[string]interface{}
	SetAlone     bool
	FormStyle    string
	data         map[string]interface{}
}

func (f *LangSetType) SetName(name string) {
	f.CurrName = name
}

func (f *LangSetType) OriginalName() string {
	return f.OrigName
}

func (f *LangSetType) SetLang(lang string) {
}

func (f *LangSetType) Lang() string {
	return ``
}

func (f *LangSetType) Clone() conf.FormElement {
	fc := *f
	return &fc
}

func (f *LangSetType) AddLanguage(language *conf.Language) {
	f.LangMap[language.ID] = len(f.Languages)
	f.Languages = append(f.Languages, language)
}

func (f *LangSetType) Language(lang string) *conf.Language {
	if ind, ok := f.LangMap[lang]; ok {
		return f.Languages[ind]
	}
	return nil
}

func (f *LangSetType) SetData(key string, value interface{}) {
	f.AppendData[key] = value
}

func (f *LangSetType) Data() map[string]interface{} {
	if len(f.data) > 0 {
		return f.data
	}
	safeParams := make(map[template.HTMLAttr]interface{})
	for k, v := range f.Params {
		safeParams[template.HTMLAttr(k)] = v
	}
	f.data = map[string]interface{}{
		"container": "langset",
		"params":    safeParams,
		"tags":      f.Tags,
		"langs":     f.Languages,
		"name":      f.CurrName,
	}
	for k, v := range f.AppendData {
		f.data[k] = v
	}
	return f.data
}

func (f *LangSetType) render() string {
	buf := bytes.NewBuffer(nil)
	tpf := common.TmplDir(f.FormStyle) + "/" + f.FormStyle + "/" + f.Tmpl + ".html"
	var err error
	tpl, ok := common.CachedTemplate(tpf)
	if !ok {
		tpl, err = common.ParseFiles(common.CreateUrl(tpf))
		if err != nil {
			return err.Error()
		}
		common.SetCachedTemplate(tpf, tpl)
	}
	err = tpl.Execute(buf, f.Data())
	if err != nil {
		return err.Error()
	}
	return buf.String()
}

// Render translates a FieldSetType into HTML code and returns it as a template.HTML object.
func (f *LangSetType) Render() template.HTML {
	return template.HTML(f.render())
}

func (f *LangSetType) String() string {
	return f.render()
}

func (f *LangSetType) SetTmpl(tmpl string) *LangSetType {
	f.Tmpl = tmpl
	return f
}

// FieldSet creates and returns a new FieldSetType with the given name and list of fields.
// Every method for FieldSetType objects returns the object itself, so that call can be chained.
func LangSet(name string, style string, languages ...*conf.Language) *LangSetType {
	ret := &LangSetType{
		Languages:    languages,
		LangMap:      map[string]int{},
		ContainerMap: make(map[string]string),
		FieldMap:     make(map[string]conf.FormElement),
		CurrName:     name,
		OrigName:     name,
		Tmpl:         "langset",
		Params:       map[string]interface{}{},
		Tags:         common.HTMLAttrValues{},
		AppendData:   map[string]interface{}{},
		FormStyle:    style,
	}
	for i, language := range languages {
		ret.LangMap[language.ID] = i
	}
	return ret
}

//SortAll("field1,field2") or SortAll("field1","field2")
func (f *LangSetType) SortAll(sortList ...string) *LangSetType {
	elem := f.Languages
	size := len(elem)
	f.Languages = make([]*conf.Language, size)
	var sortSlice []string
	if len(sortList) == 1 {
		sortSlice = strings.Split(sortList[0], ",")
	} else {
		sortSlice = sortList
	}
	for k, fieldName := range sortSlice {
		if oldIndex, ok := f.LangMap[fieldName]; ok {
			f.Languages[k] = elem[oldIndex]
			f.LangMap[fieldName] = k
		}
	}
	return f
}

// Elements adds the provided elements to the langset.
func (f *LangSetType) Elements(elems ...conf.FormElement) {
	for _, e := range elems {
		switch v := e.(type) {
		case fields.FieldInterface:
			f.addField(v)
		case *FieldSetType:
			f.addFieldSet(v)
		}
	}
}

func (f *LangSetType) addField(field fields.FieldInterface) *LangSetType {
	field.SetStyle(f.FormStyle)
	if f.SetAlone {
		if ind, ok := f.LangMap[field.Lang()]; ok {
			field.SetLang(f.Languages[ind].ID)
			field.SetName(f.Languages[ind].Name(field.OriginalName()))
			f.Languages[ind].AddField(field)
			f.FieldMap[field.Lang()+`:`+field.OriginalName()] = field
		}
		return f
	}
	for k, language := range f.Languages {
		f.LangMap[language.ID] = k
		if k == 0 {
			field.SetLang(language.ID)
			field.SetName(language.Name(field.OriginalName()))
			language.AddField(field)
			f.FieldMap[field.Lang()+`:`+field.OriginalName()] = field
			continue
		}
		fieldCopy := field.Clone()
		fieldCopy.SetLang(language.ID)
		fieldCopy.SetName(language.Name(fieldCopy.OriginalName()))
		language.AddField(fieldCopy)
		f.FieldMap[fieldCopy.Lang()+`:`+fieldCopy.OriginalName()] = fieldCopy
	}
	return f
}

func (f *LangSetType) addFieldSet(fs *FieldSetType) *LangSetType {
	if f.SetAlone {
		if ind, ok := f.LangMap[fs.Lang()]; ok {
			for _, v := range fs.FieldList {
				v.SetStyle(f.FormStyle)
				v.SetData("container", "langset")
				v.SetLang(f.Languages[ind].ID)
				v.SetName(f.Languages[ind].Name(v.OriginalName()))
				key := v.Lang() + `:` + v.OriginalName()
				f.FieldMap[key] = v
				f.ContainerMap[key] = fs.OriginalName()
			}
			fs.SetLang(f.Languages[ind].ID)
			fs.SetName(f.Languages[ind].Name(fs.OriginalName()))
			f.Languages[ind].AddField(fs)
			f.FieldMap[fs.Lang()+`:`+fs.OriginalName()] = fs
		}
		return f
	}
	for k, language := range f.Languages {
		f.LangMap[language.ID] = k
		if k == 0 {
			for _, v := range fs.FieldList {
				v.SetLang(language.ID)
				v.SetStyle(f.FormStyle)
				v.SetData("container", "langset")
				key := v.Lang() + `:` + v.OriginalName()
				f.FieldMap[key] = v
				f.ContainerMap[key] = fs.OriginalName()
				v.SetName(language.Name(v.OriginalName()))
			}
			fs.SetLang(language.ID)
			fs.SetName(language.Name(fs.OriginalName()))
			language.AddField(fs)
			f.FieldMap[fs.Lang()+`:`+fs.OriginalName()] = fs
			continue
		}
		fsCopy := fs.Clone().(*FieldSetType)
		fsCopy.FieldList = make([]fields.FieldInterface, len(fs.FieldList))
		for kk, v := range fs.FieldList {
			fieldCopy := v.Clone().(fields.FieldInterface)
			fieldCopy.SetLang(language.ID)
			fieldCopy.SetName(language.Name(fieldCopy.OriginalName()))
			key := fieldCopy.Lang() + `:` + fieldCopy.OriginalName()
			f.FieldMap[key] = fieldCopy
			f.ContainerMap[key] = fs.OriginalName()
			fsCopy.FieldList[kk] = fieldCopy
		}
		fsCopy.SetLang(language.ID)
		fsCopy.SetName(language.Name(fsCopy.OriginalName()))
		language.AddField(fsCopy)
		f.FieldMap[fsCopy.Lang()+`:`+fsCopy.OriginalName()] = fsCopy
	}
	return f
}

// Field returns the field identified by name. It returns an empty field if it is missing.
// param format: "language:name"
func (f *LangSetType) Field(name string) fields.FieldInterface {
	field, ok := f.FieldMap[name]
	if !ok {
		return &fields.Field{}
	}
	switch v := field.(type) {
	case fields.FieldInterface:
		return v
	case *FieldSetType:
		if v, ok := f.ContainerMap[name]; ok {
			r := strings.SplitN(name, `:`, 2)
			switch len(r) {
			case 2:
				return f.FieldSet(v).Field(r[1])
			case 1:
				return f.FieldSet(v).Field(r[0])
			}
		}
	}
	return &fields.Field{}
}

// FieldSet returns the fieldset identified by name.
// param format: "language:name"
func (f *LangSetType) FieldSet(name string) *FieldSetType {
	field, ok := f.FieldMap[name]
	if !ok {
		return &FieldSetType{}
	}
	switch v := field.(type) {
	case *FieldSetType:
		return v
	default:
		return &FieldSetType{}
	}
}

// NewFieldSet creates and returns a new FieldSetType with the given name and list of fields.
// Every method for FieldSetType objects returns the object itself, so that call can be chained.
func (f *LangSetType) NewFieldSet(name string, label string, elems ...fields.FieldInterface) *FieldSetType {
	return FieldSet(name, label, f.FormStyle, elems...)
}

//Sort Sort("field1:1,field2:2") or Sort("field1:1","field2:2")
func (f *LangSetType) Sort(sortList ...string) *LangSetType {
	size := len(f.Languages)
	endIdx := size - 1
	var sortSlice []string
	if len(sortList) == 1 {
		sortSlice = strings.Split(sortList[0], ",")
	} else {
		sortSlice = sortList
	}
	var index int
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
					index = endIdx + idx
				}
			}
		}
		if oldIndex, ok := f.LangMap[fieldName]; ok {
			if oldIndex != index && size > index {
				f.sortFields(index, oldIndex, endIdx, size)
			}
		}
		index++
	}
	return f
}

func (f *LangSetType) Sort2Last(fieldsName ...string) *LangSetType {
	size := len(f.Languages)
	endIdx := size - 1
	index := endIdx
	for n := len(fieldsName) - 1; n >= 0; n-- {
		fieldName := fieldsName[n]
		if oldIndex, ok := f.LangMap[fieldName]; ok {
			if oldIndex != index && index >= 0 {
				f.sortFields(index, oldIndex, endIdx, size)
			}
		}
		index--
	}
	return f
}

// Name returns the name of the langset.
func (f *LangSetType) Name() string {
	return f.CurrName
}

// SetParam saves the provided param for the langset.
func (f *LangSetType) SetParam(k string, v interface{}) *LangSetType {
	f.Params[k] = v
	return f
}

// DeleteParam removes the provided param from the langset, if it was present. Nothing is done if it was not originally present.
func (f *LangSetType) DeleteParam(k string) *LangSetType {
	delete(f.Params, k)
	return f
}

// AddTag adds a no-value parameter (e.g.: "disabled", "checked") to the langset.
func (f *LangSetType) AddTag(tag string) *LangSetType {
	f.Tags.Add(tag)
	return f
}

// RemoveTag removes a tag from the langset, if it was present.
func (f *LangSetType) RemoveTag(tag string) *LangSetType {
	f.Tags.Remove(tag)
	return f
}

// Disable adds tag "disabled" to the langset, making it unresponsive in some environment (e.g.: Bootstrap).
func (f *LangSetType) Disable() *LangSetType {
	f.AddTag("disabled")
	return f
}

// Enable removes tag "disabled" from the langset, making it responsive.
func (f *LangSetType) Enable() *LangSetType {
	f.RemoveTag("disabled")
	return f
}

func (f *LangSetType) sortFields(index, oldIndex, endIdx, size int) {
	var newFields []*conf.Language
	oldFields := make([]*conf.Language, size)
	copy(oldFields, f.Languages)
	var min, max int
	if index > oldIndex {
		//[ ][I][ ][ ][ ][ ] I:oldIndex=1
		//[ ][ ][ ][ ][I][ ] I:index=4
		if oldIndex > 0 {
			newFields = oldFields[0:oldIndex]
		}
		newFields = append(newFields, oldFields[oldIndex+1:index+1]...)
		newFields = append(newFields, f.Languages[oldIndex])
		if index+1 <= endIdx {
			newFields = append(newFields, f.Languages[index+1:]...)
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
		newFields = append(newFields, f.Languages[index:oldIndex]...)
		if oldIndex+1 <= endIdx {
			newFields = append(newFields, f.Languages[oldIndex+1:]...)
		}
		min = index
		max = oldIndex
	}
	for i := min; i <= max; i++ {
		f.LangMap[newFields[i].ID] = i
	}
	f.Languages = newFields
}
