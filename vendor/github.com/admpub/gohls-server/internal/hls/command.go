package hls

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"syscall"

	log "github.com/sirupsen/logrus"
)

func execute(cmdPath string, args []string) (data []byte, err error) {
	cmd := exec.Command(cmdPath, args...)
	stdout, err1 := cmd.StdoutPipe()
	if err1 != nil {
		err = fmt.Errorf("Error opening stdout of command: %v", err1)
		return
	}
	defer stdout.Close()
	stderr, err1 := cmd.StderrPipe()
	if err1 != nil {
		err = fmt.Errorf("Error opening stderr of command: %v", err1)
		return
	}
	defer stderr.Close()

	log.Debugf("Executing: %v %v", cmdPath, args)
	err2 := cmd.Start()
	if err2 != nil {
		err = fmt.Errorf("Error starting command: %v", err2)
		return
	}

	var buffer bytes.Buffer
	_, err3 := io.Copy(&buffer, io.MultiReader(stdout, stderr))
	if err3 != nil {
		// Ask the process to exit
		cmd.Process.Signal(syscall.SIGKILL)
		cmd.Process.Wait()
		err = fmt.Errorf("Error copying stdout to buffer: %v", err3)
		return
	}
	err4 := cmd.Wait()
	if err4 != nil {
		err = fmt.Errorf("Command failed: %v", err4)
		return
	}
	data = buffer.Bytes()
	return
}
