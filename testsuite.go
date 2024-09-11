//go:build ignore
// +build ignore

package main

import (
	"github.com/admpub/log"
	_ "github.com/admpub/nging/v5/application"
	"github.com/coscms/webcore/cmd"
	_ "github.com/coscms/webcore/initialize/manager"
	_ "github.com/coscms/webcore/library/sqlite"

	_ "github.com/coscms/webcore/cmd/testsuite"
)

// usage: go run testsuite.go testsuite --name=sqlquery

func main() {
	defer log.Sync()
	exec()
}

func exec() {
	cmd.Execute()
}
