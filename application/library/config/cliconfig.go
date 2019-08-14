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

package config

import (
	"context"
	"errors"
	"io"
	stdLog "log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/cmd/event"
	"github.com/admpub/nging/application/library/caddy"
	"github.com/admpub/nging/application/library/cron"
	"github.com/admpub/nging/application/library/frp"
	"github.com/spf13/pflag"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func NewCLIConfig() *CLIConfig {
	return &CLIConfig{cmds: map[string]*exec.Cmd{}}
}

type CLIConfig struct {
	BackendDomain  string //前台绑定域名
	FrontendDomain string //后台绑定域名
	Address        string //监听IP地址
	Port           int    //监听端口
	Conf           string
	Confx          string
	Type           string //启动类型: webserver/ftpserver/manager
	Startup        string //manager启动时同时启动的服务，可选的有webserver/ftpserver,如有多个需用半角逗号“,”隔开
	cmds           map[string]*exec.Cmd
	caddyCh        *com.CmdChanReader
}

func (c *CLIConfig) InitFlag(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&c.Address, `address`, `a`, `0.0.0.0`, `address`)
	flagSet.IntVarP(&c.Port, `port`, `p`, 9999, `port`)
	flagSet.StringVarP(&c.Conf, `config`, `c`, filepath.Join(echo.Wd(), `config/config.yaml`), `config`)
	flagSet.StringVarP(&c.Confx, `subconfig`, `u`, filepath.Join(echo.Wd(), `config/config.frpserver.yaml`), `submodule config`)
	flagSet.StringVarP(&c.Type, `type`, `t`, `manager`, `operation type`)
	flagSet.StringVarP(&c.Startup, `startup`, `s`, `webserver,task,daemon,ftpserver,frpserver,frpclient`, `startup`)
	flagSet.StringVarP(&c.FrontendDomain, `frontenddomain`, `f`, ``, `frontend domain`)
	flagSet.StringVarP(&c.BackendDomain, `backenddomain`, `b`, ``, `backend domain`)
}

func (c *CLIConfig) OnlyRunServer() bool {
	switch c.Type {
	case `webserver`:
		caddy.TrapSignals()
		c.ParseConfig()
		DefaultConfig.Caddy.Init().Start()
		return true
	case `ftpserver`:
		c.ParseConfig()
		DefaultConfig.FTP.Init().Start()
		return true
	case `frpserver`:
		id := c.GenerateIDFromConfigFileName(c.Confx)
		err := frp.StartServerByConfigFile(c.Confx, c.FRPPidFile(id, true))
		if err != nil {
			stdLog.Println(err)
			os.Exit(1)
		}
		return true
	case `frpclient`:
		id := c.GenerateIDFromConfigFileName(c.Confx)
		err := frp.StartClientByConfigFile(c.Confx, c.FRPPidFile(id, false))
		if err != nil {
			stdLog.Println(err)
			os.Exit(1)
		}
		return true
	default:
		if c.Type == `official` || !event.SupportManager {
			c.Startup = `none`
			return false
		}
	}
	return false
}

func (c *CLIConfig) ParseConfig() {
	err := ParseConfig()
	if err != nil {
		if os.IsNotExist(err) {
			panic(err)
		}
		if IsInstalled() {
			MustOK(err)
		} else {
			log.Error(err)
		}
	}
}

//RunStartup manager启动时同时启动的服务
func (c *CLIConfig) RunStartup() {
	c.ParseConfig()
	c.Startup = strings.TrimSpace(c.Startup)
	if len(c.Startup) < 1 || !IsInstalled() {
		return
	}
	for _, serverType := range strings.Split(c.Startup, `,`) {
		serverType = strings.TrimSpace(serverType)
		switch serverType {
		case `webserver`:
			if err := DefaultCLIConfig.CaddyRestart(); err != nil {
				log.Error(err)
			}

		case `ftpserver`:
			if err := DefaultCLIConfig.FTPRestart(); err != nil {
				log.Error(err)
			}

		case `task`, `cron`: // 继续上次任务
			if err := cron.InitJobs(context.Background()); err != nil {
				log.Error(err)
			}

		case `daemon`:
			RunDaemon()

		case `frpserver`:
			if err := DefaultCLIConfig.FRPRestart(); err != nil {
				log.Error(err)
			}

		case `frpclient`:
			if err := DefaultCLIConfig.FRPClientRestart(); err != nil {
				log.Error(err)
			}
		}
	}
}

func (c *CLIConfig) CmdGroupStop(groupName string) error {
	if c.cmds == nil {
		return nil
	}
	groupName += `.`
	var err error
	for key, cmd := range c.cmds {
		if !strings.HasPrefix(key, groupName) {
			continue
		}
		err := c.Kill(cmd)
		if err != nil {
			log.Error(err)
		}
	}
	return err
}

func (c *CLIConfig) CmdHasGroup(groupName string) bool {
	if c.cmds == nil {
		return false
	}
	groupName += `.`
	for key := range c.cmds {
		if !strings.HasPrefix(key, groupName) {
			continue
		}
		return true
	}
	return false
}

func (c *CLIConfig) CmdGet(typeName string) *exec.Cmd {
	if c.cmds == nil {
		return nil
	}
	cmd, _ := c.cmds[typeName]
	return cmd
}

func (c *CLIConfig) CmdStop(typeName string) error {
	cmd := c.CmdGet(typeName)
	if cmd == nil {
		return nil
	}
	return c.Kill(cmd)
}

func (c *CLIConfig) CmdSendSignal(typeName string, sig os.Signal) error {
	cmd := c.CmdGet(typeName)
	if cmd == nil {
		return nil
	}
	if cmd.Process == nil {
		return nil
	}
	err := cmd.Process.Signal(sig)
	if err != nil && (cmd.ProcessState == nil || cmd.ProcessState.Exited()) {
		err = nil
	}
	return err
}

func (c *CLIConfig) Kill(cmd *exec.Cmd) error {
	return com.CloseProcessFromCmd(cmd)
}

var ErrCmdNotRunning = errors.New(`command is not running`)

func (c *CLIConfig) SetLogWriter(cmdType string, writer ...io.Writer) error {
	if c.cmds == nil {
		return ErrCmdNotRunning
	}
	cmd, ok := c.cmds[cmdType]
	if !ok || cmd == nil {
		return ErrCmdNotRunning
	}
	var wOut, wErr io.Writer
	length := len(writer)
	if length > 0 {
		wOut = writer[0]
		if length > 1 {
			wErr = writer[1]
		} else {
			wErr = wOut
		}
	} else {
		wOut = os.Stdout
		wErr = os.Stderr
	}
	cmd.Stdout = wOut
	cmd.Stderr = wErr
	return nil
}

func (c *CLIConfig) IsRunning(ct string) bool {
	if c.cmds == nil {
		return false
	}
	cmd, ok := c.cmds[ct]
	if !ok {
		return false
	}
	return CmdIsRunning(cmd)
}

func (c *CLIConfig) Reload(cts ...string) error {
	for _, ct := range cts {
		switch ct {
		case `caddy`:
			if c.IsRunning(`caddy`) {
				c.CaddyReload()
			}
		case `ftp`:
			if c.IsRunning(`ftp`) {
				c.FTPRestart()
			}
		default:
			for _, prefix := range []string{`frpserver.`, `frpclient.`} {
				if strings.HasPrefix(ct, prefix) {
					if c.IsRunning(ct) {
						c.FRPRestartID(strings.TrimPrefix(ct, prefix))
					}
				}
			}
		}
	}
	return nil
}
