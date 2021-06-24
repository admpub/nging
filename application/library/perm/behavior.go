package perm

import (
	"fmt"

	"github.com/webx-top/echo"
)

type Behavior struct {
	Name            string      `json:",omitempty" xml:",omitempty"`
	ValueType       string      `json:",omitempty" xml:",omitempty"` // list / number / json
	VTypeOptions    echo.H      `json:",omitempty" xml:",omitempty"`
	Value           interface{} `json:",omitempty" xml:",omitempty"` // 在Behaviors中登记时，代表默认值；在BehaviorPerms中登记时代表针对某个用户设置的值
	valueInitor     func() interface{}
	formValueParser func([]string) (interface{}, error)
}

func (b Behavior) IsValid() bool {
	return len(b.Name) > 0
}

func (b *Behavior) SetValueInitor(initor func() interface{}) {
	b.valueInitor = initor
}

func (b *Behavior) SetFormValueParser(parser func([]string) (interface{}, error)) {
	b.formValueParser = parser
}

type BehaviorOption func(*Behavior)

func BehaviorOptName(name string) BehaviorOption {
	return func(a *Behavior) {
		a.Name = name
	}
}

func BehaviorOptValueInitor(initor func() interface{}) BehaviorOption {
	return func(a *Behavior) {
		a.valueInitor = initor
	}
}

func BehaviorOptFormValueParser(parser func([]string) (interface{}, error)) BehaviorOption {
	return func(a *Behavior) {
		a.formValueParser = parser
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
