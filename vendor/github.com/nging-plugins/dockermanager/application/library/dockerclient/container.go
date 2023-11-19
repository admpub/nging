package dockerclient

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/nging-plugins/dockermanager/application/library/utils"
)

func Exec(ctx context.Context, containerID string, cmd []string, env []string, outWriter io.Writer, errWriter io.Writer) error {
	c, err := Client()
	if err != nil {
		return fmt.Errorf(`client error: %w`, err)
	}
	cfg := types.ExecConfig{
		Tty:          true,
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
		Env:          env,
	}
	idResp, err := c.ContainerExecCreate(ctx, containerID, cfg)
	if err != nil {
		return err
	}
	var response types.HijackedResponse
	response, err = c.ContainerExecAttach(ctx, idResp.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer response.Close()
	buf := bytes.NewBuffer(nil)
	if errWriter == nil {
		errWriter = os.Stderr
	}
	if outWriter == nil {
		outWriter = os.Stdout
	}
	_, err = utils.StdCopy(outWriter, io.MultiWriter(errWriter, buf), response.Conn)
	if err != nil {
		return err
	}
	errMsg := buf.String()
	if len(errMsg) > 0 {
		err = errors.New(errMsg)
	}
	return err
}
