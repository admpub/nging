package cloud

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/nging/application/library/flock"
	"github.com/admpub/nging/application/library/s3manager"
	"github.com/webx-top/echo/param"
)

var (
	backupTasks  = param.NewMap()
	fileChan     chan *PutFile
	fileChanOnce sync.Once
	ctx          context.Context
	cancel       context.CancelFunc
)

type PutFile struct {
	Manager    *s3manager.S3Manager
	ObjectName string
	FilePath   string
}

func (mf *PutFile) Do() error {
	fp, err := os.Open(mf.FilePath)
	if err != nil {
		log.Error(`Open ` + mf.FilePath + `: ` + err.Error())
		return err
	}
	defer fp.Close()
	fi, err := fp.Stat()
	if err != nil {
		log.Error(`Stat ` + mf.FilePath + `: ` + err.Error())
		return err
	}
	if flock.IsCompleted(fp, fi, time.Now()) {
		err = mf.Manager.Put(fp, mf.ObjectName, fi.Size())
		if err != nil {
			log.Error(`s3manager.Put ` + mf.FilePath + `: ` + err.Error())
		}
	}
	return err
}

func FileChan() chan *PutFile {
	fileChanOnce.Do(initFileChan)
	return fileChan
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
				mf.Do()
			}
		}
	}()
}

func ResetFileChan() {
	cancel()
	fileChanOnce = sync.Once{}
}
