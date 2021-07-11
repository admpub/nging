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
