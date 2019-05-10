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
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/goforever"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/dbschema"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var Daemon = goforever.Default

func init() {
	Daemon.Name = `nging`
	Daemon.Debug = true
}

func RunDaemon() {
	if !IsInstalled() {
		return
	}
	processHook := func(p *goforever.Process) {
		/*
			if p.Status == goforever.StatusRestarted {
				return
			}
		*/
		processM := &dbschema.ForeverProcess{}
		err := processM.Get(nil, nil, `id`, p.Name)
		if err != nil {
			return
		}
		switch p.Status {
		case goforever.StatusStarted:
			processM.Lastrun = uint(time.Now().Unix())
			processM.Pid = p.Pid
		case goforever.StatusStopped:
			processM.Pid = 0
		case goforever.StatusExited:
			processM.Pid = 0
		case goforever.StatusKilled:
			processM.Pid = 0
		}
		processM.Status = p.Status
		err = processM.Edit(nil, `id`, p.Name)
		if err != nil {
			log.Error(err)
		}
	}
	Daemon.AddHook(goforever.StatusStarted, processHook)
	Daemon.AddHook(goforever.StatusStopped, processHook)
	Daemon.AddHook(goforever.StatusRunning, processHook)
	Daemon.AddHook(goforever.StatusRestarted, processHook)
	Daemon.AddHook(goforever.StatusExited, processHook)
	Daemon.AddHook(goforever.StatusKilled, processHook)
	_ = processHook
	processM := &dbschema.ForeverProcess{}
	_, err := processM.ListByOffset(nil, nil, 0, -1, `disabled`, `N`)
	if err != nil {
		log.Error(err)
		return
	}
	for _, p := range processM.Objects() {
		AddDaemon(p)
	}
	Daemon.Run()
}

func AddDaemon(p *dbschema.ForeverProcess, run ...bool) *goforever.Process {
	name := fmt.Sprint(p.Id)
	procs := goforever.NewProcess(name, p.Command, ParseArgsSlice(p.Args)...)
	procs.Debug = p.Debug == `Y`
	procs.Delay = p.Delay
	procs.Ping = p.Ping
	procs.Respawn = int(p.Respawn)

	procs.Env = ParseEnvSlice(p.Env)
	procs.Dir = p.Workdir

	procs.Errfile = p.Errfile
	procs.Logfile = p.Logfile

	pidFile := filepath.Join(echo.Wd(), `data/pid/daemon`)
	if !com.IsDir(pidFile) {
		err := os.MkdirAll(pidFile, os.ModePerm)
		if err != nil {
			log.Error(err)
		}
	}
	pidFile = filepath.Join(pidFile, fmt.Sprintf(`%d.pid`, p.Id))
	procs.Pidfile = goforever.Pidfile(pidFile)
	Daemon.Add(name, procs, run...)
	return procs
}

func ParseEnvSlice(a string) []string {
	var env []string
	a = strings.TrimSpace(a)
	if len(a) > 0 {
		for _, row := range strings.Split(a, "\n") {
			row = strings.TrimSpace(row)
			if len(row) > 0 {
				env = append(env, row)
			}
		}
	}
	return env
}

func ParseArgsSlice(a string) []string {
	var args []string
	a = strings.TrimSpace(a)
	if len(a) > 0 {
		for _, row := range strings.Split(a, "\n") {
			row = strings.TrimSpace(row)
			if len(row) > 0 {
				args = append(args, row)
			}
		}
	}
	return args
}
