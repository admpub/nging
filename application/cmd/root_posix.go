//go:build !windows && !plan9 && !nacl && !js
// +build !windows,!plan9,!nacl,!js

package cmd

import (
	"syscall"

	"github.com/admpub/nging/v5/application/library/config"
)

func init() {
	RegisterSignal(syscall.SIGHUP, func() {
		config.FromCLI().SendSignalToAllCmd(syscall.SIGQUIT)
	})
}
