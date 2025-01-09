package task

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`计划任务`),
	Action:  `task`,
	Icon:    `clock-o`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    echo.T(`任务管理`),
			Action:  `index`,
		},
		{
			Display: true,
			Name:    echo.T(`新建任务`),
			Action:  `add`,
			Icon:    `plus`,
		},
		{
			Display: true,
			Name:    echo.T(`分组管理`),
			Action:  `group`,
		},
		{
			Display: true,
			Name:    echo.T(`添加分组`),
			Action:  `group_add`,
			Icon:    `plus`,
		},
		{
			Display: true,
			Name:    echo.T(`测试邮件`),
			Action:  `email_test`,
		},
		{
			Display: false,
			Name:    echo.T(`修改任务`),
			Action:  `edit`,
		},
		{
			Display: false,
			Name:    echo.T(`删除任务`),
			Action:  `delete`,
		},
		{
			Display: false,
			Name:    echo.T(`启动任务`),
			Action:  `start`,
		},
		{
			Display: false,
			Name:    echo.T(`暂停任务`),
			Action:  `pause`,
		},
		{
			Display: false,
			Name:    echo.T(`立即执行任务`),
			Action:  `run`,
		},
		{
			Display: false,
			Name:    echo.T(`退出任务`),
			Action:  `exit`,
		},
		{
			Display: false,
			Name:    echo.T(`启动历史任务`),
			Action:  `start_history`,
		},
		{
			Display: false,
			Name:    echo.T(`修改分组`),
			Action:  `group_edit`,
		},
		{
			Display: false,
			Name:    echo.T(`删除分组`),
			Action:  `group_delete`,
		},
		{
			Display: false,
			Name:    echo.T(`日志列表`),
			Action:  `log`,
		},
		{
			Display: false,
			Name:    echo.T(`日志详情`),
			Action:  `log_view/:id`,
		},
		{
			Display: false,
			Name:    echo.T(`删除日志`),
			Action:  `log_delete`,
		},
	},
}
