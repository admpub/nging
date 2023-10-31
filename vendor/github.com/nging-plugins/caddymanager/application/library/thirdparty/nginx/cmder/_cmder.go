package cmder

import (
	"context"
	"io"

	"github.com/admpub/once"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/cmder"

	"github.com/nging-plugins/caddymanager/application/library/thirdparty"
	nginxConfig "github.com/nging-plugins/caddymanager/application/library/thirdparty/nginx/config"
)

func Initer() interface{} {
	return &nginxConfig.Config{}
}

func Get() cmder.Cmder {
	return cmder.Get(nginxConfig.Name)
}

func GetNginxConfig() *nginxConfig.Config {
	return GetNginxCmd().NginxConfig()
}

func GetNginxCmd() *nginxCmd {
	cm := cmder.Get(nginxConfig.Name).(*nginxCmd)
	return cm
}

func New() cmder.Cmder {
	return &nginxCmd{
		CLIConfig: config.FromCLI(),
		once:      once.Once{},
	}
}

type nginxCmd struct {
	CLIConfig   *config.CLIConfig
	nginxConfig *nginxConfig.Config
	once        once.Once
}

func (c *nginxCmd) Boot() error {
	err := c.NginxConfig().Init()
	if err != nil {
		return err
	}
	return c.NginxConfig().Start(context.Background())
}

func (c *nginxCmd) getConfig() *config.Config {
	if config.FromFile() == nil {
		c.CLIConfig.ParseConfig()
	}
	return config.FromFile()
}

func (c *nginxCmd) parseConfig() {
	c.nginxConfig, _ = c.getConfig().Extend.Get(nginxConfig.Name).(*nginxConfig.Config)
	if c.nginxConfig == nil {
		c.nginxConfig = &nginxConfig.Config{}
	}
}

func (c *nginxCmd) NginxConfig() *nginxConfig.Config {
	c.once.Do(c.parseConfig)
	return c.nginxConfig
}

func (c *nginxCmd) StopHistory(_ ...string) error {
	return c.Reload()
}

func (c *nginxCmd) Start(writer ...io.Writer) error {
	return c.NginxConfig().Start(thirdparty.WithStdoutStderr(context.Background(), writer...))
}

func (c *nginxCmd) Stop() error {
	return c.NginxConfig().Stop(context.Background())
}

func (c *nginxCmd) Reload() error {
	return c.NginxConfig().Reload(context.Background())
}

func (c *nginxCmd) Restart(writer ...io.Writer) error {
	return c.NginxConfig().Reload(thirdparty.WithStdoutStderr(context.Background(), writer...))
}
