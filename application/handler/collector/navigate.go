package collector

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `数据采集`,
	Action:  `collector`,
	Icon:    `truck`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `规则管理`,
			Action:  `rule`,
		},
		{
			Display: true,
			Name:    `新建规则`,
			Action:  `rule_add`,
			Icon:    `plus`,
		},
		{
			Display: true,
			Name:    `历史记录`,
			Action:  `history`,
		},
		{
			Display: false,
			Name:    `查看历史内容`,
			Action:  `history_view`,
		},
		{
			Display: false,
			Name:    `删除历史记录`,
			Action:  `history_delete`,
		},
		{
			Display: true,
			Name:    `数据导出`,
			Action:  `export`,
		},
		{
			Display: false,
			Name:    `添加导出规则`,
			Action:  `export_add`,
		},
		{
			Display: false,
			Name:    `修改导出规则`,
			Action:  `export_edit`,
		},
		{
			Display: false,
			Name:    `删除导出规则`,
			Action:  `export_delete`,
		},
		{
			Display: false,
			Name:    `导出日志管理`,
			Action:  `export_log`,
		},
		{
			Display: false,
			Name:    `查看导出日志`,
			Action:  `export_log_view/:id`,
		},
		{
			Display: false,
			Name:    `删除导出日志`,
			Action:  `export_log_delete`,
		},
		{
			Display: false,
			Name:    `更改导出日志状态`,
			Action:  `export_edit_status`,
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
			Name:    `修改规则`,
			Action:  `rule_edit`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `删除规则`,
			Action:  `rule_delete`,
			Icon:    ``,
		},
		{
			Display: false,
			Name:    `规则测试`,
			Action:  `rule_collect`,
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
			Name:    `测试正则表达式`,
			Action:  `regexp_test`,
			Icon:    ``,
		},
	},
}
