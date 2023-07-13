package chunk

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/admpub/nging/v5/application/library/errorslice"
	syncOnce "github.com/admpub/once"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
)

var (
	chunkUploadInitOnce syncOnce.Once
	chunkUpload         *uploadClient.ChunkUpload

	// ChunkTempDir 保存分片文件的临时文件夹
	ChunkTempDir = filepath.Join(echo.Wd(), `data/temp/upload/chunk_temp`)

	// MergeSaveDir 分片文件合并后保存的临时文件夹
	// 避免和最终的文件保存位置相同，上传后一般需要将合并后的文件转移(os.Rename)到最终保存位置
	MergeSaveDir = filepath.Join(echo.Wd(), `data/temp/upload/chunk_merged`)

	// TempLifetime 分片临时文件过期时间
	TempLifetime = 24 * time.Hour

	// GCInterval 垃圾清理间隔时间
	GCInterval = 48 * time.Hour
)

func MergedFileNameGenerator(uid interface{}) uploadClient.FileNameGenerator {
	return func(fileName string) (string, error) {
		return filepath.Join(param.AsString(uid), time.Now().Format("20060102"), fileName), nil
	}
}

var ownerTypes = []string{`customer`, `user`}

func CleanFileByOwner(ownerType string, ownerID uint64) error {
	if !com.InSlice(ownerType, ownerTypes) {
		return fmt.Errorf(`failed to CleanFileByOwner: unsupport ownerType %q`, ownerType)
	}
	ownerIDstr := strconv.FormatUint(ownerID, 10)
	mergedFilePath := filepath.Join(MergeSaveDir, ownerType, ownerIDstr)
	chunkFilePath := filepath.Join(ChunkTempDir, ownerType, ownerIDstr)
	errs := errorslice.New()
	for _, fpath := range []string{chunkFilePath, mergedFilePath} {
		err := os.RemoveAll(fpath)
		if err != nil && !os.IsNotExist(err) {
			errs.Add(err)
		}
	}
	return errs.ToError()
}

func NewUploader(uid interface{}, fileMaxBytes ...uint64) *uploadClient.ChunkUpload {
	chunkUploadInitOnce.Do(initChunkUploader)
	return newUploader(uid, fileMaxBytes...)
}

func newUploader(uid interface{}, fileMaxBytes ...uint64) *uploadClient.ChunkUpload {
	var fileMaxSize uint64
	if len(fileMaxBytes) > 0 {
		fileMaxSize = fileMaxBytes[0]
	}
	return chunkUpload.Clone().SetUID(uid).SetFileMaxBytes(fileMaxSize).SetFileNameGenerator(MergedFileNameGenerator(uid))
}

func initChunkUploader() { // 初始化后台实例，主要用于定时清理过期文件
	chunkUpload = uploadClient.NewChunkUpload(ChunkTempDir, MergeSaveDir, TempLifetime)
	//echo.Dump(chunkUpload)
	go chunkUpload.StartGC(GCInterval) // 定时清理过期文件
}
