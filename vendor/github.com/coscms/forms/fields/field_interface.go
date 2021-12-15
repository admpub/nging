package fields

import (
	"html/template"

	"github.com/coscms/forms/config"
)

// FieldInterface defines the interface an object must implement to be used in a form. Every method returns a FieldInterface object
// to allow methods chaining.
type FieldInterface interface {
	Name() string
	OriginalName() string
	SetName(string)
	SetLabelCols(cols int)
	SetFieldCols(cols int)
	Render() template.HTML
	AddClass(class string) FieldInterface
	RemoveClass(class string) FieldInterface
	AddTag(class string) FieldInterface
	RemoveTag(class string) FieldInterface
	SetID(id string) FieldInterface
	SetParam(key string, value interface{}) FieldInterface
	DeleteParam(key string) FieldInterface
	AddCSS(key, value string) FieldInterface
	RemoveCSS(key string) FieldInterface
	SetTheme(theme string) FieldInterface
	SetLabel(label string) FieldInterface
	AddLabelClass(class string) FieldInterface
	RemoveLabelClass(class string) FieldInterface
	SetValue(value string) FieldInterface
	Disabled() FieldInterface
	Enabled() FieldInterface
	SetTemplate(tmpl string, theme ...string) FieldInterface
	SetHelptext(text string) FieldInterface
	AddError(err string) FieldInterface
	MultipleChoice() FieldInterface
	SingleChoice() FieldInterface
	AddSelected(opt ...string) FieldInterface
	SetSelected(opt ...string) FieldInterface
	RemoveSelected(opt string) FieldInterface
	AddChoice(key, value interface{}, checked ...bool) FieldInterface
	SetChoices(choices interface{}, saveIndex ...bool) FieldInterface
	SetText(text string) FieldInterface
	SetData(key string, value interface{})
	Data() map[string]interface{}
	String() string
	SetLang(lang string)
	Lang() string
	Clone() config.FormElement

	Element() *config.Element
}
