package formbuilder

import (
	"embed"
	"html/template"

	"github.com/coscms/forms/common"
	"github.com/webx-top/echo/middleware/tplfunc"
)

//go:embed templates
var templateFS embed.FS

func init() {
	common.FileSystem.Register(templateFS)
	common.TplFuncs = func() template.FuncMap {
		return tplfunc.TplFuncMap
	}
}
