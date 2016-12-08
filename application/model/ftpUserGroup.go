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
package model

import (
	"github.com/admpub/caddyui/application/dbschema"
	"github.com/webx-top/db"
	"github.com/webx-top/echo"
)

func NewFtpUserGroup(ctx echo.Context) *FtpUserGroup {
	return &FtpUserGroup{
		FtpUserGroup: &dbschema.FtpUserGroup{},
		Base:         &Base{Context: ctx},
	}
}

type FtpUserGroup struct {
	*dbschema.FtpUserGroup
	*Base
}

func (f *FtpUserGroup) Exists(name string) (bool, error) {
	n, e := f.Param().SetArgs(db.Cond{`name`: name}).Count()
	return n > 0, e
}

func (f *FtpUserGroup) ExistsOther(name string, id uint) (bool, error) {
	n, e := f.Param().SetArgs(db.Cond{`name`: name, `id <>`: id}).Count()
	return n > 0, e
}

func (f *FtpUserGroup) ListByActive(page int, size int) (func() int64, []*dbschema.FtpUserGroup, error) {
	count, err := f.List(nil, nil, page, size, db.Cond{`disabled`: `N`})
	if err == nil {
		return count, f.Objects(), err
	}
	return count, nil, err
}
