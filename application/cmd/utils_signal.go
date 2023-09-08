package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/webx-top/echo/engine"
)

var signals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
}

var signalOperations = map[os.Signal][]func(int, engine.Engine){
	os.Interrupt:    {stopWebServer},
	syscall.SIGTERM: {stopWebServerForce},
}

func RegisterSignal(s os.Signal, op ...func(int, engine.Engine)) {
	for _, sig := range signals {
		if sig == s {
			goto ROP
		}
	}
	signals = append(signals, s)

ROP:
	if len(op) > 0 {
		if _, ok := signalOperations[s]; !ok {
			signalOperations[s] = []func(int, engine.Engine){}
		}
		signalOperations[s] = append(signalOperations[s], op...)
	}
}

func stopWebServerWithTimeout(eng engine.Engine, d time.Duration) {
	stopWebServer(0, eng)
	time.Sleep(d)
	stopWebServer(1, eng)
}

func stopWebServerForce(i int, eng engine.Engine) {
	stopWebServerWithTimeout(eng, time.Second*5)
}

func stopWebServer(i int, eng engine.Engine) {
	if i > 0 {
		err := eng.Stop()
		if err != nil {
			log.Error(err.Error())
		}
		os.Exit(2)
	}
	log.Warn("SIGINT: Shutting down")
	go func() {
		config.FromCLI().Close()
		err := eng.Shutdown(context.Background())
		var exitedCode int
		if err != nil {
			log.Error(err.Error())
			exitedCode = 4
		}
		os.Exit(exitedCode)
	}()
}

func CallSignalOperation(sig os.Signal, i int, eng engine.Engine) {
	if operations, ok := signalOperations[sig]; ok {
		for _, operation := range operations {
			operation(i, eng)
		}
	}
}

func handleSignal(eng engine.Engine) {
	shutdown := make(chan os.Signal, 1)
	// ctrl+c信号os.Interrupt
	// pkill信号syscall.SIGTERM
	signal.Notify(
		shutdown,
		signals...,
	)
	for i := 0; true; i++ {
		sig := <-shutdown
		CallSignalOperation(sig, i, eng)
	}
}
