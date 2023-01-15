package upload

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// 分片上传
func (c *ChunkUpload) Upload(r *http.Request, opts ...ChunkInfoOpter) (int64, error) {
	info := &ChunkInfo{
		FormField: `file`,
	}
	for _, opt := range opts {
		opt(info)
	}
	info.Init(r.FormValue, r.Header.Get)
	if !c.IsSupported(info) {
		return 0, ErrChunkUnsupported
	}
	// 获取上传文件
	upFile, fileHeader, err := r.FormFile(info.FormField)
	if err != nil {
		return 0, fmt.Errorf("上传文件错误: %w", err)
	}
	info.FileName = fileHeader.Filename
	info.CurrentSize = uint64(fileHeader.Size)
	defer upFile.Close()
	return c.ChunkUpload(info, upFile)
}

func (c *ChunkUpload) IsSupported(info ChunkInfor) bool {
	err := c.check(info, true)
	if err == nil {
		return true
	}
	return !errors.Is(err, ErrChunkUnsupported)
}

func (c *ChunkUpload) checkSize(info ChunkInfor) error {
	if c.FileMaxBytes > 0 && info.GetFileTotalBytes() > c.FileMaxBytes {
		return fmt.Errorf(`%w: %d>%d `, ErrFileSizeExceedsLimit, info.GetFileTotalBytes(), c.FileMaxBytes)
	}
	if info.GetChunkIndex() == info.GetFileTotalChunks()-1 {
		chunksNoLast := info.GetFileTotalChunks() - 1
		chunkSize := (info.GetFileTotalBytes() - info.GetFileChunkBytes()) / chunksNoLast
		calcTotalBytes := chunksNoLast*chunkSize + info.GetFileChunkBytes()
		if calcTotalBytes > info.GetFileTotalBytes() {

			if log.IsEnabled(log.LevelDebug) {
				log.Debug(com.Dump(echo.H{
					`chunksNoLast`: chunksNoLast, `chunkSize`: chunkSize,
					`chunkBytes`:     info.GetFileChunkBytes(),
					`totalBytes`:     info.GetFileTotalBytes(),
					`calcTotalBytes`: calcTotalBytes,
				}, false))
			}

			return fmt.Errorf(`%w: 文件的最后一个分片尺寸(%d)不正确导致总尺寸超标(%d>%d)`, ErrIncorrectSize, info.GetFileChunkBytes(), calcTotalBytes, info.GetFileTotalBytes())
		}
	} else {
		subtotal := (info.GetChunkIndex() + 1) * info.GetFileChunkBytes()
		if subtotal >= info.GetFileTotalBytes() {
			return fmt.Errorf(`%w: 文件的分片尺寸与分片数量的乘积已经超过了总尺寸(%d>=%d)`, ErrIncorrectSize, subtotal, info.GetFileTotalBytes())
		}
	}
	return nil
}

func (c *ChunkUpload) check(info ChunkInfor, ignoreCurrentSize ...bool) error {
	if info.GetFileTotalBytes() < 1 {
		return fmt.Errorf(`%w: FileTotalBytes less than 1`, ErrChunkUnsupported)
	}
	if (len(ignoreCurrentSize) == 0 || !ignoreCurrentSize[0]) && info.GetCurrentSize() < 1 {
		return fmt.Errorf(`%w: CurrentSize less than 1`, ErrChunkUnsupported)
	}
	if info.GetFileChunkBytes() < 1 {
		return fmt.Errorf(`%w: FileChunkBytes less than 1`, ErrChunkUnsupported)
	}
	if info.GetFileTotalChunks() < 1 {
		return fmt.Errorf(`%w: FileTotalChunks less than 1`, ErrChunkUnsupported)
	}
	return nil
}

func (c *ChunkUpload) ChunkFilename(chunkIndex int) string {
	return filepath.Join(c.TempDir, c.GetUIDString(), fmt.Sprintf("%s_%d", c.fileOriginalName, chunkIndex))
}

func (c *ChunkUpload) Validate(info ChunkInfor) error {
	err := c.check(info)
	if err != nil {
		return err
	}
	err = c.checkSize(info)
	return err
}

func (c *ChunkUpload) statFileDir(chunkFileDir string) string {
	statFileDir := filepath.Join(chunkFileDir, `.stat`)
	return statFileDir
}

// 分片上传
func (c *ChunkUpload) ChunkUpload(info ChunkInfor, upFile io.ReadSeeker) (int64, error) {
	if err := c.Validate(info); err != nil {
		return 0, err
	}
	fileOriginalName := filepath.Base(info.GetFileName())
	c.fileOriginalName = fileOriginalName
	if len(c.savePath) > 0 && filepath.Base(c.savePath) == fileOriginalName {
		fi, err := os.Stat(c.savePath)
		if err == nil && fi.Size() == int64(info.GetFileTotalBytes()) {
			c.setSaveSize(fi.Size())
			return 0, ErrFileUploadCompleted
		}
	}

	chunkSize := int64(info.GetCurrentSize())

	uid := c.GetUIDString()
	chunkFileDir := filepath.Join(c.TempDir, uid)
	statFileDir := c.statFileDir(chunkFileDir)

	if err := os.MkdirAll(chunkFileDir, os.ModePerm); err != nil {
		return 0, err
	}
	os.MkdirAll(statFileDir, os.ModePerm)

	// 新文件创建
	filePath := filepath.Join(chunkFileDir, fmt.Sprintf("%s_%d", fileOriginalName, info.GetChunkIndex()))
	if log.IsEnabled(log.LevelDebug) {
		log.Debug(filePath+`: `, com.Dump(info, false))
	}

	// 获取现在文件大小
	fi, err := os.Stat(filePath)
	var size int64
	if err != nil {
		if !os.IsNotExist(err) {
			return 0, err
		}
	} else {
		size = fi.Size()
	}
	// 判断文件是否传输完成
	if size == chunkSize {
		return 0, fmt.Errorf("%w: %s (size: %d bytes)", ErrChunkUploadCompleted, filepath.Base(filePath), size)
	}
	if size > chunkSize { // 清理异常尺寸分片文件
		size = 0
		os.Remove(filePath)
	}
	start := size
	saveStart := size
	if size > 0 {
		offset := int64(info.GetChunkOffsetBytes())
		if offset > 0 && offset <= chunkSize {
			start = 0
			saveStart = offset
		}
	}

	// 进行断点上传
	// 打开之前上传文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, fmt.Errorf("%w: %s: %v", ErrChunkHistoryOpenFailed, filePath, err)
	}

	// 将数据写入文件
	total, err := uploadFile(upFile, start, file, saveStart)

	file.Close()

	if err == nil && total == chunkSize {
		err = c.recordFinished(chunkFileDir, fileOriginalName, info.GetChunkIndex(), total)
		if err != nil {
			log.Error(err)
		}
		var finished bool
		finished, err = c.isFinish(info, fileOriginalName)
		if finished {
			err = c.MergeAll(info.GetFileTotalChunks(), info.GetFileChunkBytes(), info.GetFileTotalBytes(), fileOriginalName)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return total, err
}

// 上传文件
func uploadFile(upfile io.ReadSeeker, upSeek int64, file *os.File, fSeek int64) (int64, error) {
	// 设置上传偏移量
	upfile.Seek(upSeek, 0)
	// 设置文件偏移量
	file.Seek(fSeek, 0)
	return WriteTo(upfile, file)
}

func WriteTo(r io.Reader, w io.Writer) (n int64, err error) {
	defer func() {
		if e := recover(); e != nil {
			if err == nil {
				err = fmt.Errorf(`%v`, e)
				return
			}
			err = fmt.Errorf(`%w: %v`, err, e)
		}
	}()
	data := make([]byte, 1024)
	n, err = io.CopyBuffer(w, r, data)
	return
}
