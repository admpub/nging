package thumb

import (
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestRegistry(t *testing.T) {
	r := Registries{}
	r.Set(`test`, Size{Width: 100, Height: 50})
	expected := Sizes{Size{Width: 100, Height: 50}}
	test.Eq(t, expected, r.Get(`test`))

	r.Add(`test`, Size{Width: 101, Height: 51})
	expected = append(expected, Size{Width: 101, Height: 51})
	test.Eq(t, expected, r.Get(`test`))

	r.Add(`test`, Size{Width: 102, Height: 51})
	expected = append(expected, Size{Width: 103, Height: 51})
	test.NotEq(t, expected, r.Get(`test`))
}
