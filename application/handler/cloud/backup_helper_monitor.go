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
	"time"

	"github.com/admpub/nging/v4/application/library/cloudbackup"
	"github.com/admpub/nging/v4/application/library/common"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/admpub/nging/v4/application/library/msgbox"
	"github.com/admpub/nging/v4/application/library/s3manager/s3client"
	"github.com/admpub/nging/v4/application/model"
	"github.com/webx-top/com"
	"golang.org/x/sync/singleflight"
)

// 通过监控文件变动来进行备份
func monitorBackupStart(recv *model.CloudBackupExt) error {
	if err := monitorBackupStop(recv.Id); err != nil {
		return err
	}
	monitor := com.NewMonitor()
	cloudbackup.BackupTasks.Set(recv.Id, monitor)
	monitor.Debug = !config.DefaultConfig.Sys.IsEnv(`prod`)
	recv.Storage.Secret = common.Crypto().Decode(recv.Storage.Secret)
	mgr, err := s3client.New(recv.Storage, config.DefaultConfig.Sys.EditableFileMaxBytes)
	if err != nil {
		return err
	}
	filter, err := fileFilter(recv)
	if err != nil {
		return err
	}
	var delay time.Duration
	if recv.Delay > 0 {
		delay = time.Duration(recv.Delay) * time.Second
	}
	waitFillCompleted := recv.WaitFillCompleted == `Y`
	var ignoreWaitRegexp *regexp.Regexp
	if waitFillCompleted && len(recv.IgnoreWaitRule) > 0 {
		ignoreWaitRegexp, err = regexp.Compile(recv.IgnoreWaitRule)
		if err != nil {
			return err
		}
	}
	monitor.SetFilters(filter)
	sourcePath, err := filepath.Abs(recv.SourcePath)
	if err != nil {
		return err
	}

	backup := cloudbackup.New(mgr)
	backup.DestPath = recv.DestPath
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
	msgbox.Success(`Cloud-Backup`, `Watch Dir: `+recv.SourcePath)
	err = monitor.AddDir(recv.SourcePath)
	if err != nil {
		return err
	}
	monitor.Watch()
	return nil
}

func monitorBackupStop(id uint) error {
	return cloudbackup.MonitorBackupStop(id)
}
