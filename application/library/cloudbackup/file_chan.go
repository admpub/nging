package cloudbackup

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/flock"
	"github.com/admpub/nging/v5/application/library/msgbox"
	"github.com/admpub/nging/v5/application/model"
	"github.com/admpub/once"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

var (
	BackupTasks  = param.NewMap()
	fileChan     chan *PutFile
	fileChanOnce once.Once
	ctx          context.Context
	cancel       context.CancelFunc
	delay        = com.NewDelayOnce(time.Second*3, time.Minute*5)
)

type PutFile struct {
	Manager           Storager
	Config            dbschema.NgingCloudBackup
	ObjectName        string
	FilePath          string
	Operation         string
	WaitFillCompleted bool
}

func (mf *PutFile) Do(ctx context.Context) error {
	fp, err := os.OpenFile(mf.FilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Error(`Open ` + mf.FilePath + `: ` + err.Error())
		return err
	}
	defer fp.Close()
	if !mf.WaitFillCompleted || flock.IsCompleted(fp, time.Now()) {
		var fi os.FileInfo
		fi, err = fp.Stat()
		if err != nil {
			log.Error(`Stat ` + mf.FilePath + `: ` + err.Error())
			return err
		}
		size := fi.Size()
		err = RetryablePut(ctx, mf.Manager, fp, mf.ObjectName, size)
		if err != nil {
			log.Error(`s3manager.Put ` + mf.FilePath + ` (size:` + strconv.FormatInt(size, 10) + `): ` + err.Error())
		} else {
			log.Info(`s3manager.Put ` + mf.FilePath + ` (size:` + strconv.FormatInt(size, 10) + `): success`)
		}
	}
	return err
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
				exec := func() error {
					ctx := defaults.NewMockContext()
					startTime := time.Now()
					err := mf.Do(ctx)
					RecordLog(ctx, err, &mf.Config, mf.FilePath, mf.ObjectName, mf.Operation, startTime)
					return nil
				}
				var delayDur time.Duration
				if mf.Config.Delay > 0 {
					delayDur = time.Second * time.Duration(mf.Config.Delay+2)
				}
				delay.Do(ctx, mf.FilePath, exec, delayDur)
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
	startTime time.Time, backupType ...string) {
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
		msgbox.Success(`Cloud-Backup`, `Close: `+com.String(id))
	}
	return nil
}
