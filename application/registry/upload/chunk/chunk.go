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
)

func ChunkUploader() uploadClient.ChunkUpload {
	chunkUploadInitOnce.Do(func() {
		chunkUpload = &uploadClient.ChunkUpload{
			TempDir:      filepath.Join(os.TempDir(), `nging/chunk_temp`),
			SaveDir:      filepath.Join(os.TempDir(), `nging/chunk_merged`),
			TempLifetime: 24 * time.Hour,
		}
		echo.Dump(chunkUpload)
		go chunkUpload.StartGC(48 * time.Hour)
	})
	return *chunkUpload
}
