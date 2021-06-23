package perm

import (
	"fmt"

	"github.com/webx-top/echo"
)

type Behavior struct {
	Name         string
	ValueType    string // list / number / json
	VTypeOptions echo.H `json:",omitempty" xml:",omitempty"`
}

type BehaviorOption func(*Behavior)

func BehaviorOptName(name string) BehaviorOption {
	return func(a *Behavior) {
		a.Name = name
	}
}

func BehaviorOptValueType(vt string) BehaviorOption {
	return func(a *Behavior) {
		a.ValueType = vt
	}
}

func BehaviorOptVTypeOptions(opts echo.H) BehaviorOption {
	return func(a *Behavior) {
		a.VTypeOptions = opts
	}
}

func BehaviorOptVTypeOption(key string, value interface{}) BehaviorOption {
	return func(a *Behavior) {
		if a.VTypeOptions == nil {
			a.VTypeOptions = echo.H{}
		}
		a.VTypeOptions.Set(key, value)
	}
}

func NewBehavior(opts ...BehaviorOption) *Behavior {
	a := &Behavior{}
	for _, option := range opts {
		option(a)
	}
	return a
}

type Behaviors struct {
	*echo.KVData
}

func NewBehaviors() *Behaviors {
	return &Behaviors{
		KVData: echo.NewKVData(),
	}
}

func (m *Behaviors) Register(key string, value string, options ...interface{}) {
	aOpts := []BehaviorOption{}
	kOpts := []echo.KVOption{}
	for _, o := range options {
		switch v := o.(type) {
		case BehaviorOption:
			aOpts = append(aOpts, v)
		case echo.KVOption:
			kOpts = append(kOpts, v)
		default:
			panic(fmt.Sprintf(`unsupported type: %T`, v))
		}
	}
	a := NewBehavior(aOpts...)
	a.Name = key
	kOpts = append(kOpts, echo.KVOptX(a))
	m.Add(key, value, kOpts...)
}
