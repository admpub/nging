package pkg

import (
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/admpub/go-ps"
	"github.com/admpub/log"
)

var (
	logger = log.New(`startup`)
)

func isExitCode(exitCode int, exitCodes []int) bool {
	for _, code := range exitCodes {
		if exitCode == code {
			return true
		}
	}
	return false
}

func underMainProcess(mainExe string) bool {
	ppid := os.Getppid()
	if ppid == 1 {
		return false
	}
	proc, err := ps.FindProcess(ppid)
	if err != nil {
		logger.Debug(`ps.FindProcess: `, err)
		return false
	}
	if proc == nil {
		return false
	}
	name := filepath.Base(proc.Executable())
	matched := mainExe == name
	if matched {
		proc, err := os.FindProcess(ppid)
		if err != nil {
			logger.Debug(`os.FindProcess: `, err)
			return false
		}
		if err = proc.Kill(); err != nil && err != os.ErrProcessDone {
			logger.Debug(`proc.Kill: `, err)
		}
	}
	return matched
}

func Start(mainExe, exitCode string) {
	logger.Sync()
	var exitCodes []int
	for _, exitCode := range strings.Split(exitCode, `,`) {
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

	logDir := filepath.Join(workDir, `data/logs`)
	os.MkdirAll(logDir, os.ModePerm)
	ft := log.NewFileTarget()
	filepathSeparator := string([]byte{filepath.Separator})
	ft.FileName = logDir + filepathSeparator + `startup_{date:20060102}.log` //按天分割日志
	ft.MaxLevel = log.LevelInfo
	logger.AddTarget(ft)

	executable = filepath.Join(workDir, mainExe)
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
	logger.Debug(strings.Join(procArgs, ` `))
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
		logger.Debug(`received signal: `, sig.String())
		os.Exit(0)
	}()

	pidDir := filepath.Join(workDir, `data/pid`)
	os.MkdirAll(pidDir, os.ModePerm)

START:
	proc, err = os.StartProcess(executable, procArgs, &os.ProcAttr{
		Dir:   workDir,
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		logger.Fatal(err)
	}

	filenameRE := regexp.MustCompile(`[^\w]+`)
	pidFilename := filenameRE.ReplaceAllString(strings.SplitN(mainExe, `.`, 2)[0], `_`)
	os.WriteFile(pidDir+filepathSeparator+pidFilename+`_forked.pid`, []byte(strconv.Itoa(proc.Pid)), os.ModePerm)

	state, err = proc.Wait()
	if disabledLoop {
		return
	}
	if err != nil {
		logger.Error(err)
		goto START
	}
	if state.Exited() {
		if isExitCode(state.ExitCode(), exitCodes) && !underMainProcess(mainExe) {
			time.Sleep(500 * time.Millisecond)
			goto START
		}
		logger.Info(`exit code: `, state.ExitCode())
		os.Exit(state.ExitCode())
	}
}
