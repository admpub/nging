package initialize

import (
	"github.com/admpub/nging/application/library/common"
	modelFile "github.com/admpub/nging/application/model/file"
)

func init() {
	modelFile.Enable = true
	common.OnUpdateOwnerFilePath = OnUpdateOwnerFilePath
	common.OnRemoveOwnerFile = OnRemoveOwnerFile
}
