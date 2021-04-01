package upload

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

var (
	_fileRWLock com.Oncer
	_fileRWOnce sync.Once
)

func fileRWLock() com.Oncer {
	_fileRWOnce.Do(func() {
		_fileRWLock = com.NewOnce()
	})
	return _fileRWLock
}

type ChunkUpload struct {
	TempDir           string
	SaveDir           string
	TempLifetime      time.Duration
	UID               interface{} // number or string
	fileNameGenerator FileNameGenerator
	fileOriginalName  string
	savePath          string
	saveSize          int64
	merged            bool
	asyncMerge        sql.NullBool // 以否采用异步方式进行合并
	ctx               context.Context
	cancel            context.CancelFunc
}

func (c *ChunkUpload) GetUIDString() string {
	uid := param.AsString(c.UID)
	if len(uid) == 0 {
		uid = `0`
	}
	return uid
}

func (c *ChunkUpload) SetFileNameGenerator(generator FileNameGenerator) *ChunkUpload {
	c.fileNameGenerator = generator
	return c
}

func (c *ChunkUpload) FileNameGenerator() FileNameGenerator {
	if c.fileNameGenerator == nil {
		return DefaultNameGenerator
	}
	return c.fileNameGenerator
}

func (c *ChunkUpload) SetSavePath(savePath string) *ChunkUpload {
	c.savePath = savePath
	return c
}

func (c *ChunkUpload) GetSavePath() string {
	return c.savePath
}

func (c *ChunkUpload) GetSaveSize() int64 {
	return c.saveSize
}

func (c *ChunkUpload) GetFileOriginalName() string {
	return c.fileOriginalName
}

func (c *ChunkUpload) Merged() bool {
	return c.merged
}

func (c *ChunkUpload) IsAsyncMerge() bool {
	if !c.asyncMerge.Valid {
		return true
	}
	return c.asyncMerge.Bool
}

func (c *ChunkUpload) SetAsyncMerge(async bool) *ChunkUpload {
	c.asyncMerge.Bool = async
	c.asyncMerge.Valid = true
	return c
}

func (c *ChunkUpload) StartGC(inverval time.Duration) error {
	c.ctx, c.cancel = context.WithCancel(context.TODO())
	t := time.NewTicker(inverval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			if err := c.gc(); err != nil {
				return err
			}
		case <-c.ctx.Done():
			return nil
		}
	}
}

func (c *ChunkUpload) StopGC() {
	if c.cancel != nil {
		c.cancel()
	}
}

func (c *ChunkUpload) gc() error {
	err := filepath.Walk(c.TempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if time.Since(info.ModTime()) > c.TempLifetime {
			rErr := os.Remove(path)
			if rErr != nil {
				log.Warnf(`[分片文件垃圾回收] %s 删除失败: %v`, path, rErr)
			} else {
				log.Infof(`[分片文件垃圾回收] %s 删除成功`, path)
			}
		}
		return err
	})
	return err
}
