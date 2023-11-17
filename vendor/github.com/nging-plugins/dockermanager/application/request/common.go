package request

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func mapDecoder(values []string) (interface{}, error) {
	return com.SplitKVRows(values[0]), nil
}

func cmdDecoder(values []string) (interface{}, error) {
	return com.ParseArgs(values[0]), nil
}

func envDecoder(values []string) (interface{}, error) {
	var env []string
	com.SplitKVRowsCallback(values[0], func(k, v string) error {
		env = append(env, k+`=`+v)
		return nil
	})
	return env, nil
}

func makeTaskTemplateContainerSpecSecrets(ctx echo.Context) func(values []string) (interface{}, error) {
	secretIDs := param.Split(ctx.Form(`taskTemplate[containerSpec][secrets]`), `,`).Unique().Filter().String()
	secretNames := param.Split(ctx.Form(`taskTemplate[containerSpec][secrets]_text`), `,`).Unique().Filter().String()
	return func(values []string) (interface{}, error) {
		r := make([]*swarm.SecretReference, len(secretIDs))
		for index, secretID := range secretIDs {
			s := &swarm.SecretReference{
				SecretID:   secretID,
				SecretName: ``,
			}
			if index < len(secretNames) {
				s.SecretName = secretNames[index]
			}
			r[index] = s
		}
		return r, nil
	}
}

func makeTaskTemplateContainerSpecConfigs(ctx echo.Context) func(values []string) (interface{}, error) {
	configIDs := param.Split(ctx.Form(`taskTemplate[containerSpec][configs]`), `,`).Unique().Filter().String()
	configNames := param.Split(ctx.Form(`taskTemplate[containerSpec][configs]_text`), `,`).Unique().Filter().String()
	return func(values []string) (interface{}, error) {
		r := make([]*swarm.ConfigReference, len(configIDs))
		for index, secretID := range configIDs {
			s := &swarm.ConfigReference{
				ConfigID:   secretID,
				ConfigName: ``,
			}
			if index < len(configNames) {
				s.ConfigName = configNames[index]
			}
			r[index] = s
		}
		return r, nil
	}
}

func makeTaskTemplateContainerSpecMounts(ctx echo.Context) func(values []string) (interface{}, error) {
	return func(values []string) (interface{}, error) {
		types := ctx.FormValues(`taskTemplate[containerSpec][mounts][type]`)
		sources := ctx.FormValues(`taskTemplate[containerSpec][mounts][source]`)
		targets := ctx.FormValues(`taskTemplate[containerSpec][mounts][target]`)
		readOnlyList := ctx.FormValues(`taskTemplate[containerSpec][mounts][readOnly]`)
		consistencyList := ctx.FormValues(`taskTemplate[containerSpec][mounts][consistency]`)
		propagations := ctx.FormValues(`taskTemplate[containerSpec][mounts][bindOptions][propagation]`)
		nonRecursives := ctx.FormValues(`taskTemplate[containerSpec][mounts][bindOptions][nonRecursive]`)
		if len(types) != len(sources) || len(types) != len(targets) || len(types) != len(readOnlyList) || len(types) != len(consistencyList) || len(types) != len(propagations) || len(types) != len(nonRecursives) {
			return nil, ctx.NewError(code.InvalidParameter, `提交的挂载卷参数不完整`).SetZone(`taskTemplate[containerSpec][mounts]`)
		}
		r := make([]mount.Mount, len(types))
		for index, typ := range types {
			m := mount.Mount{
				Type:        mount.Type(typ),
				Source:      sources[index],
				Target:      targets[index],
				ReadOnly:    readOnlyList[index] == `Y`,
				Consistency: mount.Consistency(consistencyList[index]),
			}
			m.BindOptions = &mount.BindOptions{
				Propagation:  mount.Propagation(propagations[index]),
				NonRecursive: nonRecursives[index] == `Y`,
			}
			r[index] = m
		}
		return r, nil
	}
}

func makeSwarmServiceEndpointSpecPortsDecoder(ctx echo.Context) func(values []string) (interface{}, error) {
	return func(values []string) (interface{}, error) {
		mode := ctx.Form(`endpointSpec[mode]`)
		if len(mode) == 0 {
			return []swarm.PortConfig{}, nil
		}
		if !com.InSlice(mode, []string{`vip`, `dnsrr`}) {
			return nil, ctx.NewError(code.InvalidParameter, `入口模式值无效`).SetZone(`endpointSpec[mode]`)
		}
		names := ctx.FormValues(`endpointSpec[ports][name]`)
		protocols := ctx.FormValues(`endpointSpec[ports][protocol]`)
		targetPorts := ctx.FormValues(`endpointSpec[ports][targetPort]`)
		publishedPorts := ctx.FormValues(`endpointSpec[ports][publishedPort]`)
		publishModes := ctx.FormValues(`endpointSpec[ports][publishMode]`)
		if len(names) != len(protocols) || len(names) != len(targetPorts) || len(names) != len(publishedPorts) || len(names) != len(publishModes) {
			return nil, ctx.NewError(code.InvalidParameter, `提交的端口映射参数不完整`).SetZone(`endpointSpec[ports]`)
		}
		ports := make([]swarm.PortConfig, 0, len(names))
		for index, name := range names {
			protocol := protocols[index]
			targetPort := param.AsUint32(targetPorts[index])
			publishedPort := param.AsUint32(publishedPorts[index])
			publishMode := publishModes[index]
			if len(name) == 0 || targetPort == 0 || publishedPort == 0 {
				continue
			}
			if !com.InSlice(protocol, []string{`tcp`, `udp`, `stcp`}) {
				return nil, ctx.NewError(code.InvalidParameter, `网络协议无效`).SetZone(`endpointSpec[ports][protocol]`)
			}
			if !com.InSlice(publishMode, []string{`ingress`, `host`}) {
				return nil, ctx.NewError(code.InvalidParameter, `发布模式无效`).SetZone(`endpointSpec[ports][publishMode]`)
			}
			cfg := swarm.PortConfig{
				Name:          name,
				Protocol:      swarm.PortConfigProtocol(protocol),
				TargetPort:    targetPort,
				PublishedPort: publishedPort,
				PublishMode:   swarm.PortConfigPublishMode(publishMode),
			}
			ports = append(ports, cfg)
		}
		return ports, nil
	}
}
