package cmd

import (
	"testing"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
)

func TestFixWd(t *testing.T) {
	config.FixWd()
	assert.Equal(t, ``, echo.Wd())
}
