package handler

import "github.com/admpub/nging/v4/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `内网穿透`,
	Action:  `frp`,
	Icon:    `road`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `服务端配置`,
			Action:  `server_index`,
		},
		{
			Display: false,
			Name:    `查看服务端日志`,
			Action:  `server_log`,
		},
		{
			Display: false,
			Name:    `添加服务端配置`,
			Action:  `server_add`,
		},
		{
			Display: false,
			Name:    `修改服务端配置`,
			Action:  `server_edit`,
		},
		{
			Display: false,
			Name:    `删除服务端配置`,
			Action:  `server_delete`,
		},

		{
			Display: true,
			Name:    `账号管理`,
			Action:  `account`,
		},
		{
			Display: false,
			Name:    `添加FRP账号`,
			Action:  `account_add`,
		},
		{
			Display: false,
			Name:    `修改FRP账号`,
			Action:  `account_edit`,
		},
		{
			Display: false,
			Name:    `删除FRP账号`,
			Action:  `account_delete`,
		},

		{
			Display: true,
			Name:    `客户端配置`,
			Action:  `client_index`,
		},
		{
			Display: false,
			Name:    `查看客户端日志`,
			Action:  `client_log`,
		},
		{
			Display: false,
			Name:    `添加客户端配置`,
			Action:  `client_add`,
		},
		{
			Display: false,
			Name:    `修改客户端配置`,
			Action:  `client_edit`,
		},
		{
			Display: false,
			Name:    `删除客户端配置`,
			Action:  `client_delete`,
		},
		{
			Display: true,
			Name:    `分组管理`,
			Action:  `group_index`,
		},
		{
			Display: false,
			Name:    `添加分组`,
			Action:  `group_add`,
		},
		{
			Display: false,
			Name:    `修改分组`,
			Action:  `group_edit`,
		},
		{
			Display: false,
			Name:    `删除分组`,
			Action:  `group_delete`,
		},
		{
			Display: false,
			Name:    `统计图表`,
			Action:  `dashboard`,
			Target:  `_blank`,
		},
		{
			Display: false,
			Name:    `重启服务端`,
			Action:  `server_restart`,
		},
		{
			Display: false,
			Name:    `关闭服务端`,
			Action:  `server_stop`,
		},
		{
			Display: false,
			Name:    `重启客户端`,
			Action:  `client_restart`,
		},
		{
			Display: false,
			Name:    `关闭客户端`,
			Action:  `client_stop`,
		},
		{
			Display: false,
			Name:    `配置表单`,
			Action:  `addon_form`,
		},
	},
}
