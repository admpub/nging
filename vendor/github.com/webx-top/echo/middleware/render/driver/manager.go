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
package driver

import (
	"github.com/webx-top/echo/logger"
)

type Manager interface {
	Start() error
	Close()
	ClearCallback()
	AddCallback(rootDir string, callback func(name, typ, event string))
	DelCallback(rootDir string)
	ClearAllows()
	AddAllow(allows ...string)
	DelAllow(allow string)
	ClearIgnores()
	AddIgnore(ignores ...string)
	DelIgnore(ignore string)
	AddWatchDir(ppath string) (err error)
	CancelWatchDir(oldDir string) (err error)
	ChangeWatchDir(oldDir string, newDir string) (err error)
	SetLogger(logger.Logger)
	ClearCache()
	GetTemplate(string) ([]byte, error)
}

var _ Manager = &BaseManager{}

type BaseManager struct {
}

func (b *BaseManager) Start() error                                                       { return nil }
func (b *BaseManager) Close()                                                             {}
func (b *BaseManager) ClearCallback()                                                     {}
func (b *BaseManager) AddCallback(rootDir string, callback func(name, typ, event string)) {}
func (b *BaseManager) DelCallback(rootDir string)                                         {}
func (b *BaseManager) ClearAllows()                                                       {}
func (b *BaseManager) AddAllow(allows ...string)                                          {}
func (b *BaseManager) DelAllow(allow string)                                              {}
func (b *BaseManager) ClearIgnores()                                                      {}
func (b *BaseManager) AddIgnore(ignores ...string)                                        {}
func (b *BaseManager) DelIgnore(ignore string)                                            {}
func (b *BaseManager) AddWatchDir(ppath string) (err error)                               { return nil }
func (b *BaseManager) CancelWatchDir(oldDir string) (err error)                           { return nil }
func (b *BaseManager) ChangeWatchDir(oldDir string, newDir string) (err error)            { return nil }
func (b *BaseManager) SetLogger(logger.Logger)                                            {}
func (b *BaseManager) ClearCache()                                                        {}
func (b *BaseManager) GetTemplate(string) ([]byte, error)                                 { return nil, nil }
