package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func TestFixWd(t *testing.T) {
	fixWd()
	assert.Equal(t, ``, echo.Wd())
}
