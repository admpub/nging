package cmder

import (
	"io"

	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/config/cmder"
	conf "github.com/nging-plugins/servermanager/application/library/config"
)

func init() {
	cmder.Register(`daemon`, New())
	config.DefaultStartup += ",daemon"
}

func Get() cmder.Cmder {
	return cmder.Get(`daemon`)
}

func New() cmder.Cmder {
	return &daemonCmd{
		CLIConfig: config.FromCLI(),
	}
}

type daemonCmd struct {
	CLIConfig *config.CLIConfig
}

func (c *daemonCmd) Init() error {
	return nil
}

func (c *daemonCmd) getConfig() *config.Config {
	if config.FromFile() == nil {
		c.CLIConfig.ParseConfig()
	}
	return config.FromFile()
}

func (c *daemonCmd) StopHistory(_ ...string) error {
	return nil
}

func (c *daemonCmd) Start(writer ...io.Writer) error {
	conf.RunDaemon()
	return nil
}

func (c *daemonCmd) Stop() error {
	return nil
}

func (c *daemonCmd) Reload() error {
	return nil
}

func (c *daemonCmd) Restart(writer ...io.Writer) error {
	conf.RestartDaemon()
	return nil
}
