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

package oauth

import (
	"github.com/admpub/nging/application/dbschema"
	"github.com/admpub/nging/application/registry/settings"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func init() {
	settings.RegisterDecoder(`oauth`, func(v *dbschema.Config, r echo.H) error {
		jsonData := NewConfig()
		if len(v.Value) > 0 {
			com.JSONDecode([]byte(v.Value), jsonData)
		}
		jsonData.On = v.Disabled != `Y`
		r[`ValueObject`] = jsonData
		return nil
	})
	settings.RegisterEncoder(`oauth`, func(v *dbschema.Config, r echo.H) ([]byte, error) {
		oauthConfig := NewConfig().FromStore(v.Key, r)
		return com.JSONEncode(oauthConfig)
	})
}
