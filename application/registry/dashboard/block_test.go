package dashboard

import (
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestBlock(t *testing.T) {
	b := Blocks{}
	b.Set(-1, &Block{
		Tmpl: `test1`,
	})
	test.Eq(t, 1, b.Size())
	b.Add(&Block{
		Tmpl: `test2`,
	})
	test.Eq(t, 2, b.Size())
}
