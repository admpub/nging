package config

import "github.com/webx-top/echo"

var TTLs = echo.NewKVData().
	Add(``, `自动`).Add(`1`, `1秒`).Add(`5`, `5秒`).Add(`10`, `10秒`).
	Add(`60`, `1分钟`).Add(`120`, `2分钟`).Add(`300`, `5分钟`).
	Add(`600`, `10分钟`).Add(`1800`, `30分钟`).
	Add(`3600`, `1小时`)
