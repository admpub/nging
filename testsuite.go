// +build ignore

package main

import (
	"github.com/admpub/log"
	_ "github.com/admpub/nging/application"
	"github.com/admpub/nging/application/cmd"
	_ "github.com/admpub/nging/application/initialize/manager"
	_ "github.com/admpub/nging/application/library/sqlite"

	_ "github.com/admpub/nging/application/cmd/testsuite"
)

// usage: go run testsuite.go testsuite --name=sqlquery

func main() {
	defer log.Sync()
	exec()
}

func exec() {
	cmd.Execute()
}
