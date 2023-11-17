package request

import (
	"github.com/webx-top/echo"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

var _ echo.ValueDecodersGetter = (*NetworkAdd)(nil)
var _ echo.ValueDecodersGetter = (*NetworkConnect)(nil)
var _ echo.AfterValidate = (*NetworkAdd)(nil)

type NetworkAdd struct {
	Name string `validate:"required"`
	types.NetworkCreate
}

func (a *NetworkAdd) ValueDecoders(echo.Context) echo.BinderValueCustomDecoders {
	return map[string]echo.BinderValueCustomDecoder{
		`Options`:                mapDecoder,
		`IPAM.Options`:           mapDecoder,
		`IPAM.Config.AuxAddress`: mapDecoder,
		`Labels`:                 mapDecoder,
	}
}

func (c *NetworkAdd) AfterValidate(ctx echo.Context) error {
	if c.IPAM != nil {
		ipam := make([]network.IPAMConfig, 0, len(c.IPAM.Config))
		for _, cfg := range c.IPAM.Config {
			if len(cfg.Subnet) > 0 {
				ipam = append(ipam, cfg)
			}
		}
		c.IPAM.Config = ipam
	}
	return nil
}

type NetworkConnect struct {
	NetworkID   string `validate:"required"`
	ContainerID string `validate:"required"`
	network.EndpointSettings
}

func (a *NetworkConnect) ValueDecoders(echo.Context) echo.BinderValueCustomDecoders {
	return map[string]echo.BinderValueCustomDecoder{
		`DriverOpts`: mapDecoder,
	}
}

type NetworkDisconnect struct {
	NetworkID   string `validate:"required"`
	ContainerID string `validate:"required"`
	Force       bool
}
