// Copyright 2016 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package assets

//go:generate go get github.com/admpub/statik
//go:generate statik -src=./frps/static -dest=./frps -k=server
//go:generate statik -src=./frpc/static -dest=./frpc -k=client
//go:generate go fmt ./frps/statik/statik.go
//go:generate go fmt ./frpc/statik/statik.go

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/admpub/statik/fs"
)

var (
	// FileSystem store static files in memory by statik
	FileSystem http.FileSystem

	// if prefix is not empty, we get file content from disk
	prefixPath string

	defaultAssetKey = `server`
	fileSystems     = map[string]http.FileSystem{}
)

func getAssetKey(assetKeys ...string) string {
	var assetKey string
	if len(assetKeys) > 0 {
		assetKey = assetKeys[0]
	}
	if len(assetKey) == 0 {
		assetKey = defaultAssetKey
	}
	return assetKey
}

func FS(assetKeys ...string) http.FileSystem {
	assetKey := getAssetKey(assetKeys...)
	fs, _ := fileSystems[assetKey]
	return fs
}

// Load if path is empty, load assets in memory
// or set FileSystem using disk files
func Load(path string, assetKeys ...string) (err error) {
	assetKey := getAssetKey(assetKeys...)
	prefixPath = path
	if prefixPath != "" {
		fileSystems[assetKey] = http.Dir(prefixPath)
	} else {
		fileSystems[assetKey], err = fs.NewWithKey(assetKey)
	}
	if assetKey == defaultAssetKey {
		FileSystem = fileSystems[assetKey]
	}
	return err
}

func ReadFile(file string, assetKeys ...string) (content string, err error) {
	assetKey := getAssetKey(assetKeys...)
	_, ok := fileSystems[assetKey]
	if prefixPath == "" && ok {
		file, err := fileSystems[assetKey].Open(path.Join("/", file))
		if err != nil {
			return content, err
		}
		buf, err := ioutil.ReadAll(file)
		if err != nil {
			return content, err
		}
		content = string(buf)
	} else {
		file, err := os.Open(path.Join(prefixPath, file))
		if err != nil {
			return content, err
		}
		buf, err := ioutil.ReadAll(file)
		if err != nil {
			return content, err
		}
		content = string(buf)
	}
	return content, err
}
