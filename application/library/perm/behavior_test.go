package perm

import (
	"testing"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

type testStruct struct {
	Test string `json:"test"`
}

func TestBehavior(t *testing.T) {
	be := &Behavior{
		Value: echo.H{},
	}
	test.Eq(t, `{}`, be.String())
	be = &Behavior{
		Value: &testStruct{
			Test: `OK`,
		},
	}
	test.Eq(t, `{"test":"OK"}`, be.String())
}

func TestParseBehavior(t *testing.T) {
	behaviors := NewBehaviors()
	behaviors.Register(`test1`, `测试1`, BehaviorOptValue(&testStruct{}), BehaviorOptValueInitor(func() interface{} {
		return &testStruct{}
	}))
	behaviors.Register(`test2`, `测试2`, BehaviorOptValue(1))
	behaviors.Register(`test3`, `测试3`, BehaviorOptValue(``))

	perms, err := ParseBehavior(`{"test1":{"test":"OK"},"test2":100,"test3":"333"}`, behaviors)
	if err != nil {
		panic(err)
	}

	result := perms.Get(`test1`).Value.(*testStruct)
	test.Eq(t, `OK`, result.Test)

	item := behaviors.GetItem(`test1`)
	typeBehavior := item.X.(*Behavior)
	test.NotEq(t, typeBehavior.Value.(*testStruct).Test, result.Test)

	result2 := perms.Get(`test2`).Value.(*int)
	test.Eq(t, 100, *result2)

	result3 := perms.Get(`test3`).Value.(*string)
	test.Eq(t, `333`, *result3)

	jsonData, err := SerializeBehaviorValues(map[string][]string{
		"test1": {`{"test":"OK"}`},
		"test2": {`100`},
		"test3": {`333`},
	}, behaviors)
	if err != nil {
		panic(err)
	}
	expected := `{"test1":{"test":"OK"},"test2":100,"test3":"333"}`
	test.Eq(t, expected, jsonData)
}

func TestParseBehavior2(t *testing.T) {
	behaviors := NewBehaviors()
	behaviors.Register(`test1`, `测试1`, BehaviorOptValue(&testStruct{}))
	perms, err := ParseBehavior(`{"test1":{"test":"OK"}}`, behaviors)
	if err != nil {
		panic(err)
	}

	result := perms.Get(`test1`).Value.(*testStruct)
	test.Eq(t, `OK`, result.Test)

	jsonData, err := SerializeBehaviorValues(map[string][]string{
		"test1": {`{"test":"OK"}`},
	}, behaviors)
	if err != nil {
		panic(err)
	}
	expected := `{"test1":{"test":"OK"}}`
	test.Eq(t, expected, jsonData)
}
