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

package index

import (
	"fmt"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/engine"
)

var M3U8Template = `#EXTM3U
#EXT-X-VERSION:3
#EXT-X-MEDIA-SEQUENCE:0
#EXT-X-ALLOW-CACHE:YES
#EXT-X-TARGETDURATION:86400
#EXTINF:86400,
%s
#EXT-X-ENDLIST
`

func TS2M3U8(ctx echo.Context) error {
	ctx.Response().Header().Set("Content-Type", "application/vnd.apple.mpegURL")
	tsFile := ctx.Form(`ts`)
	m3u8Content := fmt.Sprintf(M3U8Template, tsFile)
	return ctx.Blob(engine.Str2bytes(m3u8Content))
}
