package captcha

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuote(t *testing.T) {
	r := strconv.Quote(`a"b\`)
	assert.Equal(t, `"a\"b\\"`, r)
}
