package database

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	DisplayOnMenu: true,
	Name:          `数据库`,
	Action:        `db`,
	Icon:          `table`,
	Children: &navigate.List{
		{
			DisplayOnMenu: true,
			Name:          `数据库账号`,
			Action:        `account`,
		},
		{
			DisplayOnMenu: true,
			Name:          `添加账号`,
			Action:        `account_add`,
			Icon:          `plus`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改账号`,
			Action:        `account_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除账号`,
			Action:        `account_delete`,
		},
		{
			DisplayOnMenu: true,
			Name:          `连接数据库`,
			Action:        ``,
		},
		{
			DisplayOnMenu: true,
			Name:          `表结构同步`,
			Action:        `schema_sync`,
		},
		{
			DisplayOnMenu: false,
			Name:          `新增同步方案`,
			Action:        `schema_sync_add`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改同步方案`,
			Action:        `schema_sync_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除同步方案`,
			Action:        `schema_sync_delete`,
		},
		{
			DisplayOnMenu: false,
			Name:          `预览表结构差异`,
			Action:        `schema_sync_preview`,
		},
		{
			DisplayOnMenu: false,
			Name:          `执行表结构同步`,
			Action:        `schema_sync_run`,
		},
		{
			DisplayOnMenu: false,
			Name:          `表结构同步日志列表`,
			Action:        `schema_sync_log/:id`,
		},
		{
			DisplayOnMenu: false,
			Name:          `表结构同步日志详情`,
			Action:        `schema_sync_log_view/:id`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除表结构同步日志`,
			Action:        `schema_sync_log_delete`,
		},
	},
}
