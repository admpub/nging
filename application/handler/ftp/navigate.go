package ftp

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	DisplayOnMenu: true,
	Name:          `FTP账号`,
	Action:        `ftp`,
	Icon:          `users`,
	Children: &navigate.List{
		{
			DisplayOnMenu: true,
			Name:          `账号管理`,
			Action:        `account`,
		},
		{
			DisplayOnMenu: true,
			Name:          `添加账号`,
			Action:        `account_add`,
			Icon:          `plus`,
		},
		{
			DisplayOnMenu: true,
			Name:          `用户组`,
			Action:        `group`,
		},
		{
			DisplayOnMenu: true,
			Name:          `添加用户组`,
			Action:        `group_add`,
			Icon:          `plus`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改账号`,
			Action:        `account_edit`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除账号`,
			Action:        `account_delete`,
			Icon:          ``,
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
		{
			DisplayOnMenu: false,
			Name:          `重启FTP服务`,
			Action:        `restart`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `关闭FTP服务`,
			Action:        `stop`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `查看FTP动态`,
			Action:        `log`,
			Icon:          ``,
		},
	},
}
