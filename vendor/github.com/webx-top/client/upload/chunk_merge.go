package upload

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

func (c *ChunkUpload) clearChunk(chunkTotal uint64, fileName, fileUUID string) error {
	uid := c.GetUIDString()
	chunkFileDir := filepath.Join(c.TempDir, uid)
	for i := uint64(0); i < chunkTotal; i++ {
		finishedFlag := c.finishedFlagFile(chunkFileDir, fileName, fileUUID, i)
		os.Remove(finishedFlag)
		chunkFile := c.chunkFilePath(chunkFileDir, fileName, fileUUID, i)
		log.Debugf("删除分片文件: %s", chunkFile)
		err := os.Remove(chunkFile)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("%w: %s: %v", ErrChunkFileDeleteFailed, chunkFile, err)
			}
		}
	}
	totalFile := c.totalFile(chunkFileDir, fileName, fileUUID)
	os.Remove(totalFile)
	return nil
}

func (c *ChunkUpload) totalFile(chunkFileDir, fileName, fileUUID string) string {
	totalFile := filepath.Join(c.statFileDir(chunkFileDir), fileName+".total.txt")
	return totalFile
}

func (c *ChunkUpload) finishedFlagFile(chunkFileDir, fileName, fileUUID string, chunkIndex uint64) string {
	filename := c.chunkFileNameWithoutExt(fileName, fileUUID, chunkIndex)
	finishedFlag := filepath.Join(c.statFileDir(chunkFileDir), filename+".finished")
	return finishedFlag
}

func (c *ChunkUpload) recordFinished(chunkFileDir, fileName, fileUUID string, chunkIndex uint64, size int64) error {
	finishedFlag := c.finishedFlagFile(chunkFileDir, fileName, fileUUID, chunkIndex)
	return os.WriteFile(finishedFlag, []byte(param.AsString(size)), os.ModePerm)
}

func (c *ChunkUpload) existsFinishedFlag(
	chunkFileDir, fileName, fileUUID string, chunkIndex uint64,
	chunkFileModTime time.Time,
) bool {
	finishedFlag := c.finishedFlagFile(chunkFileDir, fileName, fileUUID, chunkIndex)
	fi, err := os.Stat(finishedFlag)
	return err == nil && !fi.IsDir() && fi.ModTime().After(chunkFileModTime)
}

func (c *ChunkUpload) calcFinisedSize(info ChunkInfor, fileName string) (uint64, error) {
	fileSize := info.GetFileTotalBytes()
	uid := c.GetUIDString()
	chunkFileDir := filepath.Join(c.TempDir, uid)
	totalFile := c.totalFile(chunkFileDir, fileName, info.GetFileUUID())
	var finishedSize uint64
	b, err := os.ReadFile(totalFile)
	if err == nil {
		finishedSize = param.AsUint64(string(b))
		if finishedSize == fileSize {
			return fileSize, err
		}
	}
	err = nil
	finishedSize = 0
	var finishedCount uint64
	chunkTotal := info.GetFileTotalChunks()
	for i := uint64(0); i < chunkTotal; i++ {
		chunkFile := c.chunkFilePath(chunkFileDir, fileName, info.GetFileUUID(), i)
		// 分片大小获取
		fi, err := os.Stat(chunkFile)
		if err != nil {
			if !os.IsNotExist(err) {
				err = fmt.Errorf(`统计分片文件尺寸错误: %s: %w`, chunkFile, err)
				return finishedSize, err
			}
		} else {
			finishedSize += uint64(fi.Size())
			if c.existsFinishedFlag(chunkFileDir, fileName, info.GetFileUUID(), i, fi.ModTime()) {
				finishedCount++
			}
		}
	}
	if log.IsEnabled(log.LevelDebug) {
		log.Debug(fileName+`: `, echo.Dump(echo.H{`finishedSize`: finishedSize, `fileSize`: fileSize, `finishedCount`: finishedCount}, false))
	}
	if finishedCount == chunkTotal {
		if finishedSize != fileSize {
			finishedSize = fileSize
		}
	}
	if finishedSize > 0 {
		err = os.WriteFile(totalFile, []byte(param.AsString(finishedSize)), os.ModePerm)
	}
	return finishedSize, err
}

// 判断是否完成  根据现有文件的大小 与 上传文件大小进行匹配
func (c *ChunkUpload) isFinish(info ChunkInfor, fileName string, counter ...*int) (bool, error) {
	fileSize := info.GetFileTotalBytes()
	uid := c.GetUIDString()
	flag := `chunkUpload.saveFileSizeInfo.` + uid + `.` + fileName
	value, err, shared := chunkSg.Do(flag, func() (interface{}, error) {
		finishedSize, err := c.calcFinisedSize(info, fileName)
		return finishedSize, err
	})
	finishedSize := value.(uint64)
	finished := finishedSize == fileSize
	if err != nil || finished || !shared {
		return finished, err
	}
	if len(counter) == 0 {
		retries := 0
		counter = []*int{&retries}
	} else {
		*(counter[0])++
	}
	if log.IsEnabled(log.LevelDebug) {
		log.Debug(echo.Dump(echo.H{
			`finishedSize`: finishedSize, `fileSize`: fileSize, `wait`: true, `retries`: *(counter[0]),
			`fileName`: fileName + `_` + strconv.FormatUint(info.GetChunkIndex(), 10),
		}, false))
	}
	if *(counter[0]) > 3 {
		return false, nil
	}
	log.Debugf(`[isFinish()] %s retry: %d`, c.chunkFileNameByInfo(fileName, info), *(counter[0]))
	time.Sleep(time.Millisecond * 20)
	return c.isFinish(info, fileName, counter...)
}

func GenSavePath(saveDir string, saveFileName string, namer FileNameGenerator) (string, error) {
	saveName, err := namer(saveFileName)
	if err != nil {
		return ``, err
	}
	savePath := filepath.Join(saveDir, saveName)
	return savePath, nil
}

func (c *ChunkUpload) getOrGenSavePath(saveFileName string) (string, error) {
	if len(c.savePath) == 0 {
		savePath, err := GenSavePath(c.SaveDir, saveFileName, c.FileNameGenerator())
		if err != nil {
			return ``, err
		}
		c.savePath = savePath
	}
	return c.savePath, nil
}

func (c *ChunkUpload) prepareSavePath(saveFileName string) error {
	savePath, err := c.getOrGenSavePath(saveFileName)
	if err != nil {
		return err
	}
	saveDir := filepath.Dir(savePath)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func (c *ChunkUpload) existsDistFile(saveFileName string) (bool, error) {
	savePath, err := c.getOrGenSavePath(saveFileName)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(savePath)
	if err == nil {
		return true, err
	}
	if !os.IsNotExist(err) {
		return false, err
	}
	return false, nil
}

func (c *ChunkUpload) doMergeAllCase1(flag string, ctx context.Context, info ChunkInfor, saveFileName string) (err error) {
	done := make(chan error)
	isNew := chunkDelay.Do(ctx, flag, func() error {
		var err error
		defer func() {
			done <- err
			close(done)
		}()
		var exists bool
		exists, err = c.existsDistFile(saveFileName)
		if err != nil || exists {
			return err
		}
		if c.onBeforeMerge != nil {
			if err = c.onBeforeMerge(ctx, info, saveFileName); err != nil {
				return err
			}
		}
		err = c.mergeAll(info, saveFileName)
		if err != nil {
			return err
		}
		if c.onAfterMerge != nil {
			err = c.onAfterMerge(ctx, info, saveFileName)
		}
		return err
	})
	if isNew {
		err = <-done
	} else {
		close(done)
	}
	return
}

func (c *ChunkUpload) doMergeAllCase2(flag string, ctx context.Context, info ChunkInfor, saveFileName string) (err error) {
	_, err, _ = chunkSg.Do(flag, func() (interface{}, error) {
		exists, err := c.existsDistFile(saveFileName)
		if err != nil || exists {
			return nil, err
		}
		if c.onBeforeMerge != nil {
			if err := c.onBeforeMerge(ctx, info, saveFileName); err != nil {
				return nil, err
			}
		}
		err = c.mergeAll(info, saveFileName)
		if err != nil {
			return nil, err
		}
		if c.onAfterMerge != nil {
			err = c.onAfterMerge(ctx, info, saveFileName)
		}
		return nil, err
	})
	return
}

// 合并某个文件的所有切片
func (c *ChunkUpload) MergeAll(ctx context.Context, info ChunkInfor, saveFileName string) error {
	uid := c.GetUIDString()
	flag := `chunkUpload.mergeAll.` + uid + `.` + saveFileName
	if c.DelayMerge {
		return c.doMergeAllCase1(flag, ctx, info, saveFileName)
	}
	return c.doMergeAllCase2(flag, ctx, info, saveFileName)
}

func (c *ChunkUpload) setSaveSize(n int64) {
	c.mu.Lock()
	c.saveSize = n
	c.mu.Unlock()
}

func (c *ChunkUpload) addSaveSize(n int64) {
	c.mu.Lock()
	c.saveSize += n
	c.mu.Unlock()
}

func (c *ChunkUpload) mergeAll(info ChunkInfor, saveFileName string) (err error) {
	totalChunks := info.GetFileTotalChunks()
	c.setSaveSize(0)
	if err = os.MkdirAll(c.SaveDir, os.ModePerm); err != nil {
		return
	}
	if err = c.prepareSavePath(saveFileName); err != nil {
		return
	}
	// 打开之前上传文件
	var file *os.File
	file, err = os.OpenFile(c.savePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		err = fmt.Errorf("%w: %s: %v (mergeAll)", ErrChunkMergeFileCreateFailed, c.savePath, err)
		return
	}
	defer file.Close()
	uid := c.GetUIDString()
	chunkFileDir := filepath.Join(c.TempDir, uid)
	for chunkIndex := uint64(0); chunkIndex < totalChunks; chunkIndex++ {
		chunkFilePath := c.chunkFilePath(chunkFileDir, saveFileName, info.GetFileUUID(), chunkIndex)
		cfile, cerr := os.Open(chunkFilePath)
		if cerr != nil {
			err = fmt.Errorf("%w: %s: %v", ErrChunkFileOpenFailed, chunkFilePath, cerr)
			return
		}
		var n int64
		n, err = WriteTo(cfile, file)

		cfile.Close()

		if err != nil {
			err = fmt.Errorf("%w: %s: %v", ErrChunkFileMergeFailed, chunkFilePath, err)
			return
		}
		c.addSaveSize(n)
	}
	err = file.Sync()
	if err != nil {
		return
	}

	err = c.clearChunk(totalChunks, saveFileName, info.GetFileUUID())
	c.merged = true
	log.Debugf("分片文件合并完毕: %s", c.savePath)
	return
}
