package setup

import (
	_ "embed"
)

//go:embed install.sql
var InstallSQL string
