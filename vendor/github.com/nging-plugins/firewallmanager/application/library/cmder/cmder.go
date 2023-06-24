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

package cmder

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/admpub/gerberos"
	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/config/cmder"
	"github.com/admpub/nging/v5/application/library/config/extend"

	firewallConfig "github.com/nging-plugins/firewallmanager/application/library/config"
	"github.com/nging-plugins/firewallmanager/application/library/firewall"
	"github.com/nging-plugins/firewallmanager/application/model"
)

const Name = `firewall`
const DefaultPidFile = `firewall.pid`
const DefaultChainName = `NgingDynamic`
const DefaultTable4Name = `nging_dynamic_ip4`
const DefaultTable6Name = `nging_dynamic_ip6`

func init() {
	cmder.Register(Name, New())
	config.DefaultStartup += "," + Name
	extend.Register(Name, func() interface{} {
		return &firewallConfig.Config{}
	})
	gerberos.DefaultChainName = DefaultChainName
	gerberos.DefaultTable4Name = DefaultTable4Name
	gerberos.DefaultTable6Name = DefaultTable6Name
}

func Initer() interface{} {
	return &firewallConfig.Config{}
}

func Get() cmder.Cmder {
	return cmder.Get(Name)
}

func GetFirewallConfig() *firewallConfig.Config {
	cm := cmder.Get(Name).(*firewallCmd)
	return cm.FirewallConfig()
}

func StartOnce(writer ...io.Writer) {
	if config.FromCLI().IsRunning(Name) {
		return
	}
	Get().Start(writer...)
}

func Stop() {
	if !config.FromCLI().IsRunning(Name) {
		return
	}
	Get().Stop()
}

func New() cmder.Cmder {
	return &firewallCmd{
		CLIConfig: config.FromCLI(),
		once:      once.Once{},
	}
}

type firewallCmd struct {
	CLIConfig      *config.CLIConfig
	firewallConfig *firewallConfig.Config
	pidFile        string
	once           once.Once
}

func (c *firewallCmd) PidFile() string {
	c.FirewallConfig()
	return c.pidFile
}

func (c *firewallCmd) boot() error {
	cfg := c.FirewallConfig()
	err := com.WritePidFile(c.pidFile)
	if err != nil {
		log.Error(err.Error())
	}

	gerberosCfg := &gerberos.Configuration{
		Verbose:      cfg.Verbose,
		SaveFilePath: cfg.SaveFilePath,
		Rules:        map[string]*gerberos.Rule{},
	}
	switch cfg.Backend {
	case `nftables`:
		gerberosCfg.Backend = `nft`
	case `iptables`:
		gerberosCfg.Backend = `ipset`
	case `nft`, `ipset`:
		gerberosCfg.Backend = cfg.Backend
	case ``:
		backends := firewall.DynamicRuleBackends.Slice()
		if len(backends) > 0 {
			gerberosCfg.Backend = backends[0].K
		}
	default:
		return fmt.Errorf(`unsupported firewall backend: %q`, cfg.Backend)
	}
	//gerberosCfg.Verbose = true
	ctx := defaults.NewMockContext()
	ruleM := model.NewRuleDynamic(ctx)
	_, err = ruleM.ListByOffset(nil, nil, 0, -1, `disabled`, `N`)
	if err != nil {
		return err
	}

	for _, row := range ruleM.Objects() {
		rule, err := firewall.DynamicRuleFromDB(ctx, row)
		if err != nil {
			log.Error(err.Error())
		} else {
			gerberosCfg.Rules[param.AsString(row.Id)] = &rule
		}
	}

	// Runner
	rn := gerberos.NewRunner(gerberosCfg)
	if err := rn.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize runner: %s", err)
	}
	defer func() {
		if err := rn.Finalize(); err != nil {
			log.Fatalf("failed to finalize runner: %s", err)
		}
	}()
	rn.Run(true)
	return err
}

func (c *firewallCmd) getConfig() *config.Config {
	if config.FromFile() == nil {
		c.CLIConfig.ParseConfig()
	}
	return config.FromFile()
}

func (c *firewallCmd) parseConfig() {
	c.firewallConfig, _ = c.getConfig().Extend.Get(Name).(*firewallConfig.Config)
	if c.firewallConfig == nil {
		c.firewallConfig = &firewallConfig.Config{}
	}
	pidFile := filepath.Join(echo.Wd(), `data/pid/`+Name)
	err := com.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		log.Error(err)
	}
	pidFile = filepath.Join(pidFile, DefaultPidFile)
	c.pidFile = pidFile
}

func (c *firewallCmd) FirewallConfig() *firewallConfig.Config {
	c.once.Do(c.parseConfig)
	return c.firewallConfig
}

func (c *firewallCmd) StopHistory(_ ...string) error {
	if c.getConfig() == nil {
		return nil
	}
	return com.CloseProcessFromPidFile(c.PidFile())
}

func (c *firewallCmd) Start(writer ...io.Writer) error {
	err := c.StopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	ctx := defaults.NewMockContext()
	ruleM := model.NewRuleDynamic(ctx)
	exists, err := ruleM.ExistsAvailable()
	if err != nil {
		log.Error(err.Error())
	}
	if !exists { // 没有有效用户时无需启动
		return nil
	}
	params := []string{os.Args[0], `--config`, c.CLIConfig.Conf, `--type`, Name}
	cmd := com.RunCmdWithWriter(params, writer...)
	c.CLIConfig.CmdSet(Name, cmd)
	return nil
}

func (c *firewallCmd) Stop() error {
	c.CLIConfig.CmdSendSignal(Name, syscall.SIGINT)
	time.Sleep(time.Second)
	return c.CLIConfig.CmdStop(Name)
}

func (c *firewallCmd) Reload() error {
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

func (c *firewallCmd) Restart(writer ...io.Writer) error {
	err := c.Stop()
	if err != nil {
		log.Error(err)
	}
	return c.Start(writer...)
}
