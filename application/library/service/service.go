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
package service

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/service"
	"github.com/webx-top/com"
)

func ValidServiceAction(action string) error {
	for _, act := range service.ControlAction {
		if act == action {
			return nil
		}
	}
	return fmt.Errorf("Available actions: %q", service.ControlAction)
}

// New 以服务的方式启动nging
// 服务支持的操作有：
// nging install  	-- 安装服务
// nging uninstall  -- 卸载服务
// nging start 		-- 启动服务
// nging stop 		-- 停止服务
// nging restart 	-- 重启服务
func New(cfg *Config, action string) error {
	p := NewProgram(cfg)
	p.Config.Arguments = append([]string{`run`}, p.Args...)

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
		return service.Control(s, action)
	}
	return s.Run()
}

func getPidFiles() []string {
	pidFile := []string{}
	ftpPid := config.DefaultConfig.FTP.PidFile
	caddyPid := config.DefaultConfig.Caddy.PidFile
	if len(ftpPid) == 0 {
		ftpPid = `ftp.pid`
	}
	if len(caddyPid) == 0 {
		caddyPid = `caddy.pid`
	}
	if runtime.GOOS == `windows` {
		if !strings.Contains(ftpPid, `:`) {
			ftpPid = filepath.Join(com.SelfDir(), ftpPid)
		}
		if !strings.Contains(caddyPid, `:`) {
			caddyPid = filepath.Join(com.SelfDir(), caddyPid)
		}
	} else {
		if !strings.HasPrefix(ftpPid, `/`) {
			ftpPid = filepath.Join(com.SelfDir(), ftpPid)
		}
		if !strings.HasPrefix(caddyPid, `/`) {
			caddyPid = filepath.Join(com.SelfDir(), caddyPid)
		}
	}
	pidFile = append(pidFile, caddyPid)
	pidFile = append(pidFile, ftpPid)
	return pidFile
}

func NewProgram(cfg *Config) *program {
	return &program{
		Config:  cfg,
		pidFile: filepath.Join(com.SelfDir(), `nging.pid`),
	}
}

type program struct {
	*Config
	service  service.Service
	cmd      *exec.Cmd
	stopped  bool
	exited   chan struct{}
	fullExec string
	pidFile  string
}

func (p *program) Start(s service.Service) (err error) {
	if service.Interactive() {
		log.Println("Running in terminal.")
	} else {
		log.Println("Running under service manager.")
	}
	// Look for exec.
	// Verify home directory.
	p.fullExec, err = exec.LookPath(p.Exec)
	if err != nil {
		return fmt.Errorf("Failed to find executable %q: %v", p.Exec, err)
	}
	p.stopped = false
	p.exited = make(chan struct{})
	p.createCmd()

	go p.run()
	return nil
}

func (p *program) createCmd() {
	p.cmd = exec.Command(p.fullExec, p.Args...)
	p.cmd.Dir = p.Dir
	p.cmd.Env = append(os.Environ(), p.Env...)
	log.Printf("Running cmd: %s %#v\n", p.fullExec, p.Args)
}

func (p *program) Stop(s service.Service) error {
	p.stopped = true
	p.killCmd()
	log.Println("Stopping", p.DisplayName)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p *program) killCmd() {
	if p.cmd != nil && p.cmd.ProcessState != nil {
		if p.cmd.ProcessState.Exited() == false && p.cmd.Process != nil {
			p.cmd.Process.Kill()
		}
	}
	err := com.CloseProcessFromPidFile(p.pidFile)
	if err != nil {
		log.Println(p.pidFile+`:`, err)
	}
	for _, pidFile := range getPidFiles() {
		err = com.CloseProcessFromPidFile(pidFile)
		if err != nil {
			log.Println(pidFile+`:`, err)
		}
	}
}

func (p *program) close() {
	if service.Interactive() {
		p.Stop(p.service)
	} else {
		p.service.Stop()
	}
}

func FileWriter(file string) (io.WriteCloser, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	return f, err
}

func (p *program) run() {
	log.Println("Starting", p.DisplayName)

	//如果调用的程序停止了，则本服务同时也停止
	defer p.close()

	if p.Stderr != nil {
		p.cmd.Stderr = p.Stderr
	}
	if p.Stdout != nil {
		p.cmd.Stdout = p.Stdout
	}

	go func() {
		for i := 0; i < 10 && !p.stopped; i++ {
			err := p.cmd.Start()
			if err == nil {
				log.Println("APP PID:", p.cmd.Process.Pid)
				ioutil.WriteFile(p.pidFile, []byte(strconv.Itoa(p.cmd.Process.Pid)), os.ModePerm)
				err = p.cmd.Wait()
			}
			if err != nil {
				log.Println("Error running:", err)
				p.killCmd()
				p.createCmd()
			} else {
				i = -1
			}
		}
		p.exited <- struct{}{}
	}()
	<-p.exited
	p.killCmd()
}
