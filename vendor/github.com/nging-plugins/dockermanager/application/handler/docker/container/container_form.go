package container

import (
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/library/common"
)

func setContainerDataToFormAdd(ctx echo.Context, c *client.Client, data types.ContainerJSON) {
	imageInspect, _, _ := c.ImageInspectWithRaw(ctx, data.Image)
	if imageInspect.Config != nil {
		if len(data.Config.WorkingDir) > 0 && data.Config.WorkingDir == imageInspect.Config.WorkingDir {
			data.Config.WorkingDir = ``
		}
		if len(data.Config.Cmd) > 0 {
			data.Config.Cmd = com.StringSliceDiff(data.Config.Cmd, imageInspect.Config.Cmd)
		}
		if len(data.Config.Entrypoint) > 0 {
			data.Config.Entrypoint = com.StringSliceDiff(data.Config.Entrypoint, imageInspect.Config.Entrypoint)
		}
		if len(data.Config.Env) > 0 {
			data.Config.Env = com.StringSliceDiff(data.Config.Env, imageInspect.Config.Env)
		}
		if len(data.Config.Labels) > 0 && len(imageInspect.Config.Labels) > 0 {
			for k, v := range data.Config.Labels {
				if imageInspect.Config.Labels[k] == v {
					delete(data.Config.Labels, k)
				}
			}
		}
	}
	ctx.Request().Form().Set(`image`, data.Config.Image)
	if len(data.Config.WorkingDir) > 0 || len(data.Config.Cmd) > 0 || len(data.Config.Entrypoint) > 0 {
		ctx.Request().Form().Set(`commandEnabled`, `Y`)
		ctx.Request().Form().Set(`entrypoint`, strings.Join(data.Config.Entrypoint, ` `))
		ctx.Request().Form().Set(`command`, strings.Join(data.Config.Cmd, ` `))
		ctx.Request().Form().Set(`workingDir`, data.Config.WorkingDir)
	}
	ctx.Request().Form().Set(`env`, strings.Join(data.Config.Env, com.StrLF))
	labels := make([]string, 0, len(data.Config.Labels))
	for k, v := range data.Config.Labels {
		if strings.HasPrefix(k, `com.docker.compose.`) {
			continue
		}
		labels = append(labels, k+`=`+v)
	}
	ctx.Request().Form().Set(`labels`, strings.Join(labels, com.StrLF))
	if len(data.HostConfig.PortBindings) > 0 {
		ctx.Request().Form().Set(`portExport`, `Y`)
		for port, hosts := range data.HostConfig.PortBindings {
			for _, host := range hosts {
				to := host.HostPort
				if len(host.HostIP) > 0 {
					to = host.HostIP + `:` + to
				}
				ctx.Request().Form().Add(`containerPortTo`, to)
				ctx.Request().Form().Add(`containerPortNet`, port.Proto())
				ctx.Request().Form().Add(`containerPortFrom`, port.Port())
			}
		}
	}
	if len(data.HostConfig.Binds) > 0 {
		ctx.Request().Form().Set(`storageVolumeMount`, `Y`)
		for _, bd := range data.HostConfig.Binds {
			//bd := hostPath + `:` + containerPath + `:` + op
			var hostPath, op, containerPath string
			com.SliceExtract(strings.Split(bd, `:`), &hostPath, &containerPath, &op)
			ctx.Request().Form().Add(`containerPathOp`, op)
			ctx.Request().Form().Add(`containerPathTo`, hostPath)
			ctx.Request().Form().Add(`containerPathFrom`, containerPath)
		}
	}
	for _, v := range data.HostConfig.VolumesFrom {
		ctx.Request().Form().Add(`volumesFrom`, v)
	}
	ctx.Request().Form().Set(`tty`, common.BoolToFlag(data.Config.Tty))
	ctx.Request().Form().Set(`networkDisabled`, common.BoolToFlag(data.Config.NetworkDisabled))
	ctx.Request().Form().Set(`networkMode`, data.HostConfig.NetworkMode.NetworkName())
	if data.HostConfig.Privileged {
		ctx.Request().Form().Set(`privileged`, `true`)
	}
	if data.HostConfig.AutoRemove {
		ctx.Request().Form().Set(`autoRemove`, `true`)
	}
	capabilities := make([]string, len(data.HostConfig.CapAdd), len(data.HostConfig.CapAdd)+len(data.HostConfig.CapDrop))
	copy(capabilities, data.HostConfig.CapAdd)
	for _, ca := range data.HostConfig.CapDrop {
		capabilities = append(capabilities, `-`+ca)
	}
	ctx.Request().Form().Set(`capabilities`, strings.Join(capabilities, `,`))
	setContainerDataToFormEdit(ctx, data)
}

func setContainerDataToFormEdit(ctx echo.Context, data types.ContainerJSON) {
	ctx.Request().Form().Set(`name`, data.Name)
	ctx.Request().Form().Set(`memory`, param.AsString(data.HostConfig.Memory))
	ctx.Request().Form().Set(`cpuWeight`, param.AsString(data.HostConfig.CPUShares))
	ctx.Request().Form().Set(`restartPolicy`, data.HostConfig.RestartPolicy.Name)
	ctx.Request().Form().Set(`restartMaxRetryCount`, param.AsString(data.HostConfig.RestartPolicy.MaximumRetryCount))
}
