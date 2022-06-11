package upload

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/admpub/log"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

// 合并切片文件
func (c *ChunkUpload) merge(chunkIndex uint64, fileChunkBytes uint64, file *os.File, savePath string) (int64, error) {
	uid := c.GetUIDString()
	// 设置文件写入偏移量
	file.Seek(int64(fileChunkBytes*chunkIndex), 0)

	fileName := filepath.Base(file.Name())

	chunkFilePath := filepath.Join(c.TempDir, uid, fmt.Sprintf(`%s_%d`, fileName, chunkIndex))

	chunkFileObj, err := os.Open(chunkFilePath)
	if err != nil {
		return 0, fmt.Errorf("%w: %s: %v", ErrChunkFileOpenFailed, chunkFilePath, err)
	}
	var n int64
	n, err = WriteTo(chunkFileObj, file)

	chunkFileObj.Close()

	if err != nil {
		return n, fmt.Errorf("%w: %s: %v", ErrChunkFileMergeFailed, chunkFilePath, err)
	}

	log.Debugf("分片文件合并成功: %s", chunkFilePath)
	return n, err
}

func (c *ChunkUpload) clearChunk(chunkTotal uint64, fileName string) error {
	uid := c.GetUIDString()
	chunkFileDir := filepath.Join(c.TempDir, uid)
	for i := uint64(0); i < chunkTotal; i++ {
		chunkFile := filepath.Join(chunkFileDir, fileName+"_"+param.AsString(i))
		log.Debugf("删除分片文件: %s", chunkFile)
		err := os.Remove(chunkFile)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("%w: %s: %v", ErrChunkFileDeleteFailed, chunkFile, err)
			}
		}
	}
	totalFile := filepath.Join(chunkFileDir, fileName+".total")
	os.Remove(totalFile)
	return nil
}

// 判断是否完成  根据现有文件的大小 与 上传文件大小进行匹配
func (c *ChunkUpload) isFinish(info ChunkInfor, fileName string, counter ...*int) (bool, error) {
	fileSize := info.GetFileTotalBytes()
	uid := c.GetUIDString()
	chunkFileDir := filepath.Join(c.TempDir, uid)
	totalFile := filepath.Join(chunkFileDir, fileName+".total")
	flag := `chunkUpload.saveFileSizeInfo.` + uid + `.` + fileName
	if !fileRWLock().CanSet(flag) {
		fileRWLock().Wait(flag) // 需要等待创建完成
		b, err := ioutil.ReadFile(totalFile)
		if err != nil {
			if !os.IsNotExist(err) {
				err = fmt.Errorf(`读取分片统计结果文件出错: %s: %v`, totalFile, err)
			}
			return false, err
		}
		chunkSize := param.AsInt64(string(b))
		if log.IsEnabled(log.LevelDebug) {
			log.Debug(echo.Dump(echo.H{`chunkSize`: chunkSize, `fileSize`: fileSize, `wait`: true}, false))
		}
		if chunkSize == int64(fileSize) {
			return false, nil // 说明以前的已经判断为完成了，后面堵塞住的统一返回false避免重复执行
		}
		if len(counter) == 0 {
			retries := 0
			counter = []*int{&retries}
		} else {
			*(counter[0])++
		}
		if *(counter[0]) > 1000 {
			return false, nil
		}
		log.Debugf(`[isFinish()] %s_%d retry: %d`, fileName, info.GetChunkIndex(), *(counter[0]))
		return c.isFinish(info, fileName, counter...)
	}
	defer fileRWLock().Release(flag)
	var chunkSize int64
	chunkTotal := info.GetFileTotalChunks()
	for i := uint64(0); i < chunkTotal; i++ {
		chunkFile := filepath.Join(chunkFileDir, fileName+"_"+param.AsString(i))
		// 分片大小获取
		fi, err := os.Stat(chunkFile)
		if err != nil {
			if !os.IsNotExist(err) {
				err = fmt.Errorf(`统计分片文件尺寸错误: %s: %v`, chunkFile, err)
				return false, err
			}
		} else {
			chunkSize += fi.Size()
		}
	}
	if log.IsEnabled(log.LevelDebug) {
		log.Debug(echo.Dump(echo.H{`chunkSize`: chunkSize, `fileSize`: fileSize}, false))
	}
	err := ioutil.WriteFile(totalFile, []byte(param.AsString(chunkSize)), os.ModePerm)
	return chunkSize == int64(fileSize), err
}

func (c *ChunkUpload) prepareSavePath(saveFileName string) error {
	if len(c.savePath) == 0 {
		saveName, err := c.FileNameGenerator()(saveFileName)
		if err != nil {
			return err
		}
		c.savePath = filepath.Join(c.SaveDir, saveName)
	}
	saveDir := filepath.Dir(c.savePath)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// 合并某个文件的所有切片
func (c *ChunkUpload) MergeAll(totalChunks uint64, fileChunkBytes uint64, saveFileName string, async bool) (err error) {
	uid := c.GetUIDString()
	flag := `chunkUpload.mergeAll.` + uid + `.` + saveFileName
	if !fileRWLock().CanSet(flag) {
		return
	}

	err = c.mergeAll(totalChunks, fileChunkBytes, saveFileName, async)
	fileRWLock().Release(flag)
	return
}

func (c *ChunkUpload) mergeAll(totalChunks uint64, fileChunkBytes uint64, saveFileName string, async bool) (err error) {
	c.saveSize = 0
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
	if async {
		file.Close()
		wg := &sync.WaitGroup{}
		mu := sync.RWMutex{}
		for chunkIndex := uint64(0); chunkIndex < totalChunks; chunkIndex++ {
			wg.Add(1)
			go func(chunkIndex uint64) {
				defer wg.Done()
				file, err := os.OpenFile(c.savePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
				if err != nil {
					log.Errorf("%v: %s: %v (mergeAll)", ErrChunkMergeFileCreateFailed, c.savePath, err)
					return
				}
				n, err := c.merge(chunkIndex, fileChunkBytes, file, c.savePath)
				file.Close()
				if err != nil {
					log.Error(err)
				} else {
					mu.Lock()
					c.saveSize += n
					mu.Unlock()
				}
			}(chunkIndex)
		}
		wg.Wait()
		err = c.clearChunk(totalChunks, saveFileName)
		c.merged = true
		log.Debugf("分片文件合并完毕: %s", c.savePath)
		return
	}
	defer file.Close()
	uid := c.GetUIDString()
	chunkFileDir := filepath.Join(c.TempDir, uid)
	for chunkIndex := uint64(0); chunkIndex < totalChunks; chunkIndex++ {
		chunkFilePath := filepath.Join(chunkFileDir, fmt.Sprintf(`%s_%d`, saveFileName, chunkIndex))
		cfile, cerr := os.Open(chunkFilePath)
		if cerr != nil {
			err = fmt.Errorf("%w: %s: %v", ErrChunkFileOpenFailed, chunkFilePath, cerr)
			log.Errorf(err.Error())
			return nil
		}
		var n int64
		n, err = WriteTo(cfile, file)

		cfile.Close()

		if err != nil {
			err = fmt.Errorf("%w: %s: %v", ErrChunkFileMergeFailed, chunkFilePath, err)
			return
		}
		c.saveSize += n
	}

	err = c.clearChunk(totalChunks, saveFileName)
	c.merged = true
	log.Debugf("分片文件合并完毕: %s", c.savePath)
	return
}
