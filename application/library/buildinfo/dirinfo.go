package buildinfo

import (
	"path/filepath"

	"github.com/webx-top/com"
)

var GOPATH = com.GetGOPATHs()[0]

func NgingDir() string {
	return filepath.Join(GOPATH, `src/github.com/admpub/nging`)
}

func NgingPluginsDir() string {
	return filepath.Join(GOPATH, `src/github.com/nging-plugins`)
}
