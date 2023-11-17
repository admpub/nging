package container

import (
	"io"
	"time"

	"github.com/webx-top/echo"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/library/utils"
)

// Export 导出容器快照
func Export(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	var reader io.ReadCloser
	reader, err = c.ContainerExport(ctx, containerID)
	if err != nil {
		return err
	}
	defer reader.Close()
	fileName := `docker-container-` + utils.ShortenID(containerID) + `.tar`
	echo.SetAttachmentHeader(ctx, fileName, false)
	return ctx.ServeContent(reader, fileName, time.Now())
}
