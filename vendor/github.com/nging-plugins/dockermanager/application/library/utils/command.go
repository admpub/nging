package utils

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"sync"

	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/webx-top/com"
)

func RunCommand(ctx context.Context, command string, args []string, noticer notice.Noticer, opts ...func(*exec.Cmd)) (outStr string, errStr string, err error) {
	cmd := exec.CommandContext(ctx, command, args...)
	for _, opt := range opts {
		opt(cmd)
	}
	var errReader io.ReadCloser
	var outReader io.ReadCloser
	errReader, err = cmd.StderrPipe()
	if err != nil {
		return
	}
	defer errReader.Close()
	outReader, err = cmd.StdoutPipe()
	if err != nil {
		return
	}
	defer outReader.Close()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		com.SeekLines(outReader, func(line string) error {
			if len(line) == 0 {
				return nil
			}
			outStr += line + "\n"
			return noticer(line, notice.StateSuccess)
		})
	}()
	go func() {
		defer wg.Done()
		com.SeekLines(errReader, func(line string) error {
			if len(line) == 0 {
				return nil
			}
			errStr += line + "\n"
			//log.Debug(line)
			return noticer(line, notice.StateSuccess)
		})
	}()
	err = cmd.Run()
	wg.Wait()
	if err != nil {
		err = fmt.Errorf(`%w: %s`, err, errStr)
		return
	}
	return
}

func ExtName() string {
	var ext string
	if com.IsWindows {
		ext = `.exe`
	}
	return ext
}

func DockerPath() string {
	return com.Getenv(`DOCKER_PATH`, `docker`+ExtName())
}

var composeCmd string
var composeSub string
var composeOnce sync.Once

func initComposeCmd() {
	err := exec.Command(DockerPath(), `compose`, `ls`).Run()
	if err == nil {
		composeCmd = DockerPath()
		composeSub = `compose`
		return
	}
	composeCmd = com.Getenv(`DOCKER_COMPOSE_PATH`, `docker-compose`+ExtName())
}

func DockerCompose(args []string) (string, []string) {
	composeOnce.Do(initComposeCmd)
	if len(composeSub) > 0 {
		return composeCmd, append([]string{composeSub}, args...)
	}
	return composeCmd, args
}
