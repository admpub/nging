package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
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
	var proc *os.Process
	var n int
	var pid int
	log.Println(strings.Join(procArgs, ` `))

START:
	proc, err = os.StartProcess(executable, procArgs, &os.ProcAttr{
		Dir:   workDir,
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	})
	if err != nil {
		log.Fatal(err)
	}
	n++
	pid = proc.Pid
	//proc.Wait()

	for {
		time.Sleep(time.Second * 2)
		_, err := os.FindProcess(pid)
		if err != nil {
			if n > 5 {
				log.Fatal(err)
			}
			goto START
		}
		n = 0
	}
}
