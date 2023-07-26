/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package service

import (
	"errors"
	"fmt"
	stdLog "log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/admpub/log"
	"github.com/admpub/service"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

func ValidServiceAction(action string) error {
	for _, act := range service.ControlAction {
		if act == action {
			return nil
		}
	}
	return fmt.Errorf("available actions: %q", service.ControlAction)
}

func slog() *log.Logger {
	return log.GetLogger(`service`)
}

// New 以服务的方式启动nging
// 服务支持的操作有：
// nging service install  	-- 安装服务
// nging service uninstall  -- 卸载服务
// nging service start 		-- 启动服务
// nging service stop 		-- 停止服务
// nging service restart 	-- 重启服务
func New(cfg *Config, action string) error {
	p := NewProgram(cfg)
	s, err := service.New(p, &p.Config.Config)
	if err != nil {
		return err
	}
	p.service = s

	// Service
	if action != `run` {
		if err := ValidServiceAction(action); err != nil {
			return err
		}
		err = service.Control(s, action)
		if err != nil {
			slog().Errorf(`%s: %s`, action, err.Error())
		} else {
			slog().Okayf(`%s: success`, action)
		}
		return err
	}
	return s.Run()
}

func getPidFiles() []string {
	pidFile := []string{}
	pidFilePath := filepath.Join(com.SelfDir(), `data/pid`)
	err := filepath.Walk(pidFilePath, func(pidPath string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		if filepath.Ext(pidPath) == `.pid` {
			pidFile = append(pidFile, pidPath)
		}
		return nil
	})
	if err != nil {
		stdLog.Println(err)
	}
	return pidFile
}

func NewProgram(cfg *Config) *program {
	pidFile := filepath.Join(com.SelfDir(), `data/pid`)
	err := com.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		stdLog.Println(err)
	}
	pidFile = filepath.Join(pidFile, `nging.pid`)
	p := &program{
		Config:  cfg,
		pidFile: pidFile,
	}
	p.Config.Config.Arguments = append([]string{`service`, `run`}, p.Args...)
	p.Config.Config.WorkingDirectory = p.Dir
	return p
}

type program struct {
	*Config
	service  service.Service
	cmd      *exec.Cmd
	fullExec string
	pidFile  string
}

func (p *program) Start(s service.Service) (err error) {
	if service.Interactive() {
		p.logger.Info("Running in terminal.")
	} else {
		p.logger.Info("Running under service manager.")
	}
	if filepath.Base(p.Exec) == p.Exec {
		p.fullExec, err = exec.LookPath(p.Exec)
		if err != nil {
			return fmt.Errorf("failed to find executable %q: %v", p.Exec, err)
		}
	} else {
		p.fullExec = p.Exec
	}
	p.createCmd()

	go p.run()
	return nil
}

func (p *program) createCmd() {
	p.cmd = exec.Command(p.fullExec, p.Args...)
	p.cmd.Dir = p.Dir
	p.cmd.Env = param.StringSlice(append(os.Environ(), p.Env...)).Unique().String()
	if p.Stderr != nil {
		p.cmd.Stderr = p.Stderr
	}
	if p.Stdout != nil {
		p.cmd.Stdout = p.Stdout
	}
	p.logger.Infof("Running cmd: %s %#v", p.fullExec, p.Args)
	p.logger.Infof("Workdir: %s", p.cmd.Dir)
	//p.logger.Infof("Env var: %s", com.Dump(p.cmd.Env, false))
}

func (p *program) Stop(s service.Service) error {
	p.killCmd()
	p.logger.Infof("Stopping %s", p.DisplayName)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p *program) killCmd() {
	err := com.CloseProcessFromCmd(p.cmd)
	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		p.logger.Error(err)
	}
	err = com.CloseProcessFromPidFile(p.pidFile)
	if err != nil && !errors.Is(err, os.ErrProcessDone) {
		p.logger.Error(p.pidFile+`:`, err)
	}
	for _, pidFile := range getPidFiles() {
		err = com.CloseProcessFromPidFile(pidFile)
		if err != nil && !errors.Is(err, os.ErrProcessDone) {
			p.logger.Error(pidFile+`:`, err)
		}
	}
}

func (p *program) close() {
	if service.Interactive() {
		p.Stop(p.service)
	} else {
		p.service.Stop()
		p.killCmd()
	}
	if p.Config.OnExited != nil {
		err := p.Config.OnExited()
		if err != nil {
			p.logger.Error(err)
		}
	}
}

func (p *program) run() {
	p.logger.Infof("Starting %s", p.DisplayName)
	//如果调用的程序停止了，则本服务同时也停止
	//defer p.close()
	err := p.cmd.Start()
	if err == nil {
		stdLog.Println("APP PID:", p.cmd.Process.Pid)
		os.WriteFile(p.pidFile, []byte(strconv.Itoa(p.cmd.Process.Pid)), os.ModePerm)
		err = p.cmd.Wait()
	}
	if err != nil {
		p.logger.Error("Error running:", err)
	}
}
