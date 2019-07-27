package mysql

import (
	"testing"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/testing/test"
)

func TestEnumOptions(t *testing.T) {
	s := `'ab','bc'`
	test.True(t, reContainsEnumLength.MatchString(s))
	matches := reEnumLength.FindAllStringSubmatch(s, -1)
	echo.Dump(matches)
	var r string
	for index, values := range matches {
		if index > 0 {
			r += `,`
		}
		r += values[0]
	}
	r = "(" + r + ")"
	test.Eq(t, "('ab','bc')", r)
}
