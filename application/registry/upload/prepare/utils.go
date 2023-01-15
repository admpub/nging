package prepare

import (
	"fmt"
	"path"

	uploadLibrary "github.com/admpub/nging/v5/application/library/upload"
	modelFile "github.com/admpub/nging/v5/application/model/file"
	uploadChunk "github.com/admpub/nging/v5/application/registry/upload/chunk"
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
	cu := uploadChunk.NewUploader(fmt.Sprintf(`%s/%d`, ownerType, ownerID))
	client.SetChunkUpload(cu)
	uploadCfg := uploadLibrary.Get()
	client.SetReadBeforeHook(func(result *uploadClient.Result) error {
		extension := path.Ext(result.FileName)
		result.FileType = uploadClient.FileType(uploadCfg.DetectType(extension))
		return nil
	})
	return client
}
