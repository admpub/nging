package cloudbackup

import (
	"context"
	"os"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/dbschema"
	"github.com/admpub/nging/v5/application/library/common"
	"github.com/admpub/nging/v5/application/library/flock"
	"github.com/admpub/nging/v5/application/library/msgbox"
	"github.com/admpub/nging/v5/application/library/s3manager"
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
)

type PutFile struct {
	Manager           *s3manager.S3Manager
	Config            dbschema.NgingCloudBackup
	ObjectName        string
	FilePath          string
	WaitFillCompleted bool
}

func (mf *PutFile) Do(ctx context.Context) error {
	fp, err := os.Open(mf.FilePath)
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
		err = RetryablePut(ctx, mf.Manager, fp, mf.ObjectName, fi.Size())
		if err != nil {
			log.Error(`s3manager.Put ` + mf.FilePath + `: ` + err.Error())
		} else {
			log.Info(`s3manager.Put ` + mf.FilePath + `: success`)
		}
	}
	return err
}

func FileChan() chan *PutFile {
	fileChanOnce.Do(initFileChan)
	return fileChan
}

func RetryablePut(ctx context.Context, mgr *s3manager.S3Manager, fp *os.File, objectName string, size int64) error {
	return common.OnErrorRetry(func() error {
		err := mgr.Put(ctx, fp, objectName, size)
		if mgr.ErrIsAccessDenied(err) {
			if _, connErr := mgr.Connect(); connErr != nil {
				log.Error(`s3manager.Connect: ` + connErr.Error())
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
				err := mf.Do(ctx)
				RecordLog(ctx, err, &mf.Config, mf.FilePath, startTime)
			}
		}
	}()
}

func ResetFileChan() {
	cancel()
	fileChanOnce.Reset()
}

func RecordLog(ctx echo.Context, err error, cfg *dbschema.NgingCloudBackup, filePath string, startTime time.Time, backupType ...string) {
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
