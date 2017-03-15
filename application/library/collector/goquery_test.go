package collector

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestParseSelections(t *testing.T) {
	s := parseSelections(`$('#eee').find('.aaa').text`)
	com.Dump(s)
	for i, v := range s {
		switch i {
		case 0:
			assert.Equal(t, "$", v.Function)
			assert.Equal(t, "[#eee]", fmt.Sprint(v.Parameters))
		case 1:
			assert.Equal(t, "find", v.Function)
			assert.Equal(t, "[.aaa]", fmt.Sprint(v.Parameters))
		case 2:
			assert.Equal(t, "text", v.Function)
			assert.Equal(t, "[]", fmt.Sprint(v.Parameters))
		}
	}
	s = parseSelections(`$('#eee').find( '.aa\\\'a' ).eq( 1 ).match('/<a>(.*)</a>/i','$1')`)
	com.Dump(s)
	for i, v := range s {
		switch i {
		case 0:
			assert.Equal(t, "$", v.Function)
			assert.Equal(t, "[#eee]", fmt.Sprint(v.Parameters))
		case 1:
			assert.Equal(t, "find", v.Function)
			assert.Equal(t, "[.aa\\'a]", fmt.Sprint(v.Parameters))
		case 2:
			assert.Equal(t, "eq", v.Function)
			assert.Equal(t, "[1]", fmt.Sprint(v.Parameters))
		case 3:
			assert.Equal(t, "match", v.Function)
			assert.Equal(t, "[/<a>(.*)</a>/i $1]", fmt.Sprint(v.Parameters))
		}
	}
	s = parseSelections(`$('.aaa').each('$(".a").find(\'ddd\')')`)
	com.Dump(s)
	for i, v := range s {
		switch i {
		case 0:
			assert.Equal(t, "$", v.Function)
			assert.Equal(t, "[.aaa]", fmt.Sprint(v.Parameters))
		case 1:
			assert.Equal(t, "each", v.Function)
			assert.Equal(t, "[$(\".a\").find('ddd')]", fmt.Sprint(v.Parameters))
		}
	}
}
