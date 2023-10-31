package engine

import "github.com/webx-top/echo"

const (
	EnvironLocal     = `local`
	EnvironContainer = `container`
)

var Environs = echo.NewKVData().Add(EnvironLocal, `本机`).Add(EnvironContainer, `容器`)
