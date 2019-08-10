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

package upload

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/tplfunc"
)

type APIKey interface {
	APIKey() string
}

func Token(values ...interface{}) string {
	urlValues := tplfunc.URLValues(values)
	var apiKey string
	if cfg, ok := echo.Get(`DefaultConfig`).(APIKey); ok {
		apiKey = cfg.APIKey()
	}
	return com.SafeBase64Encode(com.Token(apiKey, com.Str2bytes(urlValues.Encode())))
}
