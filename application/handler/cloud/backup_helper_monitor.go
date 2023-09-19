/*
   Nging is a toolbox for webmasters
   Copyright (C) 2021-present Wenhui Shen <swh@admpub.com>

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

package cloud

import (
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/cloudbackup"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/msgbox"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"golang.org/x/sync/singleflight"
)

// 通过监控文件变动来进行备份
func monitorBackupStart(cfg dbschema.NgingCloudBackup, debug ...bool) error {
	if err := monitorBackupStop(cfg.Id); err != nil {
		return err
	}
	monitor := com.NewMonitor()
	cloudbackup.BackupTasks.Set(cfg.Id, monitor)
	if len(debug) > 0 {
		monitor.Debug = debug[0]
	} else {
		monitor.Debug = !config.FromFile().Sys.IsEnv(`prod`)
	}
	ctx := defaults.NewMockContext()
	mgr, err := cloudbackup.NewStorage(ctx, cfg)
	if err != nil {
		return err
	}
	if err := mgr.Connect(); err != nil {
		return err
	}
	filter, err := fileFilter(&cfg)
	if err != nil {
		return err
	}
	var delay time.Duration
	if cfg.Delay > 0 {
		delay = time.Duration(cfg.Delay) * time.Second
	}
	waitFillCompleted := cfg.WaitFillCompleted == `Y`
	var ignoreWaitRegexp *regexp.Regexp
	if waitFillCompleted && len(cfg.IgnoreWaitRule) > 0 {
		ignoreWaitRegexp, err = regexp.Compile(cfg.IgnoreWaitRule)
		if err != nil {
			return err
		}
	}
	monitor.SetFilters(filter)
	sourcePath, err := filepath.Abs(cfg.SourcePath)
	if err != nil {
		return err
	}
	sourcePath, err = filepath.EvalSymlinks(sourcePath)
	if err != nil {
		return err
	}
	if !strings.HasSuffix(sourcePath, echo.FilePathSeparator) {
		sourcePath += echo.FilePathSeparator
	}

	backup := cloudbackup.New(mgr, cfg)
	backup.DestPath = cfg.DestPath
	backup.SourcePath = sourcePath
	backup.Filter = filter
	backup.WaitFillCompleted = waitFillCompleted
	backup.IgnoreWaitRegexp = ignoreWaitRegexp

	var sg singleflight.Group
	monitor.Create = func(file string) {
		if monitor.Debug {
			msgbox.Success(`Create`, file)
		}
		sg.Do(file, func() (interface{}, error) {
			if delay > 0 {
				time.Sleep(delay)
			}
			backup.OnCreate(file)
			return nil, nil
		})
	}
	monitor.Delete = func(file string) {
		if monitor.Debug {
			msgbox.Error(`Delete`, file)
		}
		backup.OnDelete(file)
	}
	monitor.Modify = func(file string) {
		if monitor.Debug {
			msgbox.Info(`Modify`, file)
		}
		sg.Do(file, func() (interface{}, error) {
			if delay > 0 {
				time.Sleep(delay)
			}
			backup.OnModify(file)
			return nil, nil
		})
	}
	monitor.Rename = func(file string) {
		if monitor.Debug {
			msgbox.Warn(`Rename`, file)
		}
		backup.OnRename(file)
	}
	msgbox.Success(`Cloud-Backup`, `Watch Dir: `+backup.SourcePath)
	err = monitor.AddDir(backup.SourcePath)
	if err != nil {
		return err
	}
	monitor.Watch()
	return nil
}

func monitorBackupStop(id uint) error {
	return cloudbackup.MonitorBackupStop(id)
}
