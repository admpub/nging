// +build bindata

package main

import (
	"net/http"
	"strings"

	"github.com/admpub/nging/application/library/modal"
	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/webx-top/echo/middleware/bindata"
)

func init() {
	binData = true

	staticMW = bindata.Static("/public/", &assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "",
	})

	tmplMgr = bindata.NewTmplManager(&assetfs.AssetFS{
		Asset:     Asset,
		AssetDir:  AssetDir,
		AssetInfo: AssetInfo,
		Prefix:    "template",
	})

	modal.ReadConfigFile = func(file string) ([]byte, error) {
		file = strings.TrimPrefix(file, `./template`)
		return tmplMgr.GetTemplate(file)
	}

	langFSFunc = func(dir string) http.FileSystem {
		return &assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
			Prefix:    dir,
		}
	}
}
