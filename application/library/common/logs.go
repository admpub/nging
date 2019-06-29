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

package common

import (
	"github.com/admpub/tail"
	"github.com/webx-top/echo"
)

// LogShow 获取日志内容用于显示
func LogShow(ctx echo.Context, logFile string, extensions ...echo.H) error {
	data := ctx.Data()
	var result echo.H
	if len(extensions) > 0 {
		result = extensions[0]
	}
	if result == nil {
		result = echo.H{}
	}
	if len(logFile) == 0 {
		result.Set(`content`, ctx.T(`没有日志文件`))
		data.SetData(result)
		return ctx.JSON(data)
	}
	lastLines := ctx.Formx(`lastLines`).Int()
	config := tail.Config{
		MustExist: true,
		LastLines: 50,
	}
	if lastLines > 0 {
		config.LastLines = lastLines
	}
	obj, err := tail.TailFile(logFile, config)
	if err != nil {
		data.SetError(err)
	} else {
		var content string
		for line := range obj.Lines {
			content += line.Text + "\n"
		}
		result.Set(`content`, content)
		data.SetData(result)
	}
	return ctx.JSON(data)
}
