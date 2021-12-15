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

package mysql

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	md2html "github.com/russross/blackfriday"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func (m *mySQL) Analysis() error {
	sql := m.Form(`sql`)
	data := m.Data()
	if len(sql) == 0 {
		data.SetInfo(m.T(`请输入SQL语句`), 0)
		return m.JSON(data)
	}
	command := `soar`
	var extension string
	if com.IsWindows {
		extension = `.exe`
	}
	_, err := exec.LookPath(command + extension)
	if err != nil {
		files := []string{
			filepath.Join(echo.Wd(), `support`, command+`_`+runtime.GOOS+`_`+runtime.GOARCH) + extension,
			filepath.Join(echo.Wd(), `support`, command) + extension,
		}
		for _, support := range files {
			if com.FileExists(support) {
				err = nil
				command = support
				break
			}
		}

		if err != nil {
			gpath := os.Getenv("GOPATH")
			if len(gpath) > 0 {
				command = filepath.Join(gpath, `src/github.com/XiaoMi/soar`, command) + extension
				if com.FileExists(command) {
					err = nil
				}
			}
		}
	} else {
		command += extension
	}
	if err != nil {
		data.SetError(err)
		return m.JSON(data)
	}

	params := []string{
		command,
	}
	output := []byte{}
	cmd := com.CreateCmd(params, func(b []byte) error {
		output = append(output, b...)
		return nil
	})
	reader := strings.NewReader(sql)
	cmd.Stdin = reader
	err = cmd.Run()
	if err != nil {
		data.SetError(err)
	} else {
		output = md2html.MarkdownCommon(output)
		data.SetData(com.Bytes2str(output))
	}
	return m.JSON(data)
}
