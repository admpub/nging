package upload

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo/param"
	"golang.org/x/sync/singleflight"
)

var (
	chunkSg = singleflight.Group{}
)

func NewChunkUpload(tempDir, mergeSaveDir string, tempLifetime time.Duration) *ChunkUpload {
	return &ChunkUpload{
		TempDir:      tempDir,
		SaveDir:      mergeSaveDir,
		TempLifetime: tempLifetime,
	}
}

type ChunkUpload struct {
	TempDir           string
	SaveDir           string
	TempLifetime      time.Duration
	UID               interface{} // number or string
	FileMaxBytes      uint64
	BasedUUID         bool // 基于fileUUID记录切片文件
	fileNameGenerator FileNameGenerator
	fileOriginalName  string
	savePath          string
	saveSize          int64
	merged            bool
	ctx               context.Context
	cancel            context.CancelFunc
	mu                sync.RWMutex
}

func (c *ChunkUpload) Clone() *ChunkUpload {
	return &ChunkUpload{
		TempDir:           c.TempDir,
		SaveDir:           c.SaveDir,
		TempLifetime:      c.TempLifetime,
		UID:               c.UID,
		FileMaxBytes:      c.FileMaxBytes,
		BasedUUID:         c.BasedUUID,
		fileNameGenerator: c.fileNameGenerator,
	}
}

func (c *ChunkUpload) SetUID(uid interface{}) *ChunkUpload {
	c.UID = uid
	return c
}

func (c *ChunkUpload) SetFileMaxBytes(fileMaxBytes uint64) *ChunkUpload {
	c.FileMaxBytes = fileMaxBytes
	return c
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
	c.mu.RLock()
	size := c.saveSize
	c.mu.RUnlock()
	return size
}

func (c *ChunkUpload) GetFileOriginalName() string {
	return c.fileOriginalName
}

func (c *ChunkUpload) Merged() bool {
	return c.merged
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
				if os.IsNotExist(err) {
					return nil
				}
				log.Warnf(`[分片文件垃圾回收] %s 删除失败: %v`, path, rErr)
			} else {
				log.Infof(`[分片文件垃圾回收] %s 删除成功`, path)
			}
		}
		return err
	})
	return err
}
