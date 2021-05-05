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
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/flock"
	"github.com/admpub/nging/application/library/msgbox"
	"github.com/admpub/nging/application/library/s3manager/s3client"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
)

// 通过监控文件变动来进行备份
func monitorBackupStart(recv *model.CloudBackupExt) error {
	if err := monitorBackupStop(recv.Id); err != nil {
		return err
	}
	monitor := com.NewMonitor()
	backupTasks.Set(recv.Id, monitor)
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
	monitor.SetFilters(filter)
	sourcePath, err := filepath.Abs(recv.SourcePath)
	if err != nil {
		return err
	}
	monitor.Create = func(file string) {
		if monitor.Debug {
			msgbox.Success(`Create`, file)
		}
		fp, err := os.Open(file)
		if err != nil {
			log.Error(file + `: ` + err.Error())
			return
		}
		defer fp.Close()
		if err = flock.LockBlock(fp); err != nil {
			log.Error(file + `: ` + err.Error())
			return
		}
		defer flock.Unlock(fp)
		fi, err := fp.Stat()
		if err != nil {
			log.Error(file + `: ` + err.Error())
			return
		}
		if fi.IsDir() {
			err = filepath.Walk(file, func(ppath string, info os.FileInfo, werr error) error {
				if werr != nil {
					return werr
				}
				if info.IsDir() || !filter(ppath) {
					return nil
				}
				objectName := path.Join(recv.DestPath, strings.TrimPrefix(ppath, sourcePath))
				fp, err := os.Open(ppath)
				if err != nil {
					log.Error(ppath + `: ` + err.Error())
					return err
				}
				defer fp.Close()
				if err = flock.LockBlock(fp); err != nil {
					log.Error(ppath + `: ` + err.Error())
					return err
				}
				defer flock.Unlock(fp)
				return mgr.Put(fp, objectName, info.Size())
			})
		} else {
			objectName := path.Join(recv.DestPath, strings.TrimPrefix(file, sourcePath))
			err = mgr.Put(fp, objectName, fi.Size())
		}
		if err != nil {
			log.Error(err)
		}
	}
	monitor.Delete = func(file string) {
		if monitor.Debug {
			msgbox.Error(`Delete`, file)
		}
		objectName := path.Join(recv.DestPath, strings.TrimPrefix(file, sourcePath))
		err = mgr.RemoveDir(objectName)
		if err != nil {
			log.Error(file + `: ` + err.Error())
		}
		err = mgr.Remove(objectName)
		if err != nil {
			log.Error(file + `: ` + err.Error())
		}
	}
	monitor.Modify = func(file string) {
		if monitor.Debug {
			msgbox.Info(`Modify`, file)
		}
		objectName := path.Join(recv.DestPath, strings.TrimPrefix(file, sourcePath))
		fp, err := os.Open(file)
		if err != nil {
			log.Error(file + `: ` + err.Error())
			return
		}
		defer fp.Close()
		if err = flock.LockBlock(fp); err != nil {
			log.Error(file + `: ` + err.Error())
			return
		}
		defer flock.Unlock(fp)
		fi, err := fp.Stat()
		if err != nil {
			log.Error(file + `: ` + err.Error())
			return
		}
		err = mgr.Put(fp, objectName, fi.Size())
		if err != nil {
			log.Error(file + `: ` + err.Error())
		}
	}
	monitor.Rename = func(file string) {
		if monitor.Debug {
			msgbox.Warn(`Rename`, file)
		}
		objectName := path.Join(recv.DestPath, strings.TrimPrefix(file, sourcePath))
		err = mgr.RemoveDir(objectName)
		if err != nil {
			log.Error(file + `: ` + err.Error())
		}
		err = mgr.Remove(objectName)
		if err != nil {
			log.Error(file + `: ` + err.Error())
		}
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
	if monitor, ok := backupTasks.Get(id).(*com.MonitorEvent); ok {
		monitor.Close()
		backupTasks.Delete(id)
		msgbox.Success(`Cloud-Backup`, `Close: `+com.String(id))
	}
	return nil
}
