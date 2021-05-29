package config

import (
	"fmt"
)

func NewLanguage(lang, label, namefmt string) *Language {
	return &Language{
		ID:         lang,
		Label:      label,
		NameFormat: namefmt,
		fields:     make([]FormElement, 0),
		fieldMap:   make(map[string]int),
	}
}

type Language struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	NameFormat string `json:"nameFormat"`
	fields     []FormElement
	fieldMap   map[string]int
}

func (l *Language) Name(name string) string {
	if len(l.NameFormat) == 0 {
		return name
	}
	if l.NameFormat == `~` {
		l.NameFormat = `Language[` + l.ID + `][%s]`
	}
	return fmt.Sprintf(l.NameFormat, name)
}

func (l *Language) HasName(name string) bool {
	if l.fieldMap == nil {
		return false
	}
	_, ok := l.fieldMap[name]
	return ok
}

func (l *Language) AddField(f ...FormElement) {
	if l.fieldMap == nil {
		l.fieldMap = map[string]int{}
		l.fields = []FormElement{}
	}
	for _, field := range f {
		name := l.Name(field.OriginalName())
		if _, ok := l.fieldMap[name]; ok {
			continue
		}
		l.fieldMap[name] = len(l.fields)
		l.fields = append(l.fields, field)
	}
}

func (l *Language) Field(name string) FormElement {
	if l.fieldMap == nil {
		return nil
	}
	if idx, ok := l.fieldMap[l.Name(name)]; ok {
		return l.fields[idx]
	}
	return nil
}

func (l *Language) Fields() []FormElement {
	return l.fields
}

func (l *Language) Clone() *Language {
	lg := NewLanguage(l.ID, l.Label, l.NameFormat)
	copy(lg.fields, l.fields)
	for k, v := range l.fieldMap {
		lg.fieldMap[k] = v
	}
	return lg
}
