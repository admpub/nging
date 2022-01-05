package prepare

import (
	"fmt"

	"github.com/admpub/nging/v4/application/library/config"
	modelFile "github.com/admpub/nging/v4/application/model/file"
	uploadChunk "github.com/admpub/nging/v4/application/registry/upload/chunk"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/echo"
)

func NewClientWithModel(fileM *modelFile.File, clientName string, result *uploadClient.Result) uploadClient.Client {
	return NewClient(fileM.Context(), fileM.OwnerType, fileM.OwnerId, clientName, result)
}

func NewClient(ctx echo.Context, ownerType string, ownerID uint64, clientName string, results ...*uploadClient.Result) uploadClient.Client {
	var result *uploadClient.Result
	if len(results) > 0 && results[0] != nil {
		result = results[0]
	} else {
		result = &uploadClient.Result{}
	}
	client := uploadClient.Get(clientName)
	client.Init(ctx, result)
	cu := uploadChunk.ChunkUploader()
	cu.UID = fmt.Sprintf(`%s/%d`, ownerType, ownerID)
	client.SetChunkUpload(&cu)
	client.SetUploadMaxSize(int64(config.DefaultConfig.GetMaxRequestBodySize()))
	return client
}
