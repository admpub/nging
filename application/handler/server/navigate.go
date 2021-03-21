package server

import "github.com/admpub/nging/application/registry/navigate"

var LeftNavigate = &navigate.Item{
	DisplayOnMenu: true,
	Name:          `服务器`,
	Action:        `server`,
	Icon:          `desktop`,
	Children: &navigate.List{
		{
			DisplayOnMenu: true,
			Name:          `服务器信息`,
			Action:        `sysinfo`,
		},
		{
			DisplayOnMenu: true,
			Name:          `网络端口`,
			Action:        `netstat`,
		},
		{
			DisplayOnMenu: true,
			Name:          `执行命令`,
			Action:        `cmd`,
		},
		//指令集
		{
			DisplayOnMenu: true,
			Name:          `指令集`,
			Action:        `command`,
		},
		{
			DisplayOnMenu: false,
			Name:          `添加指令`,
			Action:        `command_add`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改指令`,
			Action:        `command_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除指令`,
			Action:        `command_delete`,
		},
		{
			DisplayOnMenu: true,
			Name:          `服务管理`,
			Action:        `service`,
		},
		{
			DisplayOnMenu: true,
			Name:          `hosts文件`,
			Action:        `hosts`,
		},
		{
			DisplayOnMenu: false,
			Name:          `查看Nging日志`,
			Action:        `log`,
		},
		{
			DisplayOnMenu: false,
			Name:          `查看进程详情`,
			Action:        `process/:pid`,
		},
		{
			DisplayOnMenu: false,
			Name:          `杀死进程`,
			Action:        `procskill/:pid`,
		},
		{
			DisplayOnMenu: false,
			Name:          `命令对话`,
			Action:        `cmdSend/*`,
		},
		{
			DisplayOnMenu: false,
			Name:          `发送命令`,
			Action:        `cmdSendWS`,
		},
		{
			DisplayOnMenu: true,
			Name:          `进程值守`,
			Action:        `daemon_index`,
		},
		{
			DisplayOnMenu: false,
			Name:          `进程值守日志`,
			Action:        `daemon_log`,
		},
		{
			DisplayOnMenu: false,
			Name:          `添加值守配置`,
			Action:        `daemon_add`,
		},
		{
			DisplayOnMenu: false,
			Name:          `修改值守配置`,
			Action:        `daemon_edit`,
		},
		{
			DisplayOnMenu: false,
			Name:          `删除值守配置`,
			Action:        `daemon_delete`,
		},
		{
			DisplayOnMenu: false,
			Name:          `重启值守`,
			Action:        `daemon_restart`,
		},
	},
}
