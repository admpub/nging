package prepare

import (
	"fmt"

	"github.com/admpub/nging/v4/application/library/config"
	uploadLibrary "github.com/admpub/nging/v4/application/library/upload"
	modelFile "github.com/admpub/nging/v4/application/model/file"
	uploadChunk "github.com/admpub/nging/v4/application/registry/upload/chunk"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

func NewClientWithModel(fileM *modelFile.File, clientName string, result *uploadClient.Result) uploadClient.Client {
	return NewClientWithResult(fileM.Context(), fileM.OwnerType, fileM.OwnerId, clientName, result)
}

func NewClient(ctx echo.Context, ownerType string, ownerID uint64, clientName string, fileType string) uploadClient.Client {
	result := &uploadClient.Result{
		FileType: uploadClient.FileType(fileType),
	}
	return NewClientWithResult(ctx, ownerType, ownerID, clientName, result)
}

func NewClientWithResult(ctx echo.Context, ownerType string, ownerID uint64, clientName string, result *uploadClient.Result) uploadClient.Client {
	client := uploadClient.Get(clientName)
	client.Init(ctx, result)
	cu := uploadChunk.ChunkUploader()
	cu.UID = fmt.Sprintf(`%s/%d`, ownerType, ownerID)
	client.SetChunkUpload(&cu)
	uploadCfg := uploadLibrary.Get()
	maxSize := uploadCfg.MaxSizeBytes(result.FileType.String())
	if maxSize <= 0 {
		maxSize = config.DefaultConfig.GetMaxRequestBodySize()
	}
	client.SetUploadMaxSize(int64(maxSize))
	return client
}
