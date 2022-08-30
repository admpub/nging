package cmder

import (
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/admpub/log"
	"github.com/webx-top/com"

	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/config/cmder"

	"github.com/nging-plugins/caddymanager/application/library/caddy"
)

func Initer() interface{} {
	return &caddy.Config{}
}

func Get() cmder.Cmder {
	return cmder.Get(`caddy`)
}

func GetCaddyConfig() *caddy.Config {
	cm := cmder.Get(`caddy`).(*caddyCmd)
	return cm.CaddyConfig()
}

func New() cmder.Cmder {
	return &caddyCmd{
		CLIConfig: config.FromCLI(),
		once:      sync.Once{},
	}
}

type caddyCmd struct {
	CLIConfig   *config.CLIConfig
	caddyConfig *caddy.Config
	caddyCh     *com.CmdChanReader
	once        sync.Once
}

func (c *caddyCmd) Init() error {
	caddy.TrapSignals()
	return c.CaddyConfig().Init().Start()
}

func (c *caddyCmd) getConfig() *config.Config {
	if config.FromFile() == nil {
		c.CLIConfig.ParseConfig()
	}
	return config.FromFile()
}

func (c *caddyCmd) parseConfig() {
	c.caddyConfig, _ = c.getConfig().Extend.Get(`caddy`).(*caddy.Config)
	if c.caddyConfig == nil {
		c.caddyConfig = &caddy.Config{}
	}
	if len(c.caddyConfig.Caddyfile) == 0 {
		c.caddyConfig.Caddyfile = `./Caddyfile`
	} else if strings.HasSuffix(c.caddyConfig.Caddyfile, `/`) || strings.HasSuffix(c.caddyConfig.Caddyfile, `\`) {
		c.caddyConfig.Caddyfile = path.Join(c.caddyConfig.Caddyfile, `Caddyfile`)
	}
	caddy.SetDefaults(c.caddyConfig)
}

func (c *caddyCmd) CaddyConfig() *caddy.Config {
	c.once.Do(c.parseConfig)
	return c.caddyConfig
}

func (c *caddyCmd) StopHistory(_ ...string) error {
	if c.getConfig() == nil {
		return nil
	}
	return com.CloseProcessFromPidFile(c.CaddyConfig().PidFile)
}

func (c *caddyCmd) Start(writer ...io.Writer) error {
	err := c.StopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	params := []string{os.Args[0], `--config`, c.CLIConfig.Conf, `--type`, `caddy`}
	var cmd *exec.Cmd
	if caddy.EnableReload {
		c.caddyCh = com.NewCmdChanReader()
		cmd = com.RunCmdWithReaderWriter(params, c.caddyCh, writer...)
	} else {
		cmd = com.RunCmdWithWriter(params, writer...)
	}
	c.CLIConfig.CmdSet(`caddy`, cmd)
	return nil
}

func (c *caddyCmd) Stop() error {
	defer func() {
		if c.caddyCh != nil {
			c.caddyCh.Close()
			c.caddyCh = nil
		}
	}()
	return c.CLIConfig.CmdStop("caddy")
}

func (c *caddyCmd) Reload() error {
	if c.caddyCh == nil || com.IsWindows { //windows不支持重载，需要重启
		return c.Restart()
	}
	c.caddyCh.Send(com.BreakLine)
	return nil
	//c.CmdSendSignal("caddy", os.Interrupt)
}

func (c *caddyCmd) Restart(writer ...io.Writer) error {
	err := c.Stop()
	if err != nil {
		log.Error(err)
	}
	return c.Start(writer...)
}
