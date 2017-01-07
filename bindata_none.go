// +build !bindata

package main

import (
	"fmt"
	"os"
)

func init() {
	binData = false
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
