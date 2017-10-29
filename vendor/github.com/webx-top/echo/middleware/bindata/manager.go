/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/

package bindata

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	assetfs "github.com/admpub/go-bindata-assetfs"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/render/driver"
)

func NewTmplManager(fs http.FileSystem, templateDir ...string) driver.Manager {
	var prefix string
	if len(templateDir) > 0 {
		prefix = templateDir[0]
	} else {
		switch f := fs.(type) {
		case *assetfs.AssetFS:
			prefix = f.Prefix
		default:
			prefix = echo.Wd()
		}
	}
	prefix, _ = filepath.Abs(prefix)
	return &TmplManager{
		BaseManager: &driver.BaseManager{},
		FileSystem:  fs,
		Prefix:      prefix,
	}
}

type TmplManager struct {
	*driver.BaseManager
	http.FileSystem
	Prefix string
}

func (a *TmplManager) GetTemplate(fileName string) ([]byte, error) {
	fileName = strings.TrimPrefix(fileName, a.Prefix)
	fileName = filepath.ToSlash(fileName)
	file, err := a.FileSystem.Open(fileName)
	if err != nil {
		err = errors.New(fileName + `: ` + err.Error())
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	return b, err
}
