//go:build !bindata
// +build !bindata

/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

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

package bindata

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/echo/middleware/render/driver"
	"github.com/webx-top/image"

	"github.com/admpub/nging/v4/application/cmd/event"
	"github.com/admpub/nging/v4/application/initialize/backend"
	"github.com/admpub/nging/v4/application/library/modal"
	"github.com/admpub/nging/v4/application/library/ntemplate"
	uploadLibrary "github.com/admpub/nging/v4/application/library/upload"
)

// StaticOptions static中间件选项
var StaticOptions = &middleware.StaticOptions{
	Root:     "",
	Path:     "/public/assets/",
	Fallback: []string{},
}

var PathAliases = ntemplate.PathAliases{}

// Initialize 初始化
func Initialize() {
	event.Bindata = false
	if len(StaticOptions.Root) == 0 {
		StaticOptions.Root = backend.AssetsDir
	}
	event.StaticMW = middleware.Static(StaticOptions)
	if !com.FileExists(event.FaviconPath) {
		log.Error(`not found favicon file: ` + event.FaviconPath)
	}
	event.FaviconHandler = func(c echo.Context) error {
		return c.File(event.FaviconPath)
	}
	image.WatermarkOpen = func(file string) (image.FileReader, error) {
		f, err := image.DefaultHTTPSystemOpen(file)
		if err != nil {
			if os.IsNotExist(err) {
				if strings.HasPrefix(file, uploadLibrary.UploadURLPath) || strings.HasPrefix(file, `/public/assets/`) {
					return os.Open(filepath.Join(echo.Wd(), file))
				}
			}
		}
		return f, err
	}
	modal.PathFixer = func(c echo.Context, file string) string {
		rpath := strings.TrimPrefix(file, backend.TemplateDir+`/`)
		rpath, ok := PathAliases.ParsePrefixOk(rpath)
		if ok {
			file = rpath
		}
		return file
	}
	backend.RendererDo = func(renderer driver.Driver) {
		renderer.SetTmplPathFixer(func(c echo.Context, tmpl string) string {
			rpath, ok := PathAliases.ParsePrefixOk(tmpl)
			if ok {
				return rpath
			}
			tmpl = filepath.Join(renderer.TmplDir(), tmpl)
			return tmpl
		})
	}
}
