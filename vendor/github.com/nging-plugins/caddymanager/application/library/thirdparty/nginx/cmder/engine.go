package cmder

import (
	"github.com/admpub/log"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/library/common"
	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/engine"
	nginxConfigPkg "github.com/nging-plugins/caddymanager/application/library/thirdparty/nginx/config"
)

func init() {
	engine.Engines.Add(nginxConfigPkg.Name, `Nginx`, echo.KVOptX(newEngine()))
}

func newEngine() engine.Enginer {
	return &Engine{}
}

type Engine struct{}

func (b *Engine) Name() string {
	return nginxConfigPkg.Name
}

func (b *Engine) ListConfig(ctx echo.Context) ([]engine.Configer, error) {
	ident := ctx.Internal().String(`serverIdent`)
	m := dbschema.NewNgingVhostServer(ctx)
	cond := db.NewCompounds()
	cond.AddKV(`disabled`, common.BoolN)
	cond.AddKV(`engine`, nginxConfigPkg.Name)
	if len(ident) > 0 {
		cond.AddKV(`ident`, ident)
	}
	_, err := m.ListByOffset(nil, nil, 0, -1, cond.And())
	if err != nil {
		log.Error(err)
		return nil, err
	}
	rows := m.Objects()
	result := make([]engine.Configer, len(rows))
	for idx, row := range rows {
		result[idx] = b.BuildConfig(ctx, row)
	}
	return result, err
}

func (b *Engine) BuildConfig(ctx echo.Context, m *dbschema.NgingVhostServer) engine.Configer {
	cfg := nginxConfigPkg.New()
	cfg.CopyFrom(m)
	return cfg
}

func (b *Engine) ReloadServer(ctx echo.Context, cfg engine.Configer) error {
	return cfg.(*nginxConfigPkg.Config).Reload(ctx)
}

func (b *Engine) DefaultConfigDir() string {
	return nginxConfigPkg.DefaultConfigDir()
}
