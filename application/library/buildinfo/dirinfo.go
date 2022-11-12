package buildinfo

import (
	"path/filepath"

	"github.com/webx-top/com"
)

// 以下代码仅用于开发模式

var GOPATH = com.GetGOPATHs()[0]

func NgingDir() string {
	return filepath.Join(GOPATH, `src/github.com/admpub/nging`)
}

func NgingPluginsDir() string {
	return filepath.Join(GOPATH, `src/github.com/nging-plugins`)
}
