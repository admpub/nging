package render

import "strings"

type Config struct {
	Theme        string
	Engine       string
	Style        string
	Reload       bool
	ParseStrings map[string]string
}

func (t *Config) Parser() func([]byte) []byte {
	if t.ParseStrings == nil {
		return nil
	}
	return func(b []byte) []byte {
		s := string(b)
		for oldVal, newVal := range t.ParseStrings {
			s = strings.Replace(s, oldVal, newVal, -1)
		}
		return []byte(s)
	}
}
