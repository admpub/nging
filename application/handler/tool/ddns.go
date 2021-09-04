package tool

import (
	"context"

	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/ddnsmanager"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/boot"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/config"
	_ "github.com/admpub/nging/v3/application/library/ddnsmanager/providerall"
	"github.com/webx-top/echo"
)

func DdnsSettings(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		if err = ctx.MustBindAndValidate(boot.Config); err != nil {
			goto END
		}
		boot.Reset(context.Background())
		return ctx.Redirect(`/tool/ddns`)
	}

END:
	ctx.Set(`config`, boot.Config)
	ctx.Set(`ttlList`, config.TTLs.Slice())
	ctx.Set(`providers`, ddnsmanager.All())
	ctx.Set(`title`, `DDNS`)
	return ctx.Render(`tool/ddns`, handler.Err(ctx, err))
}
