package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	BUILD_TIME string
	BUILD_OS   string
	BUILD_ARCH string
	CLOUD_GOX  string
	COMMIT     string
	VERSION    = `0.0.1`
	MAIN_EXE   = `nging`
)

func main() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		log.Fatal(err)
	}
	workDir := filepath.Dir(executable)
	executable = filepath.Join(workDir, MAIN_EXE)
	procArgs := []string{executable}
	if len(os.Args) > 1 {
		procArgs = append(procArgs, os.Args[1:]...)
	}
	var disabledLoop bool
	if len(procArgs) > 2 {
		disabledLoop = procArgs[1] == `service` && procArgs[2] != `run`
	}
	var proc *os.Process
	var state *os.ProcessState
	log.Println(strings.Join(procArgs, ` `))
	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(
			shutdown,
			os.Interrupt,    // ctrl+c 信号
			syscall.SIGTERM, // pkill 信号
		)
		sig := <-shutdown
		proc.Signal(sig)
		os.Exit(0)
	}()

START:
	proc, err = os.StartProcess(executable, procArgs, &os.ProcAttr{
		Dir:   workDir,
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		log.Fatal(err)
	}
	for {
		state, err = proc.Wait()
		if disabledLoop {
			return
		}
		if err != nil {
			goto START
		}
		if state.Exited() {
			goto START
		}
	}
}
