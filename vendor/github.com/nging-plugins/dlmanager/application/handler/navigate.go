package handler

import "github.com/admpub/nging/v4/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `离线下载`,
	Action:  `download`,
	Icon:    `download`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `下载管理`,
			Action:  `index.html`,
		},
		{
			Display: true,
			Name:    `文件管理`,
			Action:  `file`,
		},
		{
			Display: false,
			Name:    `总进度信息`,
			Action:  `progress.json`,
		},
		{
			Display: false,
			Name:    `添加任务`,
			Action:  `add_task`,
		},
		{
			Display: false,
			Name:    `删除任务`,
			Action:  `remove_task`,
		},
		{
			Display: false,
			Name:    `启动任务`,
			Action:  `start_task`,
		},
		{
			Display: false,
			Name:    `停止任务`,
			Action:  `stop_task`,
		},
		{
			Display: false,
			Name:    `启动所有任务`,
			Action:  `start_all_task`,
		},
		{
			Display: false,
			Name:    `停止所有任务`,
			Action:  `stop_all_task`,
		},
		{
			Display: false,
			Name:    `单个文件进度信息`,
			Action:  `progress/*`,
		},
	},
}
