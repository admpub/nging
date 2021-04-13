package cloud

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `云服务`,
	Action:  `cloud`,
	Icon:    `cloud`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `云存储账号`,
			Action:  `storage`,
		},
		{
			Display: true,
			Name:    `添加账号`,
			Action:  `storage_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改账号`,
			Action:  `storage_edit`,
		},
		{
			Display: false,
			Name:    `删除账号`,
			Action:  `storage_delete`,
		},
		{
			Display: false,
			Name:    `云存储文件管理`,
			Action:  `storage_file`,
		},
		{
			Display: true,
			Name:    `文件备份`,
			Action:  `backup`,
		},
		{
			Display: false,
			Name:    `添加备份配置`,
			Action:  `backup_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改备份配置`,
			Action:  `backup_edit`,
		},
		{
			Display: false,
			Name:    `删除备份配置`,
			Action:  `backup_delete`,
		},
	},
}
