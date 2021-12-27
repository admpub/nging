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

package common

import (
	"errors"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/admpub/decimal"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// Ok 操作成功
func Ok(v string) Successor {
	return NewOk(v)
}

// Err 获取错误信息
func Err(ctx echo.Context, err error) (ret interface{}) {
	if err != nil {
		return ProcessError(ctx, err)
	}
	flash := ctx.Flash()
	if flash != nil {
		if errMsg, ok := flash.(string); ok {
			ret = errors.New(errMsg)
		} else {
			ret = flash
		}
	}
	return
}

// SendOk 记录成功信息
func SendOk(ctx echo.Context, msg string) {
	if ctx.IsAjax() || ctx.Format() != echo.ContentTypeHTML {
		ctx.Data().SetInfo(msg, 1)
		return
	}
	ctx.Session().AddFlash(Ok(msg))
}

// SendFail 记录失败信息
func SendFail(ctx echo.Context, msg string) {
	if ctx.IsAjax() || ctx.Format() != echo.ContentTypeHTML {
		ctx.Data().SetInfo(msg, 0)
		return
	}
	ctx.Session().AddFlash(msg)
}

// SendErr 记录错误信息 (SendFail的别名)
func SendErr(ctx echo.Context, err error) {
	err = ProcessError(ctx, err)
	SendFail(ctx, err.Error())
}

type ConfigFromDB interface {
	ConfigFromDB() echo.H
}

// GetLocalIP 获取本机网卡IP
func GetLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIPNet bool
	)
	ipv4 = `127.0.0.1`
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIPNet = addr.(*net.IPNet); isIPNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}
	return
}

var notWordRegexp = regexp.MustCompile(`[^\w]+`)

// LookPath 获取二进制可执行文件路径
func LookPath(bin string, otherPaths ...string) (string, error) {
	envVarName := `NGING_` + notWordRegexp.ReplaceAllString(strings.ToUpper(bin), `_`) + `_PATH`
	envVarValue := os.Getenv(envVarName)
	if len(envVarValue) > 0 {
		if com.IsFile(envVarValue) {
			return envVarValue, nil
		}
		envVarValue = filepath.Join(envVarValue, bin)
		if com.IsFile(envVarValue) {
			return envVarValue, nil
		}
	}
	findPath, err := exec.LookPath(bin)
	if err == nil {
		return findPath, err
	}
	if !errors.Is(err, exec.ErrNotFound) {
		return findPath, err
	}
	for _, binPath := range otherPaths {
		binPath = filepath.Join(binPath, bin)
		if com.IsFile(binPath) {
			return binPath, nil
		}
	}
	return findPath, err
}

func SeekLinesWithoutComments(r io.Reader) (string, error) {
	var content string
	err := com.SeekLines(r, WithoutCommentsLineParser(func(line string) error {
		content += line + "\n"
		return nil
	}))
	return content, err
}

func WithoutCommentsLineParser(exec func(string) error) func(string) error {
	var commentStarted bool
	return func(line string) error {
		lineClean := strings.TrimSpace(line)
		if len(lineClean) == 0 {
			return nil
		}
		if commentStarted {
			if strings.HasSuffix(lineClean, `*/`) {
				commentStarted = false
			}
			return nil
		}
		switch lineClean[0] {
		case '#':
			return nil
		case '/':
			if len(lineClean) > 1 {
				switch lineClean[1] {
				case '/':
					return nil
				case '*':
					commentStarted = true
					return nil
				}
			}
		}

		//content += line + "\n"
		return exec(line)
	}
}

func Float64Sum(numbers ...float64) float64 {
	if len(numbers) == 0 {
		return 0
	}
	d := decimal.NewFromFloat(numbers[0])
	if len(numbers) > 1 {
		var d2 decimal.Decimal
		for _, number := range numbers[1:] {
			d2 = decimal.NewFromFloat(number)
			d = d.Add(d2)
		}
	}
	number, _ := d.Float64()
	return number
}

func Float32Sum(numbers ...float32) float32 {
	if len(numbers) == 0 {
		return 0
	}
	d := decimal.NewFromFloat32(numbers[0])
	if len(numbers) > 1 {
		var d2 decimal.Decimal
		for _, number := range numbers[1:] {
			d2 = decimal.NewFromFloat32(number)
			d = d.Add(d2)
		}
	}
	number, _ := d.Float64()
	return float32(number)
}
