package cmder

import (
	"github.com/webx-top/echo"

	"github.com/nging-plugins/caddymanager/application/dbschema"
	"github.com/nging-plugins/caddymanager/application/library/caddy"
	"github.com/nging-plugins/caddymanager/application/library/engine"
)

func init() {
	engine.Engines.Add(`default`, `默认(内置Caddy1)`, echo.KVOptX(newEngine()))
}

func newEngine() engine.Enginer {
	return &Engine{}
}

type Engine struct{}

func (b *Engine) Name() string {
	return Name
}

func (b *Engine) ListConfig(ctx echo.Context) ([]engine.Configer, error) {
	return []engine.Configer{GetCaddyConfig()}, nil
}

func (b *Engine) BuildConfig(ctx echo.Context, m *dbschema.NgingVhostServer) engine.Configer {
	return nil
}

func (b *Engine) ReloadServer(ctx echo.Context, cfg engine.Configer) error {
	return GetCaddyCmd().ReloadServer()
}

func (b *Engine) DefaultConfigDir() string {
	return caddy.DefaultConfigDir()
}
