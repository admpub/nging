//go:build !windows && !plan9 && !nacl && !js
// +build !windows,!plan9,!nacl,!js

package cmd

import (
	"syscall"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/echo/engine"
)

func init() {
	RegisterSignal(syscall.SIGHUP /*终端关闭*/, func(i int, eng engine.Engine, exitCode int) {
		config.FromCLI().SendSignalToAllCmd(syscall.SIGQUIT)
		CallSignalOperation(SignalWithExitCode{signal: syscall.SIGTERM, exitCode: exitCode}, i, eng)
	})
}
