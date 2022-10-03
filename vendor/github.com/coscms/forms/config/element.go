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
	elements := make([]*Element, len(e.Elements))
	languages := make([]*Language, len(e.Languages))
	choices := make([]*Choice, len(e.Choices))
	for index, elem := range e.Elements {
		elements[index] = elem.Clone()
	}
	for index, value := range e.Languages {
		languages[index] = value.Clone()
	}
	for index, value := range e.Choices {
		choices[index] = value.Clone()
	}
	r := &Element{
		ID:         e.ID,
		Type:       e.Type,
		Name:       e.Name,
		Label:      e.Label,
		LabelCols:  e.LabelCols,
		FieldCols:  e.FieldCols,
		Value:      e.Value,
		HelpText:   e.HelpText,
		Template:   e.Template,
		Valid:      e.Valid,
		Attributes: make([][]string, len(e.Attributes)),
		Choices:    choices,
		Elements:   elements,
		Format:     e.Format,
		Languages:  languages,
		Data:       map[string]interface{}{},
	}
	for k, v := range e.Data {
		r.Data[k] = v
	}
	for k, v := range e.Attributes {
		cv := make([]string, len(v))
		copy(cv, v)
		r.Attributes[k] = cv
	}
	return r
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
