package container

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
)

func Kill(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	err = c.ContainerKill(ctx, containerID, `SIGKILL`)
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
}

func Start(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	opts := types.ContainerStartOptions{}
	err = c.ContainerStart(ctx, containerID, opts)
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
}

func Stop(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	err = c.ContainerStop(ctx, containerID, container.StopOptions{})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
}

func Restart(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	err = c.ContainerRestart(ctx, containerID, container.StopOptions{})
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
}

func Pause(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	containerID := ctx.Param(`id`)
	pause := ctx.Formx(`pause`)
	if len(pause.String()) == 0 || pause.Bool() {
		err = c.ContainerPause(ctx, containerID)
	} else {
		err = c.ContainerUnpause(ctx, containerID)
	}
	if err == nil {
		handler.SendOk(ctx, ctx.T(`操作成功`))
	} else {
		handler.SendFail(ctx, err.Error())
	}

	return ctx.Redirect(handler.URLFor(`/docker/base/container/index`))
}
