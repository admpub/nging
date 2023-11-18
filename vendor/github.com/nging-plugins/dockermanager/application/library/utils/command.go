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
			return noticer(line, notice.StateFailure)
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

func DockerPath() string {
	return com.Getenv(`DOCKER_PATH`, `docker`)
}
