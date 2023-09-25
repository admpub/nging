package cmder

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/webx-top/com"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/cmder"

	"github.com/nging-plugins/caddymanager/application/library/caddy"
)

const Name = `caddy`

func Initer() interface{} {
	return &caddy.Config{}
}

func Get() cmder.Cmder {
	return cmder.Get(Name)
}

func GetCaddyConfig() *caddy.Config {
	return GetCaddyCmd().CaddyConfig()
}

func GetCaddyCmd() *caddyCmd {
	cm := cmder.Get(Name).(*caddyCmd)
	return cm
}

func New() cmder.Cmder {
	return &caddyCmd{
		CLIConfig: config.FromCLI(),
		once:      once.Once{},
	}
}

type caddyCmd struct {
	CLIConfig   *config.CLIConfig
	caddyConfig *caddy.Config
	caddyCh     *com.CmdChanReader
	once        once.Once
}

func (c *caddyCmd) Boot() error {
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
	c.caddyConfig, _ = c.getConfig().Extend.Get(Name).(*caddy.Config)
	if c.caddyConfig == nil {
		c.caddyConfig = &caddy.Config{}
	}
	if len(c.caddyConfig.Caddyfile) > 0 && (strings.HasSuffix(c.caddyConfig.Caddyfile, `/`) || strings.HasSuffix(c.caddyConfig.Caddyfile, `\`)) {
		c.caddyConfig.Caddyfile = filepath.Join(c.caddyConfig.Caddyfile, `Caddyfile`)
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
	params := []string{os.Args[0], `--config`, c.CLIConfig.Conf, `--type`, Name}
	var cmd *exec.Cmd
	if caddy.EnableReload {
		c.caddyCh = com.NewCmdChanReader()
		cmd = com.RunCmdWithReaderWriter(params, c.caddyCh, writer...)
	} else {
		cmd = com.RunCmdWithWriter(params, writer...)
	}
	c.CLIConfig.CmdSet(Name, cmd)
	return nil
}

func (c *caddyCmd) Stop() error {
	defer func() {
		if c.caddyCh != nil {
			c.caddyCh.Close()
			c.caddyCh = nil
		}
	}()
	return c.CLIConfig.CmdStop(Name)
}

func (c *caddyCmd) Reload() error {
	err := c.Stop()
	if err != nil {
		log.Error(err)
	}
	err = c.StopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	c.once.Reset()
	return c.Start()
}

func (c *caddyCmd) ReloadServer() error {
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
