package cloudbackup

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/admpub/checksum"
	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"

	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/flock"
	"github.com/admpub/nging/v5/application/library/msgbox"
	"github.com/admpub/nging/v5/application/model"
)

var (
	BackupTasks  = param.NewMap()
	fileChan     chan *PutFile
	fileChanOnce once.Once
	ctx          context.Context
	cancel       context.CancelFunc
)

type PutFile struct {
	Manager           Storager
	Config            dbschema.NgingCloudBackup
	ObjectName        string
	FilePath          string
	Operation         string
	WaitFillCompleted bool
}

func (mf *PutFile) Do(ctx context.Context) (size int64, lastModtime time.Time, err error) {
	var fp *os.File
	fp, err = os.OpenFile(mf.FilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error(`Open ` + mf.FilePath + `: ` + err.Error())
		return
	}
	defer fp.Close()
	if !mf.WaitFillCompleted || flock.IsCompleted(fp, time.Now()) {
		var fi os.FileInfo
		fi, err = fp.Stat()
		if err != nil {
			log.Error(`Stat ` + mf.FilePath + `: ` + err.Error())
			return
		}
		size = fi.Size()
		lastModtime = fi.ModTime()
		err = RetryablePut(ctx, mf.Manager, fp, mf.ObjectName, size)
		if err != nil {
			log.Error(`s3manager.Put ` + mf.FilePath + ` (size:` + strconv.FormatInt(size, 10) + `): ` + err.Error())
		} else {
			log.Info(`s3manager.Put ` + mf.FilePath + ` (size:` + strconv.FormatInt(size, 10) + `): success`)
		}
	}
	return
}

func FileChan() chan *PutFile {
	fileChanOnce.Do(initFileChan)
	return fileChan
}

type ErrIsAccessDenied interface {
	ErrIsAccessDenied(error) bool
}

func RetryablePut(ctx context.Context, mgr Storager, fp *os.File, objectName string, size int64) error {
	return common.OnErrorRetry(func() error {
		err := mgr.Put(ctx, fp, objectName, size)
		if cli, ok := mgr.(ErrIsAccessDenied); ok {
			if cli.ErrIsAccessDenied(err) {
				if connErr := mgr.Connect(); connErr != nil {
					log.Error(`s3manager.Connect: ` + connErr.Error())
				}
			}
		}
		fp.Seek(0, 0)
		return err
	}, 3, time.Second*2)
}

func initFileChan() {
	fileChan = make(chan *PutFile, 1000)
	ctx, cancel = context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case mf, ok := <-fileChan:
				if !ok || mf == nil {
					return
				}
				startTime := time.Now()
				ctx := defaults.NewMockContext()
				db, err := LevelDB().OpenDB(mf.Config.Id)
				if err != nil {
					err = fmt.Errorf(`failed to open levelDB file: %w`, err)
					log.Errorf(`[cloundbackup] %v`, err)
					RecordLog(ctx, err, &mf.Config, mf.FilePath, mf.ObjectName, mf.Operation, startTime, 0)
					continue
				}
				dbKey := com.Str2bytes(mf.FilePath)
				var (
					md5                                    string
					startTs, endTs, fileModifyTs, fileSize int64
				)
				var nowFileModifyTs, nowFileSize int64
				val, err := db.Get(dbKey, nil)
				if err != nil {
					if err != leveldb.ErrNotFound {
						err = fmt.Errorf(`failed to read data from levelDB: %w`, err)
						log.Errorf(`[cloundbackup] %v`, err)
						RecordLog(ctx, err, &mf.Config, mf.FilePath, mf.ObjectName, mf.Operation, startTime, 0)
						continue
					}
				} else {
					md5, startTs, endTs, fileModifyTs, fileSize = ParseDBValue(val)
					fi, err := os.Stat(mf.FilePath)
					if err == nil {
						nowFileModifyTs = fi.ModTime().Unix()
						nowFileSize = fi.Size()
						if fileModifyTs == nowFileModifyTs && fileSize == nowFileSize {
							continue
						}
					} else {
						nowFileModifyTs = fileModifyTs
						nowFileSize = fileSize
					}
					if startTs > 0 {
						continue
					}
					if endTs > 0 && endTs < startTime.Unix()-int64(mf.Config.MinModifyInterval) {
						continue
					}
				}
				parts := []string{
					md5,                              // md5
					param.AsString(startTime.Unix()), // taskStartTime
					param.AsString(0),                // taskEndTime
					param.AsString(nowFileModifyTs),  // fileModifyTime
					param.AsString(nowFileSize),      // fileSize
				}
				err = db.Put(dbKey, com.Str2bytes(strings.Join(parts, `||`)), nil)
				if err != nil {
					err = fmt.Errorf(`failed to write data to levelDB: %w`, err)
					log.Errorf(`[cloundbackup] %v`, err)
					RecordLog(ctx, err, &mf.Config, mf.FilePath, mf.ObjectName, mf.Operation, startTime, 0)
					continue
				}
				var filemtime time.Time
				fileSize, filemtime, err = mf.Do(ctx)
				RecordLog(ctx, err, &mf.Config, mf.FilePath, mf.ObjectName, mf.Operation, startTime, uint64(fileSize))
				if err == nil {
					md5, _ = checksum.MD5sum(mf.FilePath)
					parts := []string{
						md5,                               // md5
						param.AsString(0),                 // taskStartTime
						param.AsString(time.Now().Unix()), // taskEndTime
						param.AsString(filemtime.Unix()),  // fileModifyTime
						param.AsString(fileSize),          // fileSize
					}
					err := db.Put(dbKey, com.Str2bytes(strings.Join(parts, `||`)), nil)
					if err != nil {
						log.Errorf(`[cloundbackup] failed to write data to levelDB: %v`, err)
					}
				}
			}
		}
	}()
}

func ResetFileChan() {
	cancel()
	fileChanOnce.Reset()
}

func RecordLog(ctx echo.Context, err error, cfg *dbschema.NgingCloudBackup,
	filePath string, remotePath string, operation string,
	startTime time.Time, size uint64, backupType ...string) {
	if cfg.LogDisabled == `N` && (cfg.LogType == model.CloudBackupLogTypeAll || err != nil) {
		if ctx == nil {
			ctx = defaults.NewMockContext()
		}
		logM := model.NewCloudBackupLog(ctx)
		logM.BackupId = cfg.Id
		if len(backupType) > 0 && len(backupType[0]) > 0 {
			logM.BackupType = backupType[0]
		} else {
			logM.BackupType = model.CloudBackupTypeChange
		}
		logM.BackupFile = filePath
		logM.RemoteFile = remotePath
		logM.Operation = operation
		logM.Size = size
		logM.Elapsed = uint(time.Since(startTime).Milliseconds())
		if err != nil {
			logM.Error = err.Error()
			logM.Status = model.CloudBackupStatusFailure
		} else {
			logM.Status = model.CloudBackupStatusSuccess
		}
		if _, err := logM.Add(); err != nil {
			log.Error(err)
		}
	}
}

func MonitorBackupStop(id uint) error {
	if monitor, ok := BackupTasks.Get(id).(*com.MonitorEvent); ok {
		monitor.Close()
		BackupTasks.Delete(id)
		LevelDB().CloseDB(id)
		msgbox.Success(`Cloud-Backup`, `Close: `+com.String(id))
	}
	return nil
}
