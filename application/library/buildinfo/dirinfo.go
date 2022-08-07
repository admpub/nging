package buildinfo

import (
	"os"
	"path/filepath"
)

func NgingDir() string {
	return filepath.Join(os.Getenv(`GOPATH`), `src/github.com/admpub/nging`)
}

func NgingPluginsDir() string {
	return filepath.Join(os.Getenv(`GOPATH`), `src/github.com/nging-plugins`)
}
