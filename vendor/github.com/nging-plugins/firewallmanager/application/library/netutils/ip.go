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

func parseCIDR(ip string) (netip.Addr, error) {
	if strings.Contains(ip, `/`) {
		pre, err := netip.ParsePrefix(ip)
		if err != nil {
			return netip.Addr{}, err
		}
		return pre.Addr(), nil
	}
	return netip.ParseAddr(ip)
}

func ValidateIP(ctx echo.Context, ip string) (int, error) {
	parts := strings.Split(ip, `-`)
	if len(parts) > 2 {
		return 0, fmt.Errorf(`%v: %v`, ctx.T(`IP 设置不正确`), ctx.T(`不支持多个范围值`))
	}
	if len(parts) < 2 {
		ipd, err := parseCIDR(ip)
		if err != nil {
			return 0, fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 解析失败`, ip), err)
		}
		err = ValidateAddress(ipd)
		if err != nil {
			err = fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 无效`, ip), err)
			return 0, err
		}
		var ipVer int
		switch {
		case ipd.Is4():
			ipVer = 4
		case ipd.Is6():
			ipVer = 6
		}
		return ipVer, nil
	}
	start, err := parseCIDR(parts[0])
	if err != nil {
		return 0, fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 解析失败`, parts[0]), err)
	}

	end, err := parseCIDR(parts[1])
	if err != nil {
		return 0, fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 解析失败`, parts[1]), err)
	}

	if err = ValidateAddress(start); err != nil {
		err = fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 无效`, parts[0]), err)
		return 0, err
	}

	if err = ValidateAddress(end); err != nil {
		err = fmt.Errorf(`%v: %w`, ctx.T(`IP (%v) 无效`, parts[1]), err)
		return 0, err
	}

	if end.Less(start) {
		err = fmt.Errorf(ctx.T("IP 设置错误: 起始 IP (%v) 不能大于终止 IP (%v)"), start.String(), end.String())
		return 0, err
	}

	var ipVer int
	switch {
	case start.Is4():
		ipVer = 4
	case start.Is6():
		ipVer = 6
	}
	return ipVer, nil
}
