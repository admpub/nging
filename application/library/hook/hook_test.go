package hook

import (
	"bytes"
	"testing"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

func TestHooks(t *testing.T) {
	h := New()
	b := bytes.NewBufferString(``)

	h.On(`test`, func(_ echo.H) error {
		b.WriteString(`ok`)
		return nil
	})
	test.Eq(t, 1, h.Size(`test`))

	h.On(`test`, func(_ echo.H) error {
		b.WriteString(`/`)
		return nil
	})
	test.Eq(t, 2, h.Size(`test`))

	h.On(`test`, func(_ echo.H) error {
		b.WriteString(`no`)
		return nil
	})
	test.Eq(t, 3, h.Size(`test`))

	err := h.Fire(`test`, nil)
	if err != nil {
		t.Fatal(err)
	}
	test.Eq(t, `ok/no`, b.String())

	h.Off(`test`)
	test.Eq(t, 0, h.Size(`test`))
	test.Eq(t, []string{}, h.Names())
}
