package cmder

import (
	"io"
	"os"
	"sync"

	"github.com/admpub/log"
	"github.com/webx-top/com"

	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/config/cmder"
	"github.com/admpub/nging/v4/application/library/config/extend"

	"github.com/nging-plugins/ftpmanager/application/library/ftp"
)

func init() {
	cmder.Register(`ftpserver`, New())
	config.DefaultStartup += ",ftpserver"
	extend.Register(`ftpserver`, func() interface{} {
		return &ftp.Config{}
	})
}

func Initer() interface{} {
	return &ftp.Config{}
}

func Get() cmder.Cmder {
	return cmder.Get(`ftpserver`)
}

func GetCaddyConfig() *ftp.Config {
	cm := cmder.Get(`ftpserver`).(*ftpCmd)
	return cm.FTPConfig()
}

func New() cmder.Cmder {
	return &ftpCmd{
		CLIConfig: config.FromCLI(),
		once:      sync.Once{},
	}
}

type ftpCmd struct {
	CLIConfig *config.CLIConfig
	ftpConfig *ftp.Config
	once      sync.Once
}

func (c *ftpCmd) Init() error {
	return c.FTPConfig().Init().Start()
}

func (c *ftpCmd) getConfig() *config.Config {
	if config.FromFile() == nil {
		c.CLIConfig.ParseConfig()
	}
	return config.FromFile()
}

func (c *ftpCmd) parseConfig() {
	c.ftpConfig, _ = c.getConfig().Extend.Get(`ftpserver`).(*ftp.Config)
	if c.ftpConfig == nil {
		c.ftpConfig = &ftp.Config{}
	}
	ftp.SetDefaults(c.ftpConfig)
}

func (c *ftpCmd) FTPConfig() *ftp.Config {
	c.once.Do(c.parseConfig)
	return c.ftpConfig
}

func (c *ftpCmd) StopHistory(_ ...string) error {
	if c.getConfig() == nil {
		return nil
	}
	return com.CloseProcessFromPidFile(c.FTPConfig().PidFile)
}

func (c *ftpCmd) Start(writer ...io.Writer) error {
	err := c.StopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	params := []string{os.Args[0], `--config`, c.CLIConfig.Conf, `--type`, `ftpserver`}
	cmd := com.RunCmdWithWriter(params, writer...)
	c.CLIConfig.CmdSet(`ftpserver`, cmd)
	return nil
}

func (c *ftpCmd) Stop() error {
	return c.CLIConfig.CmdStop("ftpserver")
}

func (c *ftpCmd) Reload() error {
	return nil
}

func (c *ftpCmd) Restart(writer ...io.Writer) error {
	err := c.Stop()
	if err != nil {
		log.Error(err)
	}
	return c.Start(writer...)
}
