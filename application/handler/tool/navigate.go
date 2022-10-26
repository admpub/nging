package tool

import "github.com/admpub/nging/v5/application/registry/navigate"

var TopNavigate = navigate.List{
	{
		Display: true,
		Name:    `IP归属地`,
		Action:  `ip`,
	},
	{
		Display: true,
		Name:    `Base64解码`,
		Action:  `base64`,
	},
	{
		Display: true,
		Name:    `URL解码`,
		Action:  `url`,
	},
	{
		Display: true,
		Name:    `时间戳转换`,
		Action:  `timestamp`,
	},
	{
		Display: true,
		Name:    `附件网址替换`,
		Action:  `replaceurl`,
	},
	{
		Display: true,
		Name:    `生成密码`,
		Action:  `gen_password`,
		Target:  `ajax`,
	},
}
