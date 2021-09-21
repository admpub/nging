package config

import "strings"

type Element struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Label      string                 `json:"label"`
	LabelCols  int                    `json:"labelCols,omitempty"`
	FieldCols  int                    `json:"fieldCols,omitempty"`
	Value      string                 `json:"value"`
	HelpText   string                 `json:"helpText"`
	Template   string                 `json:"template"`
	Valid      string                 `json:"valid"`
	Attributes [][]string             `json:"attributes"`
	Choices    []*Choice              `json:"choices"`
	Elements   []*Element             `json:"elements"`
	Format     string                 `json:"format"`
	Languages  []*Language            `json:"languages,omitempty"`
	Data       map[string]interface{} `json:"data,omitempty"`
}

func (e *Element) Clone() *Element {
	r := *e
	return &r
}

func (e *Element) HasAttr(attrs ...string) bool {
	mk := map[string]struct{}{}
	for _, attr := range attrs {
		mk[strings.ToLower(attr)] = struct{}{}
	}
	for _, v := range e.Attributes {
		if len(v) == 0 || len(v[0]) == 0 {
			continue
		}
		v[0] = strings.ToLower(v[0])
		if _, ok := mk[v[0]]; ok {
			return true
		}
	}
	return false
}

func (e *Element) AddElement(elements ...*Element) *Element {
	e.Elements = append(e.Elements, elements...)
	return e
}

func (e *Element) AddLanguage(languages ...*Language) *Element {
	e.Languages = append(e.Languages, languages...)
	return e
}

func (e *Element) AddAttribute(attributes ...string) *Element {
	e.Attributes = append(e.Attributes, attributes)
	return e
}

func (e *Element) AddChoice(choices ...*Choice) *Element {
	e.Choices = append(e.Choices, choices...)
	return e
}

func (e *Element) Set(name string, value interface{}) *Element {
	if e.Data == nil {
		e.Data = map[string]interface{}{}
	}
	e.Data[name] = value
	return e
}
