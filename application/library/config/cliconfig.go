/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package config

import (
	"errors"
	"flag"
	"io"
	"os"
	"os/exec"

	"github.com/admpub/log"
	"github.com/webx-top/com"
)

type CLIConfig struct {
	Port int
	Conf string
	Type string //启动类型: server/manager
	cmds map[string]*exec.Cmd
}

func (c *CLIConfig) InitFlag() {
	flag.IntVar(&c.Port, `p`, 9999, `port`)
	flag.StringVar(&c.Conf, `c`, `config/config.yaml`, `config`)
	flag.StringVar(&c.Type, `t`, `manager`, `operation type`)
}

func (c *CLIConfig) CaddyStopHistory() (err error) {
	return com.CloseProcessFromPidFile(DefaultConfig.Caddy.PidFile)
}

func (c *CLIConfig) CaddyStart(writer ...io.Writer) (err error) {
	err = c.CaddyStopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	params := []string{os.Args[0], `-c`, c.Conf, `-t`, `webserver`}
	c.cmds["caddy"] = com.RunCmdWithWriter(params, writer...)
	return
}

func (c *CLIConfig) CaddyStop() error {
	if c.cmds == nil {
		return nil
	}
	cmd, ok := c.cmds["caddy"]
	if !ok || cmd == nil {
		return nil
	}
	if cmd.ProcessState != nil {
		return nil
	}
	return cmd.Process.Kill()
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

func (c *CLIConfig) CaddyRestart(writer ...io.Writer) error {
	err := c.CaddyStop()
	if err != nil {
		return err
	}
	return c.CaddyStart(writer...)
}

func (c *CLIConfig) FTPStopHistory() (err error) {
	return com.CloseProcessFromPidFile(DefaultConfig.FTP.PidFile)
}

func (c *CLIConfig) FTPStart(writer ...io.Writer) (err error) {
	err = c.FTPStopHistory()
	if err != nil {
		log.Error(err.Error())
	}
	params := []string{os.Args[0], `-c`, c.Conf, `-t`, `ftpserver`}
	c.cmds["ftp"] = com.RunCmdWithWriter(params, writer...)
	return
}

func (c *CLIConfig) FTPStop() error {
	if c.cmds == nil {
		return nil
	}
	cmd, ok := c.cmds["ftp"]
	if !ok || cmd == nil {
		return nil
	}
	if cmd.ProcessState != nil {
		return nil
	}
	return cmd.Process.Kill()
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

func (c *CLIConfig) FTPRestart(writer ...io.Writer) error {
	err := c.FTPStop()
	if err != nil {
		return err
	}
	return c.FTPStart(writer...)
}

func (c *CLIConfig) Reload() error {
	if c.IsRunning(`caddy`) {
		c.CaddyRestart()
	}
	if c.IsRunning(`ftp`) {
		c.FTPRestart()
	}
	return nil
}
