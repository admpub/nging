package handler

import "github.com/admpub/nging/v4/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `网站管理`,
	Action:  `caddy`,
	Icon:    `sitemap`,
	Children: &navigate.List{
		{
			Display: false,
			Name:    `Caddy日志`,
			Action:  `log_show`,
		},
		{
			Display: true,
			Name:    `网站列表`,
			Action:  `vhost`,
		},
		{
			Display: true,
			Name:    `添加网站`,
			Action:  `vhost_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `重启Caddy`,
			Action:  `restart`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `停止Caddy`,
			Action:  `stop`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `查看网站日志`,
			Action:  `vhost_log`,
		},
		{
			Display: false,
			Name:    `查看网站动态`,
			Action:  `log`,
		},
		{
			Display: false,
			Name:    `配置表单`,
			Action:  `addon_form`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `修改网站`,
			Action:  `vhost_edit`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `删除网站`,
			Action:  `vhost_delete`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `管理网站文件`,
			Action:  `vhost_file`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `生成Caddyfile`,
			Action:  `vhost_build`,
			Icon:    ``,
		},

		{
			Display: true,
			Name:    `分组管理`,
			Action:  `group`,
		},
		{
			Display: true,
			Name:    `添加分组`,
			Action:  `group_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改分组`,
			Action:  `group_edit`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `删除分组`,
			Action:  `group_delete`,
			Icon:    ``,
		},
	},
}
