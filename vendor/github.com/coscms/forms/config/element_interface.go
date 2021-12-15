package config

import "html/template"

// FormElement interface defines a form object (usually a Field or a FieldSet) that can be rendered as a template.HTML object.
type FormElement interface {
	Render() template.HTML
	Name() string
	OriginalName() string
	SetName(string)
	String() string
	SetData(key string, value interface{})
	Data() map[string]interface{}
	SetLang(lang string)
	Lang() string
	Clone() FormElement
}
