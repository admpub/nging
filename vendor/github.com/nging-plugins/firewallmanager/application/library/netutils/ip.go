package netutils

import (
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/webx-top/echo"
)

var ErrInvalidPort = errors.New(`端口设置不正确: 超过区间 1~65535`)

// Validates an IP address
func ValidateAddress(ip netip.Addr) error {
	if !ip.IsValid() {
		return fmt.Errorf("address is zero")
	}

	if ip.IsUnspecified() {
		return fmt.Errorf("address is unspecified %v", ip.String())
	}

	return nil
}

func ValidateIP(ctx echo.Context, ip string) error {
	parts := strings.Split(ip, `-`)
	if len(parts) > 2 {
		return fmt.Errorf(`%v: %v`, ctx.T(`IP 设置不正确`), ctx.T(`不支持多个范围值`))
	}
	if len(parts) < 2 {
		ipd, err := netip.ParseAddr(ip)
		if err != nil {
			return fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 解析失败`, ip), err)
		}
		err = ValidateAddress(ipd)
		if err != nil {
			err = fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 无效`, ip), err)
		}
		return err
	}
	start, err := netip.ParseAddr(parts[0])
	if err != nil {
		return fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 解析失败`, parts[0]), err)
	}

	end, err := netip.ParseAddr(parts[1])
	if err != nil {
		return fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 解析失败`, parts[1]), err)
	}

	if err = ValidateAddress(start); err != nil {
		err = fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 无效`, parts[0]), err)
		return err
	}

	if err = ValidateAddress(end); err != nil {
		err = fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 无效`, parts[1]), err)
		return err
	}

	if end.Less(start) {
		err = fmt.Errorf(ctx.T("IP 设置错误: start address (%v) is after end address (%v)"), start.String(), end.String())
	}
	return err
}
