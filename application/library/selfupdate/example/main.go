package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"time"
)

func main() {
	out, err := exec.Command(`go`, `build`, `-o`, `instance.test`, `./instance`).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(out))
	wd, err := filepath.Abs(`.`)
	if err != nil {
		log.Fatal(err)
	}
	executable := `./instance.test`
	os.Chmod(executable, 0755)
	log.Println(wd)
	_, err = os.StartProcess(executable, []string{}, &os.ProcAttr{
		Dir:   wd,
		Env:   os.Environ(),
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Sys:   &syscall.SysProcAttr{},
	})
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 5)
	os.Exit(0)
}
