package task

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	DisplayOnMenu: true,
	Name:          `计划任务`,
	Action:        `task`,
	Icon:          `clock-o`,
	Children: &navigate.List{
		{
			DisplayOnMenu: true,
			Name:          `任务管理`,
			Action:        `index`,
		},
		{
			DisplayOnMenu: true,
			Name:          `新建任务`,
			Action:        `add`,
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
			DisplayOnMenu: true,
			Name:          `测试邮件`,
			Action:        `email_test`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改任务`,
			Action:        `edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除任务`,
			Action:        `delete`,
		},
		{
			DisplayOnMenu: false,
			Name:          `启动任务`,
			Action:        `start`,
		},
		{
			DisplayOnMenu: false,
			Name:          `暂停任务`,
			Action:        `pause`,
		},
		{
			DisplayOnMenu: false,
			Name:          `立即执行任务`,
			Action:        `run`,
		},
		{
			DisplayOnMenu: false,
			Name:          `退出任务`,
			Action:        `exit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `启动历史任务`,
			Action:        `start_history`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改分组`,
			Action:        `group_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除分组`,
			Action:        `group_delete`,
		},
		{
			DisplayOnMenu: false,
			Name:          `日志列表`,
			Action:        `log`,
		},
		{
			DisplayOnMenu: false,
			Name:          `日志详情`,
			Action:        `log_view/:id`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除日志`,
			Action:        `log_delete`,
		},
	},
}
