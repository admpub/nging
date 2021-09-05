package tool

import (
	"github.com/admpub/nging/v3/application/handler"
	"github.com/admpub/nging/v3/application/library/ddnsmanager"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/boot"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/config"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/domain/dnsdomain"
	_ "github.com/admpub/nging/v3/application/library/ddnsmanager/providerall"
	"github.com/admpub/nging/v3/application/library/ddnsmanager/utils"
	"github.com/webx-top/echo"
)

func DdnsSettings(ctx echo.Context) error {
	var err error
	if ctx.IsPost() {
		ddnsConfig := config.New()
		if err = ctx.MustBindAndValidate(ddnsConfig); err != nil {
			goto END
		}
		echo.Dump(ddnsConfig)
		// *boot.Config = *ddnsConfig
		// boot.Reset(context.Background())
		return ctx.Redirect(`/tool/ddns`)
	}

END:
	ctx.Set(`config`, boot.Config)
	ctx.Set(`ttlList`, config.TTLs.Slice())
	ctx.Set(`providers`, ddnsmanager.AllProvoderMeta(boot.Config.DNSServices))
	ctx.Set(`title`, `DDNS`)
	ipv4NetInterfaces, ipv6NetInterfaces, _ := utils.GetNetInterface(``)
	ctx.Set(`ipv4NetInterfaces`, ipv4NetInterfaces)
	ctx.Set(`ipv6NetInterfaces`, ipv6NetInterfaces)
	ctx.Set(`tagValueDescs`, dnsdomain.TagValueDescs.Slice())
	return ctx.Render(`tool/ddns`, handler.Err(ctx, err))
}
