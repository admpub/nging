package fileperm

import (
	"encoding/json"
)

type Rules []*Rule

func (s Rules) Len() int { return len(s) }
func (s Rules) Less(i, j int) bool {
	return len(s[i].Path) < len(s[j].Path)
}
func (s Rules) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s *Rules) Add(r *Rule) (err error) {
	err = r.Init()
	*s = append(*s, r)
	return
}

func (s *Rules) Init() (err error) {
	for _, r := range *s {
		err = r.Init()
		if err != nil {
			return err
		}
	}
	return err
}

func (s Rules) IsEmpty() bool {
	return len(s) == 0
}

func (s Rules) JSONBytes() ([]byte, error) {
	return json.Marshal(s)
}

func (s Rules) JSONString() (string, error) {
	b, err := s.JSONBytes()
	return string(b), err
}
