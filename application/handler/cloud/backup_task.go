package cloud

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/checksum"
	"github.com/admpub/log"
	"github.com/admpub/nging/application/handler"
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/msgbox"
	"github.com/admpub/nging/application/library/s3manager/s3client"
	"github.com/admpub/nging/application/model"
	"github.com/syndtr/goleveldb/leveldb"
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
	if err = fullBackupStart(recv); err != nil {
		if err == ErrRunningPleaseWait {
			err = ctx.NewError(code.OperationProcessing, ctx.T(`运行中，请稍候，如果文件很多可能需要会多等一会儿`))
		}
		return err
	}
	if err = monitorBackupStart(recv); err != nil {
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
	if fullBackupIsRunning(m.Id) {
		fullBackupExit = true
	}
	if err = monitorBackupStop(m.Id); err != nil {
		return err
	}
	handler.SendOk(ctx, ctx.T(`操作成功`))
	return ctx.Redirect(handler.URLFor(`/cloud/backup`))
}

func fileFilter(recv *model.CloudBackupExt) (func(string) bool, error) {
	var (
		re  *regexp.Regexp
		err error
	)
	if len(recv.IgnoreRule) > 0 {
		re, err = regexp.Compile(recv.IgnoreRule)
		if err != nil {
			return nil, err
		}
	}
	return func(file string) bool {
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
	}, nil
}

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
		msgbox.Error(`Delete`, file)
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

var ErrRunningPleaseWait = errors.New("running, please wait")

func fullBackupIsRunning(id uint) bool {
	idKey := com.String(id)
	key := `cloud.backup-task.` + idKey
	return echo.Bool(key)
}

var fullBackupExit bool

// 全量备份
func fullBackupStart(recv *model.CloudBackupExt) error {
	idKey := com.String(recv.Id)
	key := `cloud.backup-task.` + idKey
	if echo.Bool(key) {
		return ErrRunningPleaseWait
	}
	echo.Set(key, true)
	defer func() { echo.Delete(key) }()
	cacheDir := filepath.Join(echo.Wd(), `data/cache/backup-db`)
	os.MkdirAll(cacheDir, 0777)
	cacheFile := filepath.Join(cacheDir, idKey)
	db, err := leveldb.OpenFile(cacheFile, nil)
	if err != nil {
		return err
	}
	defer db.Close()
	filter, err := fileFilter(recv)
	if err != nil {
		return err
	}
	sourcePath := recv.SourcePath
	//debug := !config.DefaultConfig.Sys.IsEnv(`prod`)
	recv.Storage.Secret = common.Crypto().Decode(recv.Storage.Secret)
	mgr, err := s3client.New(recv.Storage, config.DefaultConfig.Sys.EditableFileMaxBytes)
	if err != nil {
		return err
	}
	fullBackupExit = false
	err = filepath.Walk(sourcePath, func(ppath string, info os.FileInfo, err error) error {
		if fullBackupExit {
			return echo.ErrExit
		}
		if info.IsDir() {
			return nil
		}
		if !filter(ppath) {
			return filepath.SkipDir
		}
		var oldMd5 string
		dbKey := com.Str2bytes(ppath)
		cv, ce := db.Get(dbKey, nil)
		if ce != nil {
			if ce != leveldb.ErrNotFound {
				return ce
			}
		} else {
			oldMd5 = com.Bytes2str(cv)
		}
		md5, err := checksum.MD5sum(ppath)
		if err != nil {
			return err
		}
		if oldMd5 == md5 {
			return nil
		}

		objectName := path.Join(recv.DestPath, strings.TrimPrefix(ppath, sourcePath))
		fp, err := os.Open(ppath)
		if err != nil {
			log.Error(err)
			return err
		}
		defer func() {
			fp.Close()
			err = db.Put(dbKey, com.Str2bytes(md5), nil)
			if err != nil {
				log.Error(err)
			}
		}()
		return mgr.Put(fp, objectName, info.Size())
	})
	if err == echo.ErrExit {
		log.Info(`强制退出全量备份`)
		err = nil
	}
	fullBackupExit = false
	return err
}
