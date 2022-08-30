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

	"github.com/admpub/events"
	"github.com/admpub/goforever"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	ngingdbschema "github.com/admpub/nging/v4/application/dbschema"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/cron"
	"github.com/nging-plugins/servermanager/application/dbschema"
)

var (
	Daemon            = goforever.Default
	DaemonDefaultHook = DaemonCommonHook
)

func init() {
	Daemon.Name = `nging`
	Daemon.Debug = true
	echo.OnCallback(`nging.config.extend.unregister`, func(data events.Event) error {
		if config.FromFile() == nil {
			return nil
		}
		config.FromFile().UnregisterExtend(data.Context.String(`name`))
		return nil
	})
}

func DaemonCommonHook(p *goforever.Process) {
	if p == Daemon {
		return
	}
	/*
		if p.Status == goforever.StatusRestarted {
			return
		}
	*/
	processM := dbschema.NewNgingForeverProcess(nil)
	err := processM.Get(nil, `id`, p.Name)
	if err != nil {
		log.Errorf(`Not found ForeverProcess: %v (%v)`, p.Name, err)
		return
	}
	switch p.Status {
	case goforever.StatusStarted:
		processM.Lastrun = uint(time.Now().Unix())
		processM.Pid = p.Pid
	case goforever.StatusStopped:
		processM.Pid = 0
	case goforever.StatusExited:
		OnExitedDaemon(processM)
		processM.Pid = 0
	case goforever.StatusKilled:
		processM.Pid = 0
	}
	processM.Status = p.Status
	set := echo.H{
		`status`: processM.Status,
	}
	if p.Error() != nil {
		set[`error`] = p.Error().Error()
	}
	err = processM.UpdateFields(nil, set, `id`, p.Name)
	if err != nil {
		log.Errorf(`Update ForeverProcess: %v (%v)`, p.Name, err)
	}
}

// RestartDaemon 重启所有已登记的进程
func RestartDaemon() {
	RunDaemon()
}

// RunDaemon 运行值守程序
func RunDaemon() {
	if !config.IsInstalled() {
		return
	}
	Daemon.Reset()
	Daemon.SetHook(goforever.StatusStarted, DaemonDefaultHook)
	Daemon.SetHook(goforever.StatusStopped, DaemonDefaultHook)
	Daemon.SetHook(goforever.StatusRunning, DaemonDefaultHook)
	Daemon.SetHook(goforever.StatusRestarted, DaemonDefaultHook)
	Daemon.SetHook(goforever.StatusExited, DaemonDefaultHook)
	Daemon.SetHook(goforever.StatusKilled, DaemonDefaultHook)
	processM := dbschema.NewNgingForeverProcess(nil)
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

func AddDaemon(p *dbschema.NgingForeverProcess, run ...bool) *goforever.Process {
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
	err := com.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		log.Error(err)
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

// OnExitedDaemon 当值守程序达到最大重试次数退出时
func OnExitedDaemon(processM *dbschema.NgingForeverProcess) {
	// 发送邮件通知
	if processM.EnableNotify == 0 {
		return
	}
	user := ngingdbschema.NewNgingUser(nil)
	if processM.Uid > 0 {
		user.Get(nil, `id`, processM.Uid)
	}
	var ccList []string
	if len(processM.NotifyEmail) > 0 {
		ccList = strings.Split(processM.NotifyEmail, "\n")
		for index, email := range ccList {
			email = strings.TrimSpace(email)
			if len(email) == 0 {
				continue
			}
			ccList[index] = email
		}
	}
	if len(user.Email) == 0 {
		if len(ccList) == 0 {
			return
		}
		user.Email = ccList[0]
		user.Username = strings.SplitN(user.Email, `@`, 2)[0]
		if len(ccList) > 1 {
			ccList = ccList[1:]
		}
	}
	title := `[Nging][进程值守警报]进程[` + processM.Name + `]已经异常退出`
	content := `<h1>进程值守警报</h1><p>进程<strong>` + processM.Name + `</strong>已经于<strong>` + time.Now().Format(time.RFC3339) + `</strong>异常退出，请马上处理</p>`
	if err := cron.SendMail(user.Email, user.Username, title, com.Str2bytes(content), ccList...); err != nil {
		log.Errorf(`发送进程值守警报失败：%v`, err)
	}
}
