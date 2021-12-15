package handler

import "github.com/admpub/nging/v4/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `FTP账号`,
	Action:  `ftp`,
	Icon:    `users`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `账号管理`,
			Action:  `account`,
		},
		{
			Display: true,
			Name:    `添加账号`,
			Action:  `account_add`,
			Icon:    `plus`,
		},
		{
			Display: true,
			Name:    `用户组`,
			Action:  `group`,
		},
		{
			Display: true,
			Name:    `添加用户组`,
			Action:  `group_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改账号`,
			Action:  `account_edit`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `删除账号`,
			Action:  `account_delete`,
			Icon:    ``,
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
		{
			Display: false,
			Name:    `重启FTP服务`,
			Action:  `restart`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `关闭FTP服务`,
			Action:  `stop`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `查看FTP动态`,
			Action:  `log`,
			Icon:    ``,
		},
	},
}
