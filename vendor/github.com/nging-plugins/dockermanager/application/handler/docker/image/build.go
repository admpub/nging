package image

import (
	"github.com/docker/docker/api/types"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	"github.com/admpub/nging/v5/application/handler"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/request"
)

func Build(ctx echo.Context) error {
	user := handler.User(ctx)
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.ImageBuild)
		err = dockerclient.BuildImage(ctx, user, nil, c, types.ImageBuildOptions{
			Tags:        req.Tags,
			PullParent:  req.PullParent,
			Dockerfile:  req.Dockerfile,
			BuildArgs:   req.BuildArgs,
			AuthConfigs: req.AuthConfigs,
			Target:      req.Target,
			Version:     req.Version,
			Platform:    req.Platform,
		})
		if err != nil {
			return err
		}
		data := ctx.Data()
		data.SetInfo(ctx.T(`启动成功`), code.Success.Int())
		return ctx.JSON(data)
	}
	ctx.Set(`activeURL`, `/docker/base/image/index`)
	ctx.Set(`title`, ctx.T(`构建镜像`))
	return ctx.Render(`docker/base/image/build`, nil)
}
