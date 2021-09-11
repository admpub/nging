package frp

import (
	"math/rand"
	"os"
	"sort"
	"time"

	"github.com/admpub/frp/pkg/config"
	plugin "github.com/admpub/frp/pkg/plugin/server"
	"github.com/admpub/nging/v3/application/library/hook"
	syncOnce "github.com/admpub/once"
	"github.com/fatedier/golib/crypto"
)

type PluginGetter func(*config.ServerCommonConf) plugin.HTTPPluginOptions

type Plugin struct {
	Name   string
	Title  string
	getter PluginGetter
}

func (p *Plugin) Getter() PluginGetter {
	return p.getter
}

var (
	once      syncOnce.Once
	kcpDoneCh chan struct{}
	// 全局插件
	serverPlugins = map[string]*Plugin{}

	Hook = hook.New()
)

func onceInit() {
	crypto.DefaultSalt = os.Getenv(`FRP_CRYPTO_SALT`)
	if len(crypto.DefaultSalt) == 0 {
		crypto.DefaultSalt = `frp`
	}
	rand.Seed(time.Now().UnixNano())
}

func ServerPluginRegister(name string, title string, plug PluginGetter) {
	serverPlugins[name] = &Plugin{Name: name, Title: title, getter: plug}
}

func ServerPluginUnregister(names ...string) {
	for _, name := range names {
		delete(serverPlugins, name)
	}
}

func ServerPluginGet(name string) (plug *Plugin) {
	plug, _ = serverPlugins[name]
	return
}

func ServerPluginExists(name string) bool {
	_, ok := serverPlugins[name]
	return ok
}

func ServerPluginSlice() []*Plugin {
	var names []string
	for name := range serverPlugins {
		names = append(names, name)
	}
	sort.Strings(names)
	res := make([]*Plugin, len(names))
	for i, name := range names {
		res[i] = ServerPluginGet(name)
	}
	return res
}
