package download

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	DisplayOnMenu: true,
	Name:          `离线下载`,
	Action:        `download`,
	Icon:          `download`,
	Children: &navigate.List{
		{
			DisplayOnMenu: true,
			Name:          `下载管理`,
			Action:        `index.html`,
		},
		{
			DisplayOnMenu: true,
			Name:          `文件管理`,
			Action:        `file`,
		},
		{
			DisplayOnMenu: false,
			Name:          `总进度信息`,
			Action:        `progress.json`,
		},
		{
			DisplayOnMenu: false,
			Name:          `添加任务`,
			Action:        `add_task`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除任务`,
			Action:        `remove_task`,
		},
		{
			DisplayOnMenu: false,
			Name:          `启动任务`,
			Action:        `start_task`,
		},
		{
			DisplayOnMenu: false,
			Name:          `停止任务`,
			Action:        `stop_task`,
		},
		{
			DisplayOnMenu: false,
			Name:          `启动所有任务`,
			Action:        `start_all_task`,
		},
		{
			DisplayOnMenu: false,
			Name:          `停止所有任务`,
			Action:        `stop_all_task`,
		},
		{
			DisplayOnMenu: false,
			Name:          `单个文件进度信息`,
			Action:        `progress/*`,
		},
	},
}
