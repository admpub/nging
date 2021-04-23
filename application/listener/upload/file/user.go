/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present  Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

// Package file 监听文件更改，并记录到file表
package file

import (
	"fmt"

	"github.com/webx-top/db"
	"github.com/webx-top/db/lib/factory"

	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/library/fileupdater/listener"
)

func init() {
	// - 监听数据表事件

	// - user 表事件
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.NgingUser)
		tableID = fmt.Sprint(fm.Id) //! 表中数据的行ID
		content = fm.Avatar         //! 这里使用保存原始图片网址的字段
		property = listener.NewPropertyWith(
			fm,
			db.Cond{`id`: fm.Id}, //! 更新行内保存原始图片网址的字段(这里为avatar字段)的条件
			//! ... 指定更新行数据时，需要额外更新的字段及其值的生成方式
		)
		return
	}, false).SetTable(`nging_user`, `avatar`).ListenDefault()

	// - config 表事件
	listener.New(func(m factory.Model) (tableID string, content string, property *listener.Property) {
		fm := m.(*dbschema.NgingConfig)
		tableID = fm.Group + `.` + fm.Key
		content = fm.Value
		property = getConfigEventAttrs(fm).GenUpdater(fm, db.And(
			db.Cond{`key`: fm.Key},
			db.Cond{`group`: fm.Group},
		))
		return
	}, false).SetTable(`nging_config`, `value`).ListenDefault()
}

func getConfigEventAttrs(confM *dbschema.NgingConfig) *listener.Property {
	property := &listener.Property{}
	switch confM.Type {
	case `html`, `json`:
		property.SetEmbedded(true)
	case `image`, `video`, `audio`, `file`:
	case `list`:
		property.SetSeperator(`,`)
	default:
		property.SetExit(true)
	}
	return property
}
