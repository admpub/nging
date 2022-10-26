package standard

import (
	"html/template"

	"github.com/webx-top/echo/param"
)

type RenderData struct {
	Func   template.FuncMap
	Data   interface{}
	Stored param.MapReadonly
}
