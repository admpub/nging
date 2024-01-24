package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/admpub/events"
	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"

	"github.com/admpub/nging/v5/application/library/config"
)

var signals = []os.Signal{
	os.Interrupt,    // ctrl+c 信号
	syscall.SIGTERM, // pkill 信号
}

var signalOperations = map[os.Signal][]func(int, engine.Engine, int){
	os.Interrupt:    {stopWebServer},
	syscall.SIGTERM: {stopWebServerForce},
}

func RegisterSignal(s os.Signal, op ...func(int, engine.Engine, int)) {
	for _, sig := range signals {
		if sig == s {
			goto ROP
		}
	}
	signals = append(signals, s)

ROP:
	if len(op) > 0 {
		if _, ok := signalOperations[s]; !ok {
			signalOperations[s] = []func(int, engine.Engine, int){}
		}
		signalOperations[s] = append(signalOperations[s], op...)
	}
}

func stopWebServerWithTimeout(eng engine.Engine, d time.Duration, exitCode int) {
	stopWebServer(0, eng, exitCode)
	time.Sleep(d)
	stopWebServer(1, eng, exitCode)
}

func stopWebServerForce(_ int, eng engine.Engine, exitCode int) {
	stopWebServerWithTimeout(eng, time.Second*5, exitCode)
}

const (
	ExitCodeStopFailed     = 2
	ExitCodeShutdownFailed = 4
	ExitCodeDefaultError   = 3
	ExitCodeSelfRestart    = 124
)

func stopWebServer(i int, eng engine.Engine, exitCode int) {
	if i > 0 {
		err := eng.Stop()
		if err != nil {
			log.Errorf(`failed to engine.Stop: %v`, err.Error())
		}
		if exitCode > 0 {
			os.Exit(exitCode)
		} else {
			os.Exit(ExitCodeStopFailed)
		}
	}
	log.Warn("SIGINT: Shutting down")
	go func() {
		config.FromCLI().Close()
		err := eng.Shutdown(context.Background())
		exitedCode := exitCode
		if err != nil {
			log.Errorf(`failed to engine.Shutdown: %v`, err.Error())
			exitedCode = ExitCodeShutdownFailed
		}
		os.Exit(exitedCode)
	}()
}

func CallSignalOperation(sig os.Signal, i int, eng engine.Engine) {
	var exitCode int
	if ec, ok := sig.(ExitCoder); ok {
		exitCode = ec.ExitCode()
	} else {
		exitCode = ExitCodeDefaultError
	}
	if operations, ok := signalOperations[sig]; ok {
		for _, operation := range operations {
			operation(i, eng, exitCode)
		}
	}
}

func SendSignal(sig os.Signal, exitCode int) {
	echo.FireByNameWithMap(`nging.signal`, events.Map{`signal`: sig, `exitCode`: exitCode})
}

func NewSignalWithExitCode(sig os.Signal, exitCode int) SignalWithExitCode {
	return SignalWithExitCode{signal: sig, exitCode: exitCode}
}

type SignalWithExitCode struct {
	signal   os.Signal
	exitCode int
}

func (s SignalWithExitCode) Signal() {
	s.signal.Signal()
}

func (s SignalWithExitCode) String() string {
	return s.signal.String()
}

func (s SignalWithExitCode) ExitCode() int {
	return s.exitCode
}

type ExitCoder interface {
	ExitCode() int
}

func handleSignal(eng engine.Engine) {
	signal.Reset(signals...)
	shutdown := make(chan os.Signal, 1)
	echo.OnCallback(`nging.signal`, func(e events.Event) error {
		sig, ok := e.Context.Get(`signal`).(os.Signal)
		if !ok {
			sig = os.Interrupt
		}
		exitCode, ok := e.Context.Get(`exitCode`).(int)
		if ok {
			shutdown <- &SignalWithExitCode{signal: sig, exitCode: exitCode}
			return nil
		}
		shutdown <- sig
		return nil
	}, `nging.signal`)
	signal.Notify(
		shutdown,
		signals...,
	)
	for i := 0; true; i++ {
		sig := <-shutdown
		log.Info(`received signal: ` + sig.String())
		fmt.Println(`received signal: ` + sig.String())
		CallSignalOperation(sig, i, eng)
	}
}
