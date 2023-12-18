package main

import (
	"log"
	"os"
	"syscall"
	"time"
)

func main() {
	/*/
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
	//*/
	wd := `/Users/hank/go/src/github.com/admpub/nging/dist/localtest/nging_v5.2.0/nging_darwin_amd64`
	executable := wd + `/nging`
	_, err := os.StartProcess(executable, []string{executable, `-p`, `19999`}, &os.ProcAttr{
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
