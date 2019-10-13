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
