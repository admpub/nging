package cloud

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/admpub/checksum"
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/cloudbackup"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/model"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

var (
	// ErrRunningPleaseWait 正在运行中
	ErrRunningPleaseWait = errors.New("running, please wait")
	fullBackupExit       atomic.Bool
)

func fullBackupIsRunning(id uint) bool {
	idKey := com.String(id)
	key := `cloud.backup-task.` + idKey
	return echo.Bool(key)
}

func fileFilter(cfg *dbschema.NgingCloudBackup) (func(string) bool, error) {
	var (
		re  *regexp.Regexp
		err error
	)
	if len(cfg.IgnoreRule) > 0 {
		re, err = regexp.Compile(cfg.IgnoreRule)
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

// 全量备份
func fullBackupStart(cfg dbschema.NgingCloudBackup) error {
	idKey := com.String(cfg.Id)
	key := `cloud.backup-task.` + idKey
	if echo.Bool(key) {
		return ErrRunningPleaseWait
	}
	echo.Set(key, true)
	sourcePath, err := filepath.Abs(cfg.SourcePath)
	if err != nil {
		return err
	}
	sourcePath, err = filepath.EvalSymlinks(sourcePath)
	if err != nil {
		return err
	}
	debug := !config.FromFile().Sys.IsEnv(`prod`)
	filter, err := fileFilter(&cfg)
	if err != nil {
		return err
	}
	cacheDir := filepath.Join(echo.Wd(), `data/cache/backup-db`)
	if err := com.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return err
	}
	cacheFile := filepath.Join(cacheDir, idKey)
	ctx := defaults.NewMockContext()
	mgr, err := cloudbackup.NewStorage(ctx, cfg)
	if err != nil {
		return err
	}
	if err := mgr.Connect(); err != nil {
		return err
	}
	go func() {
		ctx := defaults.NewMockContext()
		var err error
		defer func() {
			echo.Delete(key)
		}()
		recv := cfg
		recv.SetContext(ctx)
		var db *leveldb.DB
		db, err = leveldb.OpenFile(cacheFile, nil)
		if err != nil {
			recv.UpdateFields(nil, echo.H{
				`result`: err.Error(),
				`status`: `failure`,
			}, `id`, recv.Id)
			return
		}
		defer db.Close()
		fullBackupExit.Store(false)
		err = filepath.Walk(sourcePath, func(ppath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if fullBackupExit.Load() {
				return echo.ErrExit
			}
			if info.IsDir() {
				return nil
			}
			if !filter(ppath) {
				return filepath.SkipDir
			}
			var oldMd5 string
			var md5 string
			var cv []byte
			var operation string
			dbKey := com.Str2bytes(ppath)
			cv, err = db.Get(dbKey, nil)
			var oldValParts []string
			if err != nil {
				if err != leveldb.ErrNotFound {
					return err
				}
				operation = model.CloudBackupOperationCreate
				oldValParts = []string{
					``,                                    // md5
					``,                                    // taskStartTime
					``,                                    // taskEndTime
					param.AsString(info.ModTime().Unix()), // fileModeTime
					param.AsString(info.Size()),           // fileSize
				}
			} else {
				oldValParts = strings.Split(string(cv), `||`)
				oldMd5 = oldValParts[0]
				if len(oldValParts) < 5 {
					temp := make([]string, 5)
					copy(temp, oldValParts)
					oldValParts = temp
				}
				oldValParts[3] = param.AsString(info.ModTime().Unix())
				oldValParts[4] = param.AsString(info.Size())
				operation = model.CloudBackupOperationUpdate
			}
			if len(oldMd5) > 0 {
				md5, err = checksum.MD5sum(ppath)
				if err != nil {
					return err
				}
				if oldMd5 == md5 {
					if debug {
						log.Info(ppath, `: 文件备份过并且没有改变【跳过】`)
					}
					return nil
				}
			}
			if debug {
				if len(oldMd5) > 0 {
					log.Info(ppath, `: 文件备份过并且有更改【更新】`)
				} else {
					log.Info(ppath, `: 文件未曾备份过【添加】`)
				}
			}

			objectName := path.Join(recv.DestPath, strings.TrimPrefix(ppath, sourcePath))
			startTime := time.Now()
			defer func() {
				cloudbackup.RecordLog(ctx, err, &cfg, ppath, objectName, operation, startTime, uint64(info.Size()), model.CloudBackupTypeFull)
			}()
			var fp *os.File
			fp, err = os.Open(ppath)
			if err != nil {
				log.Error(err)
				return err
			}
			defer func() {
				fp.Close()
				if err != nil {
					return
				}
				oldValParts[0] = md5
				err = db.Put(dbKey, com.Str2bytes(strings.Join(oldValParts, `||`)), nil)
				if err != nil {
					log.Error(err)
				}
			}()
			err = cloudbackup.RetryablePut(ctx, mgr, fp, objectName, info.Size())
			return err
		})
		if err != nil {
			if err == echo.ErrExit {
				log.Info(`强制退出全量备份`)
			} else {
				recv.UpdateFields(nil, echo.H{
					`result`: err.Error(),
					`status`: `failure`,
				}, `id`, recv.Id)
			}
		} else {
			recv.UpdateFields(nil, echo.H{
				`result`: ctx.T(`全量备份完成`),
				`status`: `idle`,
			}, `id`, recv.Id)
		}
		fullBackupExit.Store(false)
	}()
	return nil
}
