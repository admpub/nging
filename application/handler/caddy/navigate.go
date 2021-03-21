package caddy

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	DisplayOnMenu: true,
	Name:          `网站管理`,
	Action:        `caddy`,
	Icon:          `sitemap`,
	Children: &navigate.List{
		{
			DisplayOnMenu: false,
			Name:          `Caddy日志`,
			Action:        `log_show`,
		},
		{
			DisplayOnMenu: true,
			Name:          `网站列表`,
			Action:        `vhost`,
		},
		{
			DisplayOnMenu: true,
			Name:          `添加网站`,
			Action:        `vhost_add`,
			Icon:          `plus`,
		},
		{
			DisplayOnMenu: false,
			Name:          `重启Caddy`,
			Action:        `restart`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `停止Caddy`,
			Action:        `stop`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `查看网站日志`,
			Action:        `vhost_log`,
		},
		{
			DisplayOnMenu: false,
			Name:          `查看网站动态`,
			Action:        `log`,
		},
		{
			DisplayOnMenu: false,
			Name:          `配置表单`,
			Action:        `addon_form`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改网站`,
			Action:        `vhost_edit`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除网站`,
			Action:        `vhost_delete`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `管理网站文件`,
			Action:        `vhost_file`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `生成Caddyfile`,
			Action:        `vhost_build`,
			Icon:          ``,
		},

		{
			DisplayOnMenu: true,
			Name:          `分组管理`,
			Action:        `group`,
		},
		{
			DisplayOnMenu: true,
			Name:          `添加分组`,
			Action:        `group_add`,
			Icon:          `plus`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改分组`,
			Action:        `group_edit`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除分组`,
			Action:        `group_delete`,
			Icon:          ``,
		},
	},
}
