//go:build ignore
// +build ignore

package main

import (
	"github.com/admpub/log"
	_ "github.com/admpub/nging/v4/application"
	"github.com/admpub/nging/v4/application/cmd"
	_ "github.com/admpub/nging/v4/application/initialize/manager"
	_ "github.com/admpub/nging/v4/application/library/sqlite"

	_ "github.com/admpub/nging/v4/application/cmd/testsuite"
)

// usage: go run testsuite.go testsuite --name=sqlquery

func main() {
	defer log.Sync()
	exec()
}

func exec() {
	cmd.Execute()
}
