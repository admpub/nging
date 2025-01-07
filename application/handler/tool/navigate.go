package tool

import (
	"github.com/coscms/webcore/library/navigate"
	"github.com/webx-top/echo"
)

var TopNavigate = navigate.List{
	{
		Display: true,
		Name:    echo.T(`IP归属地`),
		Action:  `ip`,
	},
	{
		Display: true,
		Name:    echo.T(`Base64解码`),
		Action:  `base64`,
	},
	{
		Display: true,
		Name:    echo.T(`URL解码`),
		Action:  `url`,
	},
	{
		Display: true,
		Name:    echo.T(`时间戳转换`),
		Action:  `timestamp`,
	},
	{
		Display: true,
		Name:    echo.T(`正则表达式测试`),
		Action:  `regexp_test`,
	},
	{
		Display: true,
		Name:    echo.T(`附件网址替换`),
		Action:  `replaceurl`,
	},
	{
		Display: true,
		Name:    echo.T(`生成密码`),
		Action:  `gen_password`,
		Target:  `ajax`,
	},
}
