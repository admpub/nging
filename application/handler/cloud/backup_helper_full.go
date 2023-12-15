package cloud

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/admpub/checksum"
	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/cloudbackup"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/nging/v5/application/library/notice"
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
	fileSources          = map[string]*FileSource{}
)

type FileSource struct {
	Name        string
	Description string
	fileSystem  func() http.FileSystem
}

func RegisterFileSource(name string, description string, fileSystem func() http.FileSystem) {
	fileSources[name] = &FileSource{Name: name, Description: description, fileSystem: fileSystem}
}

func GetFileSources() []FileSource {
	results := make([]FileSource, 0, len(fileSources))
	names := make([]string, 0, len(fileSources))
	for name := range fileSources {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		names = append(names, name)
		fss := fileSources[name]
		r := *fss
		r.Name = name + `:`
		results = append(results, r)
	}
	return results
}

func GetFileSource(name string) (fs *FileSource) {
	fs = fileSources[name]
	return
}

func fullBackupIsRunning(id uint) bool {
	idKey := com.String(id)
	key := `cloud.backup-task.` + idKey
	return echo.Bool(key)
}

func fileFilter(rootPath string, cfg *dbschema.NgingCloudBackup) (func(file string) bool, error) {
	var (
		ignoreRE *regexp.Regexp
		matchRE  *regexp.Regexp
		err      error
	)
	if len(cfg.IgnoreRule) > 0 {
		ignoreRE, err = regexp.Compile(cfg.IgnoreRule)
		if err != nil {
			return nil, err
		}
	}
	if len(cfg.MatchRule) > 0 {
		matchRE, err = regexp.Compile(cfg.MatchRule)
		if err != nil {
			return nil, err
		}
	}
	return func(file string) bool {
		relPath := strings.TrimPrefix(file, rootPath)
		if len(relPath) == 0 || relPath == `/` || relPath == `\` {
			return true
		}
		switch filepath.Ext(relPath) {
		case ".swp":
			return false
		case ".tmp", ".TMP":
			return false
		}
		if strings.Contains(relPath, echo.FilePathSeparator+`.`) { // 忽略所有以点号开头的文件
			return false
		}
		if matchRE != nil {
			return matchRE.MatchString(relPath)
		}
		if ignoreRE != nil {
			return !ignoreRE.MatchString(relPath)
		}
		return true
	}, nil
}

// 全量备份
func fullBackupStart(cfg dbschema.NgingCloudBackup, username string, msgType string) error {
	idKey := com.String(cfg.Id)
	key := `cloud.backup-task.` + idKey
	if echo.Bool(key) {
		return ErrRunningPleaseWait
	}
	echo.Set(key, true)
	parts := strings.SplitN(cfg.SourcePath, `:`, 2)
	var fileSystem http.FileSystem
	var sourcePath string
	var err error
	if len(parts) == 2 {
		if fss := GetFileSource(parts[0]); fss != nil {
			fileSystem = fss.fileSystem()
			sourcePath = parts[1]
		}
	}
	if fileSystem == nil {
		sourcePath, err = filepath.Abs(cfg.SourcePath)
		if err != nil {
			return err
		}
		sourcePath, err = filepath.EvalSymlinks(sourcePath)
		if err != nil {
			return err
		}
	}
	debug := !config.FromFile().Sys.IsEnv(`prod`)
	filter, err := fileFilter(sourcePath, &cfg)
	if err != nil {
		return err
	}
	ctx := defaults.NewMockContext()
	mgr, err := cloudbackup.NewStorage(ctx, cfg)
	if err != nil {
		return err
	}
	if err := mgr.Connect(); err != nil {
		return err
	}
	noticeTitle := ctx.T(`全量备份`)
	go func() {
		ctx := defaults.NewMockContext()
		var err error
		defer func() {
			mgr.Close()
			echo.Delete(key)
		}()
		recv := cfg
		recv.SetContext(ctx)
		var db *leveldb.DB
		db, err = cloudbackup.LevelDB().OpenDB(cfg.Id)
		if err != nil {
			recv.UpdateFields(nil, echo.H{
				`result`: err.Error(),
				`status`: `failure`,
			}, `id`, recv.Id)
			return
		}
		fullBackupExit.Store(false)
		putFile := func(ppath string, info os.FileInfo) error {
			if fullBackupExit.Load() {
				return echo.ErrExit
			}
			if !filter(ppath) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
			if info.IsDir() {
				return nil
			}
			var md5 string
			var (
				oldMd5                 string
				fileModifyTs, fileSize int64
			)
			var cv []byte
			var operation string
			dbKey := com.Str2bytes(ppath)
			cv, err = db.Get(dbKey, nil)
			if err != nil {
				if err != leveldb.ErrNotFound {
					return err
				}
				operation = model.CloudBackupOperationCreate
			} else {
				oldMd5, _, _, fileModifyTs, fileSize = cloudbackup.ParseDBValue(cv)
				if info.Size() == fileSize && fileModifyTs == info.ModTime().Unix() {
					if debug {
						log.Info(ppath, `: 文件备份过并且没有改变【跳过】`)
					}
					return nil
				}
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
			var seekReader io.ReadSeekCloser
			if fileSystem == nil {
				seekReader, err = os.Open(ppath)
				if err != nil {
					log.Error(err)
					return err
				}
			} else {
				seekReader, err = fileSystem.Open(ppath)
				if err != nil {
					log.Error(err)
					return err
				}
			}
			defer func() {
				seekReader.Close()
				if err != nil {
					return
				}
				parts := []string{
					md5,                                   // md5
					param.AsString(0),                     // taskStartTime
					param.AsString(time.Now().Unix()),     // taskEndTime
					param.AsString(info.ModTime().Unix()), // fileModifyTime
					param.AsString(info.Size()),           // fileSize
				}
				err = db.Put(dbKey, com.Str2bytes(strings.Join(parts, `||`)), nil)
				if err != nil {
					log.Error(err)
				}
			}()
			err = cloudbackup.RetryablePut(ctx, mgr, seekReader, objectName, info.Size())
			return err
		}
		if fileSystem == nil {
			err = filepath.Walk(sourcePath, func(ppath string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				return putFile(ppath, info)
			})
		} else {
			err = recursiveDir(sourcePath, fileSystem, putFile)
		}
		if err != nil {
			if err == echo.ErrExit {
				errMsg := ctx.T(`强制退出全量备份`)
				notice.Send(username, notice.NewMessageWithValue(msgType, noticeTitle, errMsg, notice.StateFailure))
			} else {
				notice.Send(username, notice.NewMessageWithValue(msgType, noticeTitle, err.Error(), notice.StateFailure))
				recv.UpdateFields(nil, echo.H{
					`result`: err.Error(),
					`status`: `failure`,
				}, `id`, recv.Id)
			}
		} else {
			successMsg := ctx.T(`全量备份完成`)
			notice.Send(username, notice.NewMessageWithValue(msgType, noticeTitle, successMsg, notice.StateSuccess))
			recv.UpdateFields(nil, echo.H{
				`result`: successMsg,
				`status`: `idle`,
			}, `id`, recv.Id)
		}
		fullBackupExit.Store(false)
	}()
	return nil
}

func recursiveDir(ppath string, fileSystem http.FileSystem, fileFn func(string, os.FileInfo) error) error {
	fp, err := fileSystem.Open(ppath)
	if err != nil {
		return err
	}
	var infos []fs.FileInfo
	infos, err = fp.Readdir(-1)
	fp.Close()
	if err != nil {
		return err
	}
	for _, info := range infos {
		filePath := path.Join(ppath, info.Name())
		if info.IsDir() {
			err = recursiveDir(filePath, fileSystem, fileFn)
			if err != nil && err != filepath.SkipDir {
				return err
			}
			continue
		}
		fmt.Println(filePath)
		err = fileFn(filePath, info)
		if err != nil && err != filepath.SkipDir {
			return err
		}
	}
	return nil
}
