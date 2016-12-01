package errors

import (
	"encoding/gob"
)

func init() {
	gob.Register(&Success{})
}

func NewOk(v string) Successor {
	return &Success{
		Value: v,
	}
}

type Success struct {
	Value string
}

func (s *Success) Success() string {
	return s.Value
}

func (s *Success) String() string {
	return s.Value
}

type Successor interface {
	Success() string
}

func IsOk(err interface{}) bool {
	_, y := err.(Successor)
	return y
}

func Ok(err interface{}) string {
	if v, y := err.(Successor); y {
		return v.Success()
	}
	return ``
}
