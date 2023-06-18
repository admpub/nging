package netutils

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/webx-top/echo"
)

func ValidatePort(ctx echo.Context, port string) error {
	port = strings.ReplaceAll(port, `:`, `-`)
	parts := strings.Split(port, `-`)
	var isRange bool
	if len(parts) == 1 {
		parts = strings.Split(port, `,`)
	} else {
		isRange = true
	}
	partsUint16 := make([]uint16, len(parts))
	for k, p := range parts {
		i, err := strconv.ParseUint(p, 10, 16)
		if err != nil {
			return fmt.Errorf(`%v: %v`, ctx.T(`端口设置不正确`), err.Error())
		}
		if i < 1 || i > 65535 {
			return ErrInvalidPort
		}
		partsUint16[k] = uint16(i)
	}
	if isRange {
		if len(parts) > 2 {
			return fmt.Errorf(`%v: %v`, ctx.T(`端口设置不正确`), ctx.T(`不支持多个范围值`))
		}
		if len(partsUint16) == 2 {
			if partsUint16[1] < partsUint16[0] {
				return fmt.Errorf(ctx.T("端口设置不正确: starting port (%v) is higher than ending port (%v)"), partsUint16[0], partsUint16[1])
			}
		}
	}
	return nil
}
