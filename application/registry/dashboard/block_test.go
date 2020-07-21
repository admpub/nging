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
	b.Add(-1, &Block{
		Tmpl: `test2`,
	})
	test.Eq(t, 2, b.Size())
	b.Add(-1, &Block{
		Tmpl: `test3`,
	})
	test.Eq(t, 3, b.Size())
	b.Add(0, &Block{
		Tmpl: `test0`,
	})
	test.Eq(t, 4, b.Size())
	test.Eq(t, `test0`, b[0].Tmpl)
	test.Eq(t, `test1`, b[1].Tmpl)
	test.Eq(t, `test2`, b[2].Tmpl)
	test.Eq(t, `test3`, b[3].Tmpl)
}
