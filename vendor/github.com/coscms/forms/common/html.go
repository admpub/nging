package common

import (
	"html/template"
	"strings"

	"github.com/webx-top/com"
)

type (
	HTMLAttrValues []string
	HTMLAttributes map[template.HTMLAttr]interface{}
	HTMLData       map[string]interface{}
)

func (s HTMLAttrValues) String() string {
	return strings.Join([]string(s), ` `)
}

func (s HTMLAttrValues) IsEmpty() bool {
	return len(s) == 0
}

func (s HTMLAttrValues) Size() int {
	return len(s)
}

func (s HTMLAttrValues) Exists(attr string) bool {
	return com.InSlice(attr, s)
}

func (s *HTMLAttrValues) Add(value string) {
	(*s) = append((*s), value)
}

func (s *HTMLAttrValues) Remove(value string) {
	ind := -1
	for i, v := range *s {
		if v == value {
			ind = i
			break
		}
	}

	if ind != -1 {
		*s = append((*s)[:ind], (*s)[ind+1:]...)
	}
}

func (s HTMLAttributes) Exists(attr string) bool {
	_, ok := s[template.HTMLAttr(attr)]
	return ok
}

func (s HTMLAttributes) FillFrom(data map[string]interface{}) {
	for k, v := range data {
		s[template.HTMLAttr(k)] = v
	}
}

func (s HTMLAttributes) FillFromStringMap(data map[string]string) {
	for k, v := range data {
		s[template.HTMLAttr(k)] = v
	}
}

func (s HTMLData) Exists(key string) bool {
	_, ok := s[key]
	return ok
}
