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
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/cmd/event"
	"github.com/admpub/nging/v4/application/library/config/cmder"
	"github.com/admpub/nging/v4/application/library/config/startup"
	"github.com/spf13/pflag"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func NewCLIConfig() *CLIConfig {
	cli := &CLIConfig{
		cmds:     map[string]*exec.Cmd{},
		pid:      os.Getpid(),
		envFiles: findEnvFile(),
	}
	cli.InitEnviron()
	//cli.WatchEnvConfig()
	return cli
}

var DefaultStartup = `webserver,task,daemon,ftpserver,frpserver,frpclient`
var DefaultPort = 9999

type CLIConfig struct {
	BackendDomain  string //后台绑定域名
	FrontendDomain string //前台绑定域名
	Address        string //监听IP地址
	Port           int    //监听端口
	Conf           string
	Confx          string
	Type           string //启动类型: webserver/ftpserver/manager
	Startup        string //manager启动时同时启动的服务，可选的有webserver/ftpserver,如有多个需用半角逗号“,”隔开
	cmds           map[string]*exec.Cmd
	pid            int
	envVars        map[string]string // 从文件“.env”中读取到的自定义环境变量（对于系统环境变量中已经存在的变量，会自动忽略）
	envFiles       []string
	envMonitor     *com.MonitorEvent
	envLock        sync.RWMutex
}

func (c *CLIConfig) InitFlag(flagSet *pflag.FlagSet) {
	flagSet.StringVarP(&c.Address, `address`, `a`, `0.0.0.0`, `address`)
	flagSet.IntVarP(&c.Port, `port`, `p`, DefaultPort, `port`)
	flagSet.StringVarP(&c.Conf, `config`, `c`, filepath.Join(echo.Wd(), `config/config.yaml`), `config`)
	flagSet.StringVarP(&c.Confx, `subconfig`, `u`, filepath.Join(echo.Wd(), `config/config.frpserver.yaml`), `submodule config`)
	flagSet.StringVarP(&c.Type, `type`, `t`, `manager`, `operation type`)
	flagSet.StringVarP(&c.Startup, `startup`, `s`, DefaultStartup, `startup`)
	flagSet.StringVarP(&c.FrontendDomain, `frontend.domain`, `f`, ``, `frontend domain`)
	flagSet.StringVarP(&c.BackendDomain, `backend.domain`, `b`, ``, `backend domain`)
}

func (c *CLIConfig) OnlyRunServer() bool {
	cm := cmder.Get(c.Type)
	if cm != nil {
		startup.FireBefore(c.Type)
		err := cm.Init()
		if err != nil {
			com.ExitOnFailure(err.Error())
		}
		startup.FireAfter(c.Type)
		return true
	}

	// manager mode
	if c.Type == `official` || !event.SupportManager {
		c.Startup = `none`
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
	for _, serverType := range param.StringSlice(strings.Split(c.Startup, `,`)).Unique() {
		serverType = strings.TrimSpace(serverType)
		cm := cmder.Get(serverType)
		if cm != nil {
			startup.FireBefore(serverType)
			if err := cm.Restart(); err != nil {
				log.Error(err)
			}
			startup.FireAfter(serverType)
			continue
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

func (c *CLIConfig) CmdSet(name string, cmd *exec.Cmd) {
	c.cmds[name] = cmd
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
	return c.SendSignal(cmd, sig)
}

func (c *CLIConfig) SendSignal(cmd *exec.Cmd, sig os.Signal) error {
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

type ReloadByConfig interface {
	ReloadByConfig(cfg *Config, args ...string) error
}

func (c *CLIConfig) Reload(cfg *Config, cts ...string) error {
	for _, ct := range cts {
		if !c.IsRunning(ct) {
			continue
		}
		typeAndID := strings.SplitN(ct, ".", 2)
		serverType := typeAndID[0]
		cm := cmder.Get(serverType)
		if cm != nil {
			if rd, ok := cm.(ReloadByConfig); ok {
				if len(typeAndID) == 2 {
					rd.ReloadByConfig(cfg, typeAndID[1])
					continue
				}
				rd.ReloadByConfig(cfg)
				continue
			}
			if len(typeAndID) == 1 {
				cm.Reload()
				continue
			}
			if rd, ok := cm.(cmder.RestartBy); ok {
				if len(typeAndID) == 2 {
					rd.RestartBy(typeAndID[1])
					continue
				}
			}
			cm.Reload()
			continue
		}
	}
	return nil
}

func (c *CLIConfig) GenerateIDFromConfigFileName(configFile string, musts ...bool) string {
	baseName := filepath.Base(configFile)
	index := strings.LastIndex(baseName, `.`)
	var id string
	if index > 0 {
		id = baseName[0:index]
	}
	if len(musts) == 0 || !musts[0] {
		if len(id) == 0 {
			id = time.Now().Format(`020060102150405`)
		}
	}
	return id
}

func (c *CLIConfig) Close() error {
	for _, cmd := range c.cmds {
		err := c.Kill(cmd)
		if err != nil {
			log.Error(err)
		}
	}
	return nil
}

func (c *CLIConfig) Pid() int {
	return c.pid
}
