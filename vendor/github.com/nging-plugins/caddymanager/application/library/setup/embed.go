package setup

import (
	_ "embed"

	"github.com/admpub/nging/v5/application/library/config"
)

//go:embed install.sql
var installSQL string

func RegisterSQL(sc *config.SQLCollection) {
	sc.RegisterInstall(`nging`, installSQL)
}
