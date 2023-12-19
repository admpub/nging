package main

import (
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/admpub/go-ps"
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

func isExitCode(exitCode int, exitCodes []int) bool {
	for _, code := range exitCodes {
		if exitCode == code {
			return true
		}
	}
	return false
}

func underMainProcess() bool {
	ppid := os.Getppid()
	if ppid == 1 {
		return false
	}
	proc, err := ps.FindProcess(ppid)
	if err != nil {
		log.Println(`ps.FindProcess: `, err)
		return false
	}
	name := filepath.Base(proc.Executable())
	matched := MAIN_EXE == name
	if matched {
		proc, err := os.FindProcess(ppid)
		if err != nil {
			log.Println(`os.FindProcess: `, err)
			return false
		}
		if err = proc.Kill(); err != nil {
			log.Println(`proc.Kill: `, err)
		}
	}
	return matched
}

func main() {
	var exitCodes []int
	for _, exitCode := range strings.Split(EXIT_CODE, `,`) {
		exitCodeN, _ := strconv.Atoi(exitCode)
		if exitCodeN > 0 {
			exitCodes = append(exitCodes, exitCodeN)
		}
	}
	if len(exitCodes) == 0 {
		exitCodes = append(exitCodes, 124)
	}
	executable, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	executable, err = filepath.EvalSymlinks(executable)
	if err != nil {
		log.Fatal(err)
	}
	var disabledLoop bool
	workDir := filepath.Dir(executable)
	executable = filepath.Join(workDir, MAIN_EXE)
	procArgs := []string{executable}
	if len(os.Args) > 1 {
		disabledLoop = os.Args[1] != `service` && !strings.HasPrefix(os.Args[1], `-`)
		procArgs = append(procArgs, os.Args[1:]...)
	}
	if !disabledLoop && len(procArgs) > 2 {
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
		if proc != nil {
			proc.Signal(sig)
		}
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
			if isExitCode(state.ExitCode(), exitCodes) {
				if underMainProcess() {
					log.Println(`[UnderMainProcess]exitCode:`, state.ExitCode())
					proc.Signal(syscall.SIGTERM)
					os.Exit(0)
				}
				goto START
			}
			log.Println(`exitCode:`, state.ExitCode())
			os.Exit(0)
		}
	}
}
