package main

import (
	_ "github.com/admpub/nging/v5/tool/startup/ico"
	"github.com/admpub/nging/v5/tool/startup/pkg"
)

var (
	BUILD_TIME string
	BUILD_OS   string
	BUILD_ARCH string
	CLOUD_GOX  string
	COMMIT     string
	VERSION    = `0.0.1`
	MAIN_EXE   = `nging`
	EXIT_CODE  = `124`
)

func main() {
	pkg.Start(MAIN_EXE, EXIT_CODE)
}
