/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/caddyserver/certmagic"
	"github.com/spf13/cobra"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/engine/standard"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/subdomains"

	"github.com/admpub/events/emitter"
	figure "github.com/admpub/go-figure"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/cmd/event"
	"github.com/admpub/nging/application/handler/setup"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/config/startup"
	"github.com/admpub/nging/application/library/license"
	"github.com/admpub/nging/application/library/msgbox"
)

// Nging 启动入口

// rootCmd represents the base command when called without any subcommands
var rootCmd = NewRoot()

func NewRoot() *cobra.Command {
	// rootCmd represents the base command when called without any subcommands
	return &cobra.Command{
		Use:          filepath.Base(os.Args[0]),
		Short:        ``,
		Long:         ``,
		SilenceUsage: true,
		RunE:         rootRunE,
	}
}

func rootRunE(cmd *cobra.Command, args []string) error {
	if !event.Licensed {
		machineID, _ := license.MachineID()
		message := `Invalid license!
授权无效!
		
If you have already purchased a license, please place the ` + license.FileName() + ` to:
如果您已经购买了授权，请将授权证书` + license.FileName() + `放到：
		
%s`
		msgbox.Error("WARNING",
			message,
			license.FilePath())

		fmt.Println(``)
		fmt.Println(`To purchase a license, please go to our official website:`)
		fmt.Println(`购买授权请前往官方网站：`)
		fmt.Println(``)
		fmt.Println(license.ProductURL() + `?version=` + config.Version.Number + `&machineID=` + machineID)
		if event.MustLicensed {
			return nil
		}
	}

	//独立模块
	if config.DefaultCLIConfig.OnlyRunServer() {
		return nil
	}

	//Manager
	config.DefaultCLIConfig.RunStartup()

	if config.IsInstalled() {
		if err := setup.Upgrade(); err != nil && os.ErrNotExist != err {
			log.Error(`upgrade.sql: `, err)
		}
	}

	// LOGO
	fmt.Println(strings.TrimSuffix(figure.NewFigure(event.SoftwareName, `big`, false).String(), "\n"), config.Version.VString()+"\n")

	event.Start()
	startup.FireBefore(`web`)
	defer startup.FireAfter(`web`)

	c := &engine.Config{
		ReusePort:   true,
		TLSAuto:     config.DefaultConfig.Sys.SSLAuto,
		TLSEmail:    config.DefaultConfig.Sys.SSLEmail,
		TLSHosts:    config.DefaultConfig.Sys.SSLHosts,
		TLSCacheDir: config.DefaultConfig.Sys.SSLCacheDir,
		TLSCertFile: config.DefaultConfig.Sys.SSLCertFile,
		TLSKeyFile:  config.DefaultConfig.Sys.SSLKeyFile,
	}
	c.Address = fmt.Sprintf(`%s:%v`, config.DefaultCLIConfig.Address, config.DefaultCLIConfig.Port)
	hasCert := (len(c.TLSCertFile) > 0 && len(c.TLSKeyFile) > 0)
	//c.TLSAuto = true
	if c.TLSAuto || hasCert {
		if config.DefaultCLIConfig.Port == 80 {
			if c.TLSAuto {
				echo.PanicIf(initCertMagic(c))
				//c.SupportAutoTLS(nil, config.DefaultConfig.Sys.SSLHosts...)
			} else {
				c.Address = fmt.Sprintf(`%s:443`, config.DefaultCLIConfig.Address)
				e2 := echo.New()
				e2.Use(middleware.HTTPSRedirect(), middleware.Log(), middleware.Recover())
				go e2.Run(standard.New(fmt.Sprintf(`%s:80`, config.DefaultCLIConfig.Address)))
			}
			subdomains.Default.Protocol = `https`
		}
	}
	if len(event.Welcome) > 0 {
		now := time.Now()
		msgbox.Success(`Welcome`,
			event.Welcome,
			config.Version.VString(),
			now.Format("Monday, 02 Jan 2006"))
	}
	subdomains.Default.SetDebug(config.DefaultConfig.Debug)
	emitter.DefaultCondEmitter.Fire(`beforeRun`, -1)
	subdomains.Default.Run(standard.NewWithConfig(c))
	return c.Listener.Close()
}

func initCertMagic(c *engine.Config) error {
	fileStorage := &certmagic.FileStorage{
		Path: filepath.Join(echo.Wd(), `data`, `cache`, `certmagic`),
	}
	if err := os.MkdirAll(fileStorage.Path, 0777); err != nil {
		return err
	}
	if event.Develop { // use the staging endpoint while we're developing
		certmagic.DefaultACME.CA = certmagic.LetsEncryptStagingCA
	} else {
		certmagic.DefaultACME.CA = certmagic.LetsEncryptProductionCA
	}
	certmagic.DefaultACME.Email = c.TLSEmail
	certmagic.DefaultACME.Agreed = true
	certmagic.Default.Storage = fileStorage
	ln, err := certmagic.Listen(c.TLSHosts)
	if err == nil {
		c.SetListener(ln)
	}
	return err
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	Init()
	if len(rootCmd.Use) == 0 {
		rootCmd.Use = os.Args[0]
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	config.DefaultCLIConfig.InitFlag(rootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
}
