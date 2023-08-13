package cmd

import (
	"testing"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/echo"
)

func TestFixWd(t *testing.T) {
	config.FixWd()
	t.Logf(`Wd: %s`, echo.Wd())
}
