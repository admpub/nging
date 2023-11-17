package request

import (
	"strings"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/middleware/bytes"
	"github.com/webx-top/echo/param"
)

type ContainerAdd struct {
	Name string // 容器名称

	Image      string `validate:"required"`
	Entrypoint string
	Command    string // 命令
	WorkingDir string
	Env        string
	Labels     string
	Memory     int64
	MemorySwap int64
	CpuWeight  int64

	ContainerPortFrom []string
	ContainerPortTo   []string
	ContainerPortNet  []string

	ContainerPathFrom []string
	ContainerPathTo   []string
	ContainerPathOp   []string
	VolumesFrom       []string

	Capabilities []string `form_seperator:","`

	Privileged bool
	AutoRemove bool
	StartNow   bool

	Tty                  string `validate:"oneof=Y N"`
	NetworkDisabled      string `validate:"oneof=Y N"`
	NetworkMode          string
	RestartPolicy        string `validate:"omitempty,oneof=no always on-failure unless-stopped"`
	RestartMaxRetryCount int

	container container.Config
	host      container.HostConfig
	network   network.NetworkingConfig
	platform  specs.Platform
}

func (c *ContainerAdd) AfterValidate(ctx echo.Context) error {
	if len(c.ContainerPathFrom) != len(c.ContainerPathTo) ||
		len(c.ContainerPathFrom) != len(c.ContainerPathOp) {
		return ctx.NewError(code.InvalidParameter, `存储卷配置参数不完整`).SetZone(`containerPathFrom`)
	}
	if len(c.ContainerPortFrom) != len(c.ContainerPortTo) ||
		len(c.ContainerPortFrom) != len(c.ContainerPortNet) {
		return ctx.NewError(code.InvalidParameter, `端口配置参数不完整`).SetZone(`containerPortFrom`)
	}
	c.host.Memory = c.Memory * bytes.MB
	c.host.MemorySwap = c.MemorySwap * bytes.MB
	if len(c.Command) > 0 {
		c.container.Cmd = com.ParseArgs(c.Command)
	}
	for _, ca := range c.Capabilities {
		if strings.HasPrefix(ca, `-`) {
			c.host.CapDrop = append(c.host.CapDrop, strings.TrimPrefix(ca, `-`))
		} else {
			c.host.CapAdd = append(c.host.CapAdd, strings.TrimPrefix(ca, `+`))
		}
	}
	if len(c.Entrypoint) > 0 {
		c.container.Entrypoint = com.ParseArgs(c.Entrypoint)
	}
	c.Env = strings.TrimSpace(c.Env)
	if len(c.Env) > 0 {
		com.SplitKVRowsCallback(c.Env, func(k, v string) error {
			c.container.Env = append(c.container.Env, k+`=`+v)
			return nil
		})
	}
	if len(c.Labels) > 0 {
		c.container.Labels = com.SplitKVRows(c.Labels)
	}
	c.container.WorkingDir = c.WorkingDir
	c.container.Image = c.Image
	c.host.Memory = c.Memory
	c.host.CPUShares = c.CpuWeight
	c.container.Tty = c.Tty == common.BoolY
	c.container.NetworkDisabled = c.NetworkDisabled == common.BoolY
	c.host.NetworkMode = container.NetworkMode(c.NetworkMode)
	if c.RestartPolicy != `on-failure` {
		c.RestartMaxRetryCount = 0
	}
	c.host.RestartPolicy = container.RestartPolicy{
		Name:              c.RestartPolicy,
		MaximumRetryCount: c.RestartMaxRetryCount,
	}
	c.host.Privileged = c.Privileged
	c.host.AutoRemove = c.AutoRemove
	c.host.VolumesFrom = param.StringSlice(c.VolumesFrom).Unique().Filter().String()
	var err error
	if ctx.Form(`portExport`) == common.BoolY {
		var ports []string
		var allowedProtos = []string{"tcp", "udp", "sctp"}
		for index, containerPort := range c.ContainerPortFrom {
			hostPort := c.ContainerPortTo[index]
			proto := c.ContainerPortNet[index]
			if com.StrIsNumeric(hostPort) {
				hostPort = `:` + hostPort
			}
			if !com.InSlice(proto, allowedProtos) {
				return ctx.NewError(code.InvalidParameter, `网络协议值无效`).SetZone(`containerPortNet`)
			}
			ports = append(ports, hostPort+`:`+containerPort+`/`+proto)
		}
		c.container.ExposedPorts, c.host.PortBindings, err = nat.ParsePortSpecs(ports)
		if err != nil {
			return err
		}
	}
	if ctx.Form(`storageVolumeMount`) == common.BoolY {
		var allowedOps = []string{"rw", "ro"}
		for index, containerPath := range c.ContainerPathFrom {
			hostPath := c.ContainerPathTo[index]
			op := c.ContainerPathOp[index]
			if !com.InSlice(op, allowedOps) {
				return ctx.NewError(code.InvalidParameter, `读写权限值无效`).SetZone(`containerPathOp`)
			}
			key := hostPath + `:` + containerPath + `:` + op
			c.host.Binds = append(c.host.Binds, key)
		}
	}
	return nil
}

func (c *ContainerAdd) Container() *container.Config {
	return &c.container
}

func (c *ContainerAdd) Host() *container.HostConfig {
	return &c.host
}

func (c *ContainerAdd) Network() *network.NetworkingConfig {
	return &c.network
}

func (c *ContainerAdd) Platform() *specs.Platform {
	return &c.platform
}

type ContainerEdit struct {
	Name                 string `validate:"required"`
	Memory               int64
	MemorySwap           int64
	CpuWeight            int64
	RestartPolicy        string `validate:"omitempty,oneof=no always on-failure unless-stopped"`
	RestartMaxRetryCount int

	resources container.Resources
}

func (c *ContainerEdit) AfterValidate(ctx echo.Context) error {
	if c.RestartPolicy != `on-failure` {
		c.RestartMaxRetryCount = 0
	}
	c.resources.CPUShares = c.CpuWeight
	c.resources.Memory = c.Memory * bytes.MB
	c.resources.MemorySwap = c.MemorySwap * bytes.MB
	return nil
}

func (c *ContainerEdit) Resources() container.Resources {
	return c.resources
}

type ContainerResize struct {
	Height uint `validate:"required,gt=0"`
	Width  uint `validate:"required,gt=0"`
}
