package chunk

import (
	"os"
	"path/filepath"
	"time"

	syncOnce "github.com/admpub/once"
	uploadClient "github.com/webx-top/client/upload"
)

var (
	chunkUploadInitOnce syncOnce.Once
	chunkUpload         *uploadClient.ChunkUpload

	// ChunkTempDir 保存分片文件的临时文件夹
	ChunkTempDir = filepath.Join(os.TempDir(), `nging/chunk_temp`)

	// MergeSaveDir 分片文件合并后保存的临时文件夹
	// 避免和最终的文件保存位置相同，上传后一般需要将合并后的文件转移(os.Rename)到最终保存位置
	MergeSaveDir = filepath.Join(os.TempDir(), `nging/chunk_merged`)

	// TempLifetime 分片临时文件过期时间
	TempLifetime = 24 * time.Hour

	// GCInterval 垃圾清理间隔时间
	GCInterval = 48 * time.Hour
)

func ChunkUploader() uploadClient.ChunkUpload {
	chunkUploadInitOnce.Do(initChunkUploader)
	return *chunkUpload
}

func initChunkUploader() {
	chunkUpload = &uploadClient.ChunkUpload{
		TempDir:      ChunkTempDir,
		SaveDir:      MergeSaveDir,
		TempLifetime: TempLifetime,
	}
	//echo.Dump(chunkUpload)
	go chunkUpload.StartGC(GCInterval)
}
