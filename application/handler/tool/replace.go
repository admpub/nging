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

package tool

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/webx-top/com"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

type replaceDef struct {
	Find       string //查找字串
	ReplaceAs  string //替换为
	Source     string //来源文件或文件夹
	SaveAs     string //另存到文件夹(为空时代表保存到原文件)
	Extensions string //扩展名(Source为文件夹时有效)
	IsRegexp   bool   //是否使用正则表达式
}

func Replace(c echo.Context) (err error) {
	if c.IsPost() {
		repDef := &replaceDef{}
		err = c.MustBind(repDef)
		if err != nil {
			return
		}
		fi, err := os.Stat(repDef.Source)
		if err != nil {
			return err
		}
		data := c.Data()
		exts := []string{}
		for _, ext := range strings.Split(repDef.Extensions, `,`) {
			ext = strings.TrimSpace(ext)
			if len(ext) == 0 {
				continue
			}
			exts = append(exts, strings.ToLower(ext))
		}
		var re *regexp.Regexp
		if repDef.IsRegexp {
			re, err = regexp.Compile(repDef.Find)
			if err != nil {
				return err
			}
		}
		replaFn := func(file string) error {
			f, e := os.Open(file)
			if e != nil {
				return e
			}
			defer func() {
				f.Close()
			}()
			b, e := ioutil.ReadAll(f)
			if e != nil {
				return e
			}
			content := engine.Bytes2str(b)
			if re != nil {
				if !re.MatchString(content) {
					return nil
				}
				content = re.ReplaceAllString(content, repDef.ReplaceAs)
			} else {
				if strings.Index(content, repDef.Find) == -1 {
					return nil
				}
				content = strings.Replace(content, repDef.Find, repDef.ReplaceAs, -1)
			}
			if len(repDef.SaveAs) > 0 {
				fi, err := os.Stat(repDef.SaveAs)
				var isDir bool
				if err != nil {
					if os.IsNotExist(err) {
						err = os.MkdirAll(repDef.SaveAs, os.ModePerm)
						if err != nil {
							return err
						}
						isDir = true
					} else {
						return err
					}
				} else {
					isDir = fi.IsDir()
				}
				saveAs := repDef.SaveAs
				if isDir {
					saveAs = filepath.Join(repDef.SaveAs, filepath.Base(file))
				}
				e = com.Copy(file, saveAs)
				if e != nil {
					return e
				}
				f.Close()
				f, e = os.Open(saveAs)
				if e != nil {
					return e
				}
			}
			_, e = f.WriteString(content)
			return e
		}
		if fi.IsDir() {
			err = filepath.Walk(repDef.Source, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.IsDir() {
					return nil
				}
				if len(exts) > 0 {
					ext := strings.ToLower(filepath.Ext(info.Name()))
					var allow bool
					for _, ex := range exts {
						if ex == ext {
							allow = true
							break
						}
					}
					if !allow {
						return nil
					}
				}
				return replaFn(path)
			})
		} else {
			err = replaFn(repDef.Source)
		}
		if err != nil {
			data.SetError(err)
		}
		return c.JSON(data)
	}
	return c.Render(`/tool/replace`, nil)
}
