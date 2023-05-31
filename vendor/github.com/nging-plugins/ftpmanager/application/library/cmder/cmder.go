package cmder

import (
	"io"
	"os"

	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/defaults"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/cmder"
	"github.com/admpub/nging/v5/application/library/config/extend"

	"github.com/nging-plugins/ftpmanager/application/library/ftp"
	"github.com/nging-plugins/ftpmanager/application/model"
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

func GetFTPConfig() *ftp.Config {
	cm := cmder.Get(`ftpserver`).(*ftpCmd)
	return cm.FTPConfig()
}

func StartOnce() {
	if config.FromCLI().IsRunning(`ftpserver`) {
		return
	}
	Get().Start()
}

func Stop() {
	if !config.FromCLI().IsRunning(`ftpserver`) {
		return
	}
	Get().Stop()
}

func New() cmder.Cmder {
	return &ftpCmd{
		CLIConfig: config.FromCLI(),
		once:      once.Once{},
	}
}

type ftpCmd struct {
	CLIConfig *config.CLIConfig
	ftpConfig *ftp.Config
	once      once.Once
}

func (c *ftpCmd) Boot() error {
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
	ctx := defaults.NewMockContext()
	userM := model.NewFtpUser(ctx)
	exists, err := userM.ExistsAvailable()
	if err != nil {
		log.Error(err.Error())
	}
	if !exists { // 没有有效用户时无需启动
		return nil
	}
	params := []string{os.Args[0], `--config`, c.CLIConfig.Conf, `--type`, `ftpserver`}
	cmd := com.RunCmdWithWriter(params, writer...)
	c.CLIConfig.CmdSet(`ftpserver`, cmd)
	return nil
}

func (c *ftpCmd) Stop() error {
	c.FTPConfig().Stop()
	return c.CLIConfig.CmdStop("ftpserver")
}

func (c *ftpCmd) Reload() error {
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

func (c *ftpCmd) Restart(writer ...io.Writer) error {
	err := c.Stop()
	if err != nil {
		log.Error(err)
	}
	return c.Start(writer...)
}
