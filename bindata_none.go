// +build !bindata

package main

import (
	"fmt"
	"os"

	"github.com/webx-top/echo/middleware"
)

func init() {
	binData = false
	staticMW = middleware.Static(&middleware.StaticOptions{
		Root: "./public",
		Path: "/public/",
	})
}

func Asset(name string) ([]byte, error) {
	return nil, fmt.Errorf("Asset %s not found", name)
}

func AssetDir(name string) ([]string, error) {
	return nil, fmt.Errorf("Asset %s not found", name)
}

func AssetInfo(name string) (os.FileInfo, error) {
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}
