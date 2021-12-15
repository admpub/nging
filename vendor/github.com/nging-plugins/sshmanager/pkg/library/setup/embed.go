package setup

import (
	_ "embed"

	"github.com/admpub/nging/v4/application/handler/setup"
)

//go:embed install.sql
var installSQL string

func init() {
	setup.RegisterInstallSQL(`nging`, installSQL)
}
