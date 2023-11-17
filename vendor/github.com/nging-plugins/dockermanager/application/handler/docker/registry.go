package docker

import (
	"github.com/docker/docker/api/types/registry"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/handler"

	"github.com/nging-plugins/dockermanager/application/library/dockerclient"
	"github.com/nging-plugins/dockermanager/application/request"
)

func Login(ctx echo.Context) error {
	c, err := dockerclient.Client()
	if err != nil {
		return err
	}
	if ctx.IsPost() {
		req := echo.GetValidated(ctx).(*request.Login)
		var result registry.AuthenticateOKBody
		result, err = c.RegistryLogin(ctx, req.AuthConfig)
		if err != nil {
			goto END
		}

		ctx.Logger().Debugf(`RegistryLogin: %+v`, result)
		return ctx.Redirect(handler.URLFor(`/docker/base/info/index`))
	}

END:
	ctx.Set(`activeURL`, `/docker/base/registry/index`)
	ctx.Set(`title`, ctx.T(`仓库登录`))
	return ctx.Render(`docker/base/registry/login`, handler.Err(ctx, err))
}
