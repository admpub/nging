package handler

import "github.com/admpub/nging/v5/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `数据库`,
	Action:  `db`,
	Icon:    `table`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `数据库账号`,
			Action:  `account`,
		},
		{
			Display: true,
			Name:    `添加账号`,
			Action:  `account_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    `修改账号`,
			Action:  `account_edit`,
		},
		{
			Display: false,
			Name:    `删除账号`,
			Action:  `account_delete`,
		},
		{
			Display: true,
			Name:    `连接数据库`,
			Action:  ``,
		},
		{
			Display: true,
			Name:    `表结构同步`,
			Action:  `schema_sync`,
		},
		{
			Display: false,
			Name:    `新增同步方案`,
			Action:  `schema_sync_add`,
		},
		{
			Display: false,
			Name:    `修改同步方案`,
			Action:  `schema_sync_edit`,
		},
		{
			Display: false,
			Name:    `删除同步方案`,
			Action:  `schema_sync_delete`,
		},
		{
			Display: false,
			Name:    `预览表结构差异`,
			Action:  `schema_sync_preview`,
		},
		{
			Display: false,
			Name:    `执行表结构同步`,
			Action:  `schema_sync_run`,
		},
		{
			Display: false,
			Name:    `表结构同步日志列表`,
			Action:  `schema_sync_log/:id`,
		},
		{
			Display: false,
			Name:    `表结构同步日志详情`,
			Action:  `schema_sync_log_view/:id`,
		},
		{
			Display: false,
			Name:    `删除表结构同步日志`,
			Action:  `schema_sync_log_delete`,
		},
	},
}
