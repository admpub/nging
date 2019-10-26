package initialize

import (
	"github.com/admpub/nging/application/library/common"
	modelFile "github.com/admpub/nging/application/model/file"
	"github.com/admpub/nging/application/model/file/helper"
)

func init() {
	modelFile.Enable = true
	common.OnUpdateOwnerFilePath = helper.OnUpdateOwnerFilePath
	common.OnRemoveOwnerFile = helper.OnRemoveOwnerFile
}
