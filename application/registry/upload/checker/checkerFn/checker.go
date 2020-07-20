package checkerFn

import (
	"github.com/admpub/nging/application/registry/upload/table"
	"github.com/webx-top/echo"
)

// Checker 验证并生成子文件夹名称和文件名称
type Checker func(echo.Context, table.TableInfoStorer) (subdir string, name string, err error)
