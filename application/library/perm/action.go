package perm

import (
	"fmt"

	"github.com/webx-top/echo"
)

type Action struct {
	Name         string
	ValueType    string // list / number / json
	VTypeOptions echo.H `json:",omitempty" xml:",omitempty"`
}

type ActionOption func(*Action)

func ActionOptName(name string) ActionOption {
	return func(a *Action) {
		a.Name = name
	}
}

func ActionOptValueType(vt string) ActionOption {
	return func(a *Action) {
		a.ValueType = vt
	}
}

func ActionOptVTypeOptions(opts echo.H) ActionOption {
	return func(a *Action) {
		a.VTypeOptions = opts
	}
}

func ActionOptVTypeOption(key string, value interface{}) ActionOption {
	return func(a *Action) {
		if a.VTypeOptions == nil {
			a.VTypeOptions = echo.H{}
		}
		a.VTypeOptions.Set(key, value)
	}
}

func NewAction(opts ...ActionOption) *Action {
	a := &Action{}
	for _, option := range opts {
		option(a)
	}
	return a
}

type Actions struct {
	*echo.KVData
}

func NewActions() *Actions {
	return &Actions{
		KVData: echo.NewKVData(),
	}
}

func (m *Actions) Register(key string, value string, options ...interface{}) {
	aOpts := []ActionOption{}
	kOpts := []echo.KVOption{}
	for _, o := range options {
		switch v := o.(type) {
		case ActionOption:
			aOpts = append(aOpts, v)
		case echo.KVOption:
			kOpts = append(kOpts, v)
		default:
			panic(fmt.Sprintf(`unsupported type: %T`, v))
		}
	}
	a := NewAction(aOpts...)
	a.Name = key
	kOpts = append(kOpts, echo.KVOptX(a))
	m.Add(key, value, kOpts...)
}
