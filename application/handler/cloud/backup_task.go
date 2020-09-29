package cloud

import (
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/msgbox"
	"github.com/admpub/nging/application/library/s3manager/s3client"
	"github.com/admpub/nging/application/model"
	"github.com/webx-top/com"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

var backupTasks = param.NewMap()

func BackupStart(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	recv := &model.CloudBackupExt{}
	err := m.NewParam().SetArgs(db.Cond{`id`: id}).SetRecv(recv).One()
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, ctx.T(`数据不存在`))
		}
		return err
	}
	if len(recv.Storage.Endpoint) == 0 {
		return ctx.NewError(code.InvalidParameter, ctx.T(`Endpoint无效`))
	}
	if err = backupStart(recv); err != nil {
		return err
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(handler.URLFor(`/cloud/backup`))
}

func BackupStop(ctx echo.Context) error {
	id := ctx.Formx(`id`).Uint()
	m := model.NewCloudBackup(ctx)
	err := m.Get(nil, db.Cond{`id`: id})
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = ctx.NewError(code.DataNotFound, ctx.T(`数据不存在`))
		}
		return err
	}
	if err = backupStop(m.Id); err != nil {
		return err
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(handler.URLFor(`/cloud/backup`))
}

func backupStart(recv *model.CloudBackupExt) error {
	if err := backupStop(recv.Id); err != nil {
		return err
	}
	monitor := com.NewMonitor()
	monitor.Debug = !config.DefaultConfig.Sys.IsEnv(`prod`)
	recv.Storage.Secret = common.Crypto().Decode(recv.Storage.Secret)
	mgr, err := s3client.New(recv.Storage, config.DefaultConfig.Sys.EditableFileMaxBytes)
	if err != nil {
		return err
	}
	var re *regexp.Regexp
	if len(recv.IgnoreRule) > 0 {
		re, err = regexp.Compile(recv.IgnoreRule)
		if err != nil {
			return err
		}
	}
	filter := func(file string) bool {
		switch filepath.Ext(file) {
		case ".swp":
			return false
		case ".tmp", ".TMP":
			return false
		default:
			if strings.Contains(file, echo.FilePathSeparator+`.`) { // 忽略所有以点号开头的文件
				return false
			}
			if re != nil {
				if re.MatchString(file) {
					return false
				}
			}
			return true
		}
	}
	monitor.SetFilters(filter)
	sourcePath, err := filepath.Abs(recv.SourcePath)
	if err != nil {
		return err
	}
	monitor.Create = func(file string) {
		msgbox.Success(`Create`, file)
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
				objectName := path.Join(recv.DestPath, strings.TrimPrefix(file, sourcePath))
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
		msgbox.Error(`Delete`, file)
		objectName := path.Join(recv.DestPath, strings.TrimPrefix(file, sourcePath))
		if com.IsDir(file) {
			err = mgr.RemoveDir(objectName)
		} else {
			err = mgr.Remove(objectName)
		}
		if err != nil {
			log.Error(err)
		}
	}
	monitor.Modify = func(file string) {
		msgbox.Info(`Modify`, file)
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
		msgbox.Warn(`Rename`, file)
		objectName := path.Join(recv.DestPath, strings.TrimPrefix(file, sourcePath))
		if com.IsDir(file) {
			err = mgr.RemoveDir(objectName)
		} else {
			err = mgr.Remove(objectName)
		}
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

func backupStop(id uint) error {
	if monitor, ok := backupTasks.Get(id).(*com.MonitorEvent); ok {
		monitor.Close()
		backupTasks.Delete(id)
		msgbox.Success(`Cloud-Backup`, `Close: `+com.String(id))
	}
	return nil
}
