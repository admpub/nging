package frp

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	DisplayOnMenu: true,
	Name:          `内网穿透`,
	Action:        `frp`,
	Icon:          `road`,
	Children: &navigate.List{
		{
			DisplayOnMenu: true,
			Name:          `服务端配置`,
			Action:        `server_index`,
		},
		{
			DisplayOnMenu: false,
			Name:          `查看服务端日志`,
			Action:        `server_log`,
		},
		{
			DisplayOnMenu: false,
			Name:          `添加服务端配置`,
			Action:        `server_add`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改服务端配置`,
			Action:        `server_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除服务端配置`,
			Action:        `server_delete`,
		},
		{
			DisplayOnMenu: true,
			Name:          `客户端配置`,
			Action:        `client_index`,
		},
		{
			DisplayOnMenu: false,
			Name:          `查看客户端日志`,
			Action:        `client_log`,
		},
		{
			DisplayOnMenu: false,
			Name:          `添加客户端配置`,
			Action:        `client_add`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改客户端配置`,
			Action:        `client_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除客户端配置`,
			Action:        `client_delete`,
		},
		{
			DisplayOnMenu: true,
			Name:          `分组管理`,
			Action:        `group_index`,
		},
		{
			DisplayOnMenu: false,
			Name:          `添加分组`,
			Action:        `group_add`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改分组`,
			Action:        `group_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除分组`,
			Action:        `group_delete`,
		},
		{
			DisplayOnMenu: true,
			Name:          `统计图表`,
			Action:        `dashboard`,
			Target:        `_blank`,
		},
		{
			DisplayOnMenu: false,
			Name:          `重启服务端`,
			Action:        `server_restart`,
		},
		{
			DisplayOnMenu: false,
			Name:          `关闭服务端`,
			Action:        `server_stop`,
		},
		{
			DisplayOnMenu: false,
			Name:          `重启客户端`,
			Action:        `client_restart`,
		},
		{
			DisplayOnMenu: false,
			Name:          `关闭客户端`,
			Action:        `client_stop`,
		},
		{
			DisplayOnMenu: false,
			Name:          `配置表单`,
			Action:        `addon_form`,
		},
	},
}
