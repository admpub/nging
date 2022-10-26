package model

import (
	"context"

	"github.com/admpub/nging/v5/application/model/file/storer"
	"github.com/webx-top/echo/defaults"
	"github.com/webx-top/echo/param"
)

const StorageAccountIDKey = `storerID`

func GetCloudStorage(ctx context.Context) (*CloudStorage, error) {
	var cloudAccountID string
	eCtx := defaults.MustGetContext(ctx)
	cloudAccountID = eCtx.Internal().String(StorageAccountIDKey)
	m := NewCloudStorage(eCtx)
	if len(cloudAccountID) == 0 {
		cloudAccountID = param.AsString(ctx.Value(StorageAccountIDKey))
	}
	if len(cloudAccountID) == 0 || cloudAccountID == `0` {
		storerConfig, ok := storer.GetOk()
		if ok {
			cloudAccountID = storerConfig.ID
		}
	}
	err := m.Get(nil, `id`, cloudAccountID)
	return m, err
}
