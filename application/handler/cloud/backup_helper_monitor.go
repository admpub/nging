package cloud

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
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
			log.Error(err)
			return
		}
		defer fp.Close()
		fi, err := fp.Stat()
		if err != nil {
			log.Error(err)
			return
		}
		if fi.IsDir() {
			err = filepath.Walk(file, func(ppath string, info os.FileInfo, err error) error {
				if info.IsDir() || !filter(ppath) {
					return nil
				}
				objectName := path.Join(recv.DestPath, strings.TrimPrefix(ppath, sourcePath))
				fp, err := os.Open(ppath)
				if err != nil {
					log.Error(err)
					return err
				}
				defer fp.Close()
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
			log.Error(err)
		}
		err = mgr.Remove(objectName)
		if err != nil {
			log.Error(err)
		}
	}
	monitor.Modify = func(file string) {
		if monitor.Debug {
			msgbox.Info(`Modify`, file)
		}
		objectName := path.Join(recv.DestPath, strings.TrimPrefix(file, sourcePath))
		fp, err := os.Open(file)
		if err != nil {
			log.Error(err)
			return
		}
		defer fp.Close()
		fi, err := fp.Stat()
		if err != nil {
			log.Error(err)
			return
		}
		err = mgr.Put(fp, objectName, fi.Size())
		if err != nil {
			log.Error(err)
		}
	}
	monitor.Rename = func(file string) {
		if monitor.Debug {
			msgbox.Warn(`Rename`, file)
		}
		objectName := path.Join(recv.DestPath, strings.TrimPrefix(file, sourcePath))
		err = mgr.RemoveDir(objectName)
		if err != nil {
			log.Error(err)
		}
		err = mgr.Remove(objectName)
		if err != nil {
			log.Error(err)
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
