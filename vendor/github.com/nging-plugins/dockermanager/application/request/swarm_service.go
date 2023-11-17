package request

import (
	"strings"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var _ echo.ValueDecodersGetter = (*SwarmServiceAdd)(nil)
var _ echo.ValueDecodersGetter = (*SwarmServiceEdit)(nil)
var _ echo.FiltersGetter = (*SwarmServiceAdd)(nil)
var _ echo.FiltersGetter = (*SwarmServiceEdit)(nil)

var swarmServiceValueFilters = []echo.FormDataFilter{
	//echo.ExcludeFieldName(`TaskTemplate`),
}
var swarmServiceValueDecoders echo.BinderValueCustomDecoders = map[string]echo.BinderValueCustomDecoder{
	`Labels`:                             mapDecoder,
	`Networks.DriverOpts`:                mapDecoder,
	`TaskTemplate.LogDriver.Options`:     mapDecoder,
	`TaskTemplate.Networks.DriverOpts`:   mapDecoder,
	`TaskTemplate.ContainerSpec.Labels`:  mapDecoder,
	`TaskTemplate.ContainerSpec.Sysctls`: mapDecoder,
	`TaskTemplate.ContainerSpec.Command`: cmdDecoder,
	`TaskTemplate.ContainerSpec.Args`:    cmdDecoder,
	`TaskTemplate.ContainerSpec.Env`:     envDecoder,
}

type SwarmServiceAdd struct {
	types.ServiceCreateOptions
	swarm.ServiceSpec
}

func (a *SwarmServiceAdd) Filters(echo.Context) []echo.FormDataFilter {
	return swarmServiceValueFilters
}

func (a *SwarmServiceAdd) ValueDecoders(ctx echo.Context) echo.BinderValueCustomDecoders {
	r := echo.BinderValueCustomDecoders{}
	for k, v := range swarmServiceValueDecoders {
		r[k] = v
	}
	r[`EndpointSpec.Ports`] = makeSwarmServiceEndpointSpecPortsDecoder(ctx)
	r[`TaskTemplate.ContainerSpec.Configs`] = makeTaskTemplateContainerSpecConfigs(ctx)
	r[`TaskTemplate.ContainerSpec.Secrets`] = makeTaskTemplateContainerSpecSecrets(ctx)
	r[`TaskTemplate.ContainerSpec.Mounts`] = makeTaskTemplateContainerSpecMounts(ctx)
	return r
}

func (a *SwarmServiceAdd) AfterValidate(ctx echo.Context) error {
	if a.EndpointSpec != nil && len(a.EndpointSpec.Mode) == 0 {
		a.EndpointSpec = nil
	}
	if a.TaskTemplate.NetworkAttachmentSpec != nil && len(a.TaskTemplate.NetworkAttachmentSpec.ContainerID) == 0 {
		a.TaskTemplate.NetworkAttachmentSpec = nil
	}
	//a.UpdateConfig = nil
	//a.RollbackConfig = nil
	return nil
}

type SwarmServiceEdit struct {
	types.ServiceUpdateOptions
	swarm.ServiceSpec
}

func (a *SwarmServiceEdit) Filters(echo.Context) []echo.FormDataFilter {
	return swarmServiceValueFilters
}

func (a *SwarmServiceEdit) ValueDecoders(ctx echo.Context) echo.BinderValueCustomDecoders {
	r := echo.BinderValueCustomDecoders{}
	for k, v := range swarmServiceValueDecoders {
		r[k] = v
	}
	r[`EndpointSpec.Ports`] = makeSwarmServiceEndpointSpecPortsDecoder(ctx)
	r[`TaskTemplate.ContainerSpec.Configs`] = makeTaskTemplateContainerSpecConfigs(ctx)
	r[`TaskTemplate.ContainerSpec.Secrets`] = makeTaskTemplateContainerSpecSecrets(ctx)
	r[`TaskTemplate.ContainerSpec.Mounts`] = makeTaskTemplateContainerSpecMounts(ctx)
	return r
}

func (a *SwarmServiceEdit) FormNameFormatter(_ echo.Context) echo.FieldNameFormatter {
	return echo.MakeArrayFieldNameFormatter(com.LowerCaseFirst)
}

func (a *SwarmServiceEdit) ValueEncoders(ctx echo.Context) echo.BinderValueCustomEncoders {
	return echo.BinderValueCustomEncoders{
		`labels`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
		`networks[driverOpts]`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
		`taskTemplate[logDriver][options]`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
		`taskTemplate[networks][driverOpts]`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
		`taskTemplate[containerSpec][labels]`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
		`taskTemplate[containerSpec][sysctls]`: func(v interface{}) []string {
			return []string{com.JoinKVRows(v)}
		},
		`taskTemplate[containerSpec][command]`: func(v interface{}) []string {
			return []string{strings.Join(v.([]string), ` `)}
		},
		`taskTemplate[containerSpec][args]`: func(v interface{}) []string {
			return []string{strings.Join(v.([]string), ` `)}
		},
		`taskTemplate[containerSpec][env]`: func(v interface{}) []string {
			return []string{strings.Join(v.([]string), "\n")}
		},
		`taskTemplate[containerSpec][secrets]`: func(v interface{}) []string {
			secrets := v.([]*swarm.SecretReference)
			r := make([]string, len(secrets))
			for k, v := range secrets {
				r[k] = v.SecretID
			}
			return []string{strings.Join(r, ",")}
		},
		`taskTemplate[containerSpec][configs]`: func(v interface{}) []string {
			secrets := v.([]*swarm.ConfigReference)
			r := make([]string, len(secrets))
			for k, v := range secrets {
				r[k] = v.ConfigID
			}
			return []string{strings.Join(r, ",")}
		},
		`taskTemplate[containerSpec][mounts]`: func(v interface{}) []string {
			mounts := v.([]mount.Mount)
			f := ctx.Request().Form()
			for index, v := range mounts {
				echo.SetFormValue(f, `taskTemplate[containerSpec][mounts][type]`, index, v.Type)
				echo.SetFormValue(f, `taskTemplate[containerSpec][mounts][source]`, index, v.Source)
				echo.SetFormValue(f, `taskTemplate[containerSpec][mounts][target]`, index, v.Target)
				echo.SetFormValue(f, `taskTemplate[containerSpec][mounts][readOnly]`, index, common.BoolToFlag(v.ReadOnly))
				echo.SetFormValue(f, `taskTemplate[containerSpec][mounts][consistency]`, index, v.Consistency)
				if v.BindOptions != nil {
					echo.SetFormValue(f, `taskTemplate[containerSpec][mounts][bindOptions][propagation]`, index, v.BindOptions.Propagation)
					echo.SetFormValue(f, `taskTemplate[containerSpec][mounts][bindOptions][nonRecursive]`, index, common.BoolToFlag(v.BindOptions.NonRecursive))
				}
			}
			return nil
		},
		`endpointSpec[ports]`: func(v interface{}) []string {
			ports := v.([]swarm.PortConfig)
			f := ctx.Request().Form()
			for index, v := range ports {
				echo.SetFormValue(f, `endpointSpec[ports][name]`, index, v.Name)
				echo.SetFormValue(f, `endpointSpec[ports][protocol]`, index, v.Protocol)
				echo.SetFormValue(f, `endpointSpec[ports][targetPort]`, index, v.TargetPort)
				echo.SetFormValue(f, `endpointSpec[ports][publishedPort]`, index, v.PublishedPort)
				echo.SetFormValue(f, `endpointSpec[ports][publishMode]`, index, v.PublishMode)
			}
			return nil
		},
	}
}
