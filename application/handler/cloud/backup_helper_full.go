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
	"github.com/admpub/nging/application/library/common"
	"github.com/admpub/nging/application/library/config"
	"github.com/admpub/nging/application/library/s3manager/s3client"
	"github.com/admpub/nging/application/model"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
)

var (
	// ErrRunningPleaseWait 正在运行中
	ErrRunningPleaseWait = errors.New("running, please wait")
	fullBackupExit       bool
)

func fullBackupIsRunning(id uint) bool {
	idKey := com.String(id)
	key := `cloud.backup-task.` + idKey
	return echo.Bool(key)
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

// 全量备份
func fullBackupStart(recv *model.CloudBackupExt) error {
	idKey := com.String(recv.Id)
	key := `cloud.backup-task.` + idKey
	if echo.Bool(key) {
		return ErrRunningPleaseWait
	}
	echo.Set(key, true)
	sourcePath := recv.SourcePath
	debug := !config.DefaultConfig.Sys.IsEnv(`prod`)
	recv.Storage.Secret = common.Crypto().Decode(recv.Storage.Secret)
	filter, err := fileFilter(recv)
	if err != nil {
		return err
	}
	cacheDir := filepath.Join(echo.Wd(), `data/cache/backup-db`)
	if err := com.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return err
	}
	cacheFile := filepath.Join(cacheDir, idKey)
	mgr, err := s3client.New(recv.Storage, config.DefaultConfig.Sys.EditableFileMaxBytes)
	if err != nil {
		return err
	}
	go func() {
		defer func() {
			echo.Delete(key)
		}()
		ctx := defaults.NewMockContext()
		recv.SetContext(ctx)
		db, err := leveldb.OpenFile(cacheFile, nil)
		if err != nil {
			recv.SetFields(nil, echo.H{
				`result`: err.Error(),
				`status`: `failure`,
			}, `id`, recv.Id)
			return
		}
		defer db.Close()
		fullBackupExit = false
		err = filepath.Walk(sourcePath, func(ppath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
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
				if debug {
					log.Info(ppath, `: 文件备份过并且没有改变【跳过】`)
				}
				return nil
			}
			if debug {
				if len(oldMd5) > 0 {
					log.Info(ppath, `: 文件备份过并且有更改【更新】`)
				} else {
					log.Info(ppath, `: 文件未曾备份过【添加】`)
				}
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
		if err != nil {
			if err == echo.ErrExit {
				log.Info(`强制退出全量备份`)
			} else {
				recv.SetFields(nil, echo.H{
					`result`: err.Error(),
					`status`: `failure`,
				}, `id`, recv.Id)
			}
		} else {
			recv.SetFields(nil, echo.H{
				`result`: ctx.T(`全量备份完成`),
				`status`: `idle`,
			}, `id`, recv.Id)
		}
		fullBackupExit = false
	}()
	return nil
}
