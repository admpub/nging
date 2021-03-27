package upload

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/admpub/errors"
	"github.com/admpub/log"
	"github.com/webx-top/com"
)

var (
	ErrChunkUploadCompleted = errors.New("文件分片已经上传完成")
	ErrChunkUnsupported     = errors.New("不支持分片上传")
)

// 分片上传
func (c *ChunkUpload) Upload(r *http.Request, mapping map[string]string, formFields ...string) (int64, error) {
	formField := `file`
	if len(formFields) > 0 {
		formField = formFields[0]
	}
	// 获取上传文件
	upFile, fileHeader, err := r.FormFile(formField)
	if err != nil {
		return 0, fmt.Errorf("上传文件错误: %w", err)
	}
	defer upFile.Close()
	info := &ChunkInfo{
		Mapping:     mapping,
		FileName:    fileHeader.Filename,
		CurrentSize: uint64(fileHeader.Size),
	}
	info.BatchSetByURLValues(r.Form)
	return c.ChunkUpload(info, upFile)
}

// 分片上传
func (c *ChunkUpload) ChunkUpload(info ChunkInfor, upFile io.ReadSeeker) (int64, error) {
	if info.GetFileChunkBytes() < 1 {
		return 0, fmt.Errorf(`%w: FileChunkBytes less than 1`, ErrChunkUnsupported)
	}
	if info.GetFileTotalBytes() < 1 {
		return 0, fmt.Errorf(`%w: FileTotalBytes less than 1`, ErrChunkUnsupported)
	}
	if info.GetFileChunkBytes() < 1 {
		return 0, fmt.Errorf(`%w: FileChunkBytes less than 1`, ErrChunkUnsupported)
	}

	c.fileOriginalName = info.GetFileName()
	chunkSize := int64(info.GetCurrentSize())

	uid := c.GetUIDString()
	chunkFileDir := filepath.Join(c.TempDir, uid)

	if err := os.MkdirAll(chunkFileDir, os.ModePerm); err != nil {
		return 0, err
	}

	// 新文件创建
	filePath := filepath.Join(chunkFileDir, fmt.Sprintf("%s_%d", c.fileOriginalName, info.GetChunkIndex()))
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
	start := size
	saveStart := size
	offset := int64(info.GetChunkOffsetBytes())
	if offset > 0 {
		start = 0
		saveStart = offset
	}

	// 进行断点上传
	// 打开之前上传文件
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return 0, fmt.Errorf("打开之前上传文件不存在: %w", err)
	}

	defer file.Close()

	// 将数据写入文件
	total, err := uploadFile(upFile, start, file, saveStart)
	if err == nil && total == chunkSize {
		if c.isFinish(info, c.fileOriginalName) {
			_, err = c.MergeAll(info.GetFileTotalChunks(), info.GetFileChunkBytes(), c.fileOriginalName, c.IsAsyncMerge())
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

func WriteTo(r io.Reader, w io.Writer) (int64, error) {
	data := make([]byte, 1024)
	return io.CopyBuffer(w, r, data)
}
