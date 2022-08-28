package model

import "github.com/webx-top/echo"

var SafeItems = echo.NewKVData().
	Add(`gauth_bind`, `两步验证`).
	Add(`password`, `修改密码`)
