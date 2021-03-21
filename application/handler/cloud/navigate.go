package cloud

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	DisplayOnMenu: true,
	Name:          `云服务`,
	Action:        `cloud`,
	Icon:          `cloud`,
	Children: &navigate.List{
		{
			DisplayOnMenu: true,
			Name:          `云存储账号`,
			Action:        `storage`,
		},
		{
			DisplayOnMenu: true,
			Name:          `添加账号`,
			Action:        `storage_add`,
			Icon:          `plus`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改账号`,
			Action:        `storage_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除账号`,
			Action:        `storage_delete`,
		},
		{
			DisplayOnMenu: false,
			Name:          `云存储文件管理`,
			Action:        `storage_file`,
		},
		{
			DisplayOnMenu: true,
			Name:          `文件备份`,
			Action:        `backup`,
		},
		{
			DisplayOnMenu: false,
			Name:          `添加备份配置`,
			Action:        `backup_add`,
			Icon:          `plus`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改备份配置`,
			Action:        `backup_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除备份配置`,
			Action:        `backup_delete`,
		},
	},
}
