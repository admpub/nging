package frp

import (
	"sort"

	plugin "github.com/admpub/frp/pkg/plugin/server"
)

var (
	// 全局插件
	serverPlugins = map[string]*Plugin{}
)

type PluginGetter func() plugin.HTTPPluginOptions

type Plugin struct {
	Name   string
	Title  string
	getter PluginGetter
}

func (p *Plugin) Getter() PluginGetter {
	return p.getter
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

func ServerPluginOptions(pluginNames ...string) map[string]plugin.HTTPPluginOptions {
	r := map[string]plugin.HTTPPluginOptions{}
	if len(pluginNames) == 0 {
		for name, plugin := range serverPlugins {
			if plugin != nil && plugin.getter != nil {
				r[name] = plugin.getter()
			}
		}
	} else {
		for _, name := range pluginNames {
			if plugin, ok := serverPlugins[name]; ok && plugin != nil && plugin.getter != nil {
				r[name] = plugin.getter()
			}
		}
	}
	return r
}
