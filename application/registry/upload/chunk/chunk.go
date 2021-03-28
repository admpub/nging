package chunk

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

var (
	chunkUploadInitOnce sync.Once
	chunkUpload         *uploadClient.ChunkUpload

	// ChunkTempDir 保存分片文件的临时文件夹
	ChunkTempDir = filepath.Join(os.TempDir(), `nging/chunk_temp`)

	// MergeSaveDir 分片文件合并后保存的临时文件夹
	// 避免和最终的文件保存位置相同，上传后一般需要将合并后的文件转移(os.Rename)到最终保存位置
	MergeSaveDir = filepath.Join(os.TempDir(), `nging/chunk_merged`)
)

func ChunkUploader() uploadClient.ChunkUpload {
	chunkUploadInitOnce.Do(func() {
		chunkUpload = &uploadClient.ChunkUpload{
			TempDir:      ChunkTempDir,
			SaveDir:      MergeSaveDir,
			TempLifetime: 24 * time.Hour,
		}
		echo.Dump(chunkUpload)
		go chunkUpload.StartGC(48 * time.Hour)
	})
	return *chunkUpload
}
