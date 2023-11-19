package backend

import (
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v5/application/library/codec"
	"github.com/admpub/nging/v5/application/library/sessionguard"
)

func DecryptPassword(c echo.Context, pass string) (string, error) {
	var err error
	pass, err = codec.DefaultSM2DecryptHex(pass)
	if err != nil {
		return pass, err
	}
	pass, err = sessionguard.Unpack(c, pass)
	return pass, err
}
