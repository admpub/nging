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

package driver

import (
	"time"
)

type Driver interface {
	AddDefault() error
	Enabled(on bool) error
	Reset() error // 重置（恢复到出厂设置）
	Clear() error // 清空本系统创建的所有规则
	Import(wfwFile string) error
	Export(wfwFile string) error
	Insert(rules ...Rule) error
	AsWhitelist(table, chain string) error
	Append(rules ...Rule) error
	Update(rule Rule) error
	Delete(rules ...Rule) error
	Exists(rule Rule) (bool, error)
	Ban(ips []string, expires time.Duration) error
	FindPositionByID(table, chain string, id uint) (uint, error)
}
