// +build !bindata

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

package bindata

import (
	"os"
	"strings"
	"path/filepath"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware"
	"github.com/webx-top/image"

	"github.com/admpub/nging/application/cmd/event"
	"github.com/admpub/nging/application/initialize/backend"
	"github.com/admpub/nging/application/registry/upload/helper"
)

var StaticOptions = &middleware.StaticOptions{
	Root:     "",
	Path:     "/public/assets/",
	Fallback: []string{},
}

func Initialize() {
	event.Bindata = false
	if len(StaticOptions.Root) == 0 {
		StaticOptions.Root = backend.AssetsDir
	}
	event.StaticMW = middleware.Static(StaticOptions)
	faviconPath := filepath.Join(echo.Wd(), event.FaviconPath)
	event.FaviconHandler = func(c echo.Context) error {
		return c.File(faviconPath)
	}
	image.WatermarkOpen = func(file string) (image.FileReader, error) {
		f, err := image.DefaultHTTPSystemOpen(file)
		if err != nil {
			if os.IsNotExist(err) {
				if strings.HasPrefix(file, helper.DefaultUploadURLPath) || strings.HasPrefix(file, `/public/assets/`) {
					return os.Open(filepath.Join(echo.Wd(), file))
				}
			}
		}
		return f, err
	}
}
