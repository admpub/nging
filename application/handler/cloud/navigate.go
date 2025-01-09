package cloud

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var LeftNavigate = &navigate.Item{
	Display: true,
	Name:    echo.T(`云服务`),
	Action:  `cloud`,
	Icon:    `cloud`,
	Children: &navigate.List{
		{
			Display: true,
			Name:    echo.T(`云存储账号`),
			Action:  `storage`,
		},
		{
			Display: true,
			Name:    echo.T(`添加账号`),
			Action:  `storage_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改账号`),
			Action:  `storage_edit`,
		},
		{
			Display: false,
			Name:    echo.T(`删除账号`),
			Action:  `storage_delete`,
		},
		{
			Display: false,
			Name:    echo.T(`云存储文件管理`),
			Action:  `storage_file`,
		},
		{
			Display: true,
			Name:    echo.T(`文件备份`),
			Action:  `backup`,
		},
		{
			Display: false,
			Name:    echo.T(`添加备份配置`),
			Action:  `backup_add`,
			Icon:    `plus`,
		},
		{
			Display: false,
			Name:    echo.T(`修改备份配置`),
			Action:  `backup_edit`,
		},
		{
			Display: false,
			Name:    echo.T(`删除备份配置`),
			Action:  `backup_delete`,
		},
		{
			Display: false,
			Name:    echo.T(`启动备份任务`),
			Action:  `backup_start`,
		},
		{
			Display: false,
			Name:    echo.T(`停止备份任务`),
			Action:  `backup_stop`,
		},
		{
			Display: false,
			Name:    echo.T(`恢复备份文件`),
			Action:  `backup_restore`,
		},
		{
			Display: false,
			Name:    echo.T(`云备份日志列表`),
			Action:  `backup_log`,
		},
		{
			Display: false,
			Name:    echo.T(`云备份日志删除`),
			Action:  `backup_log_delete`,
		},
	},
}
