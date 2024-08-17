package captcha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixTemplatePath(t *testing.T) {
	r, s := fixTemplatePath(TypeDefault, `#default#default`)
	assert.Equal(t, `#default#captcha/default/default`, r)
	assert.Equal(t, `default`, s)
}
