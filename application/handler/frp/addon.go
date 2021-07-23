package frp

import (
	"net/http"
	"regexp"

	"github.com/admpub/nging/v3/application/handler/caddy"
	"github.com/admpub/nging/v3/application/library/common"
	"github.com/webx-top/echo"
)

var regexNumEnd = regexp.MustCompile(`_[\d]+$`)

type Section struct {
	Section string
	Addon   string
}

func AddonForm(ctx echo.Context) error {
	addon := ctx.Query(`addon`)
	if len(addon) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, ctx.T("参数 addon 的值不能为空"))
	}
	if !caddy.ValidAddonName(addon) {
		return echo.NewHTTPError(http.StatusBadRequest, ctx.T("参数 addon 的值包含非法字符"))
	}
	section := ctx.Query(`section`, addon)
	setAddonFunc(ctx)
	return ctx.Render(`frp/client/form/`+addon, section)
}

func setAddonFunc(ctx echo.Context) {
	prefix := `extra`
	formKey := func(key string, keys ...string) string {
		key = prefix + `[` + key + `]`
		for _, k := range keys {
			key += `[` + k + `]`
		}
		return key
	}
	ctx.SetFunc(`Val`, func(key string, keys ...string) string {
		return ctx.Form(formKey(key, keys...))
	})
	ctx.SetFunc(`Vals`, func(key string, keys ...string) []string {
		return ctx.FormValues(formKey(key, keys...))
	})
	ctx.SetFunc(`Key`, formKey)
	ipv4, _ := common.GetLocalIP()
	ctx.Set(`localIP`, ipv4)
}
