package handler

import "github.com/admpub/nging/v4/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    `服务器`,
	Action:  `server`,
	Icon:    `desktop`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    `服务器信息`,
			Action:  `sysinfo`,
		},
		{
			Display: false,
			Name:    `服务器进程`,
			Action:  `processes`,
		},
		{
			Display: true,
			Name:    `网络端口`,
			Action:  `netstat`,
		},
		{
			Display: true,
			Name:    `执行命令`,
			Action:  `cmd`,
		},
		//指令集
		{
			Display: true,
			Name:    `指令集`,
			Action:  `command`,
		},
		{
			Display: false,
			Name:    `添加指令`,
			Action:  `command_add`,
		},
		{
			Display: false,
			Name:    `修改指令`,
			Action:  `command_edit`,
		},
		{
			Display: false,
			Name:    `删除指令`,
			Action:  `command_delete`,
		},
		{
			Display: true,
			Name:    `服务管理`,
			Action:  `service`,
		},
		{
			Display: false,
			Name:    `查看服务日志`,
			Action:  `log/:category`,
		},
		{
			Display: true,
			Name:    `hosts文件`,
			Action:  `hosts`,
		},
		{
			Display: false,
			Name:    `查看Nging日志`,
			Action:  `log`,
		},
		{
			Display: false,
			Name:    `查看进程详情`,
			Action:  `process/:pid`,
		},
		{
			Display: false,
			Name:    `杀死进程`,
			Action:  `procskill/:pid`,
		},
		{
			Display: false,
			Name:    `命令对话`,
			Action:  `cmdSend/*`,
		},
		{
			Display: false,
			Name:    `发送命令`,
			Action:  `cmdSendWS`,
		},
		{
			Display: true,
			Name:    `进程值守`,
			Action:  `daemon_index`,
		},
		{
			Display: false,
			Name:    `进程值守日志`,
			Action:  `daemon_log`,
		},
		{
			Display: false,
			Name:    `添加值守配置`,
			Action:  `daemon_add`,
		},
		{
			Display: false,
			Name:    `修改值守配置`,
			Action:  `daemon_edit`,
		},
		{
			Display: false,
			Name:    `删除值守配置`,
			Action:  `daemon_delete`,
		},
		{
			Display: false,
			Name:    `重启值守`,
			Action:  `daemon_restart`,
		},
	},
}
