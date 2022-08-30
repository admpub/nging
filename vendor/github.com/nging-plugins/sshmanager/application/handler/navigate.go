package handler

import "github.com/admpub/nging/v4/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `SSH管理`,
	Action:  `term`,
	Icon:    `terminal`,
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
			Name:    `SSH操作`,
			Action:  `client`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `SFTP操作`,
			Action:  `sftp`,
			Icon:    ``,
		},
	},
}
