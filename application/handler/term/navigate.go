package term

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	DisplayOnMenu: true,
	Name:          `SSH管理`,
	Action:        `term`,
	Icon:          `terminal`,
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
			Name:          `SSH操作`,
			Action:        `client`,
			Icon:          ``,
		},
		{
			DisplayOnMenu: false,
			Name:          `SFTP操作`,
			Action:        `sftp`,
			Icon:          ``,
		},
	},
}
