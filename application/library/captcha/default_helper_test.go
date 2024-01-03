package captcha

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFixTemplatePath(t *testing.T) {
	r := fixTemplatePath(TypeDefault, `#default#default`)
	assert.Equal(t, `#default#captcha/default/default`, r)
}
