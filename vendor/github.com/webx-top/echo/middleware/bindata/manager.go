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
	"io/ioutil"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/webx-top/echo/logger"
)

func NewTmplManager(fs *assetfs.AssetFS) *TmplManager {
	return &TmplManager{
		AssetFS: fs,
	}
}

type TmplManager struct {
	*assetfs.AssetFS
}

func (a *TmplManager) Close()                                            {}
func (a *TmplManager) SetOnChangeCallback(func(name, typ, event string)) {}
func (a *TmplManager) SetLogger(logger.Logger)                           {}
func (a *TmplManager) ClearCache()                                       {}
func (a *TmplManager) GetTemplate(fileName string) ([]byte, error) {
	file, err := a.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	return b, err
}
func (a *TmplManager) Init(logger logger.Logger, rootDir string, reload bool, allows ...string) {}
