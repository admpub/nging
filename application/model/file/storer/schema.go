package storer

import (
	"github.com/webx-top/echo"
)

const (
	StorerInfoKey = `NgingStorer`
)

func NewInfo() *Info {
	return &Info{}
}

type Info struct {
	Name string `json:"name" xml:"name"`
	ID string `json:"id" xml:"id"`
}

func (s *Info) FromStore(v echo.H) *Info {
	s.Name = v.String("name")
	s.ID = v.String("id")
	if s.ID == `0` {
		s.ID = ``
	}
	return s
}
