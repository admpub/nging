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

package license

import (
	"errors"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/admpub/license_gen/lib"
	"github.com/admpub/once"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"golang.org/x/net/publicsuffix"

	"github.com/admpub/nging/v3/application/library/config"
)

type Mode int

const (
	ModeMachineID Mode = iota
	ModeDomain
)

var (
	trackerURL      = `https://www.webx.top/product/script/nging/tracker.js`
	productURL      = `https://www.webx.top/product/detail/nging`
	licenseURL      = `https://www.webx.top/product/license/nging`
	versionURL      = `https://www.webx.top/product/version/nging`
	licenseMode     = ModeMachineID
	licenseFileName = `license.key`
	licenseFile     = filepath.Join(echo.Wd(), licenseFileName)
	licenseExists   bool
	licenseError    = lib.UnlicensedVersion
	licenseData     *lib.LicenseData
	licenseVersion  string
	licensePackage  string
	machineID       string
	domain          string
	emptyLicense    = lib.LicenseData{}
	downloadOnce    once.Once
	downloadError   error
	downloadTime    time.Time
	// ErrLicenseNotFound 授权证书不存在
	ErrLicenseNotFound = errors.New(`License does not exist`)
	// SkipLicenseCheck 跳过授权检测
	SkipLicenseCheck = true
)

type ServerURL struct {
	Tracker         string //用于统计分析的js地址
	Product         string //该产品的详情介绍页面网址
	License         string //许可证验证和许可证下载API网址
	Version         string //该产品最新版本信息API网址
	LicenseFileName string //许可证文件名称
}

func (s *ServerURL) Apply() {
	if len(s.Tracker) > 0 {
		trackerURL = s.Tracker
	}
	if len(s.Product) > 0 {
		productURL = s.Product
	}
	if len(s.License) > 0 {
		licenseURL = s.License
	}
	if len(s.Version) > 0 {
		versionURL = s.Version
	}
	if len(s.LicenseFileName) > 0 {
		licenseFileName = s.LicenseFileName
		licenseFile = filepath.Join(echo.Wd(), licenseFileName)
	}
}

func SetServerURL(s *ServerURL) {
	if s != nil {
		s.Apply()
	}
}

func SetProductName(name string, domains ...string) {
	domain := `www.webx.top`
	if len(domains) > 0 && len(domains[0]) > 0 {
		domain = domains[0]
	}
	trackerURL = `https://` + domain + `/product/script/` + name + `/tracker.js`
	productURL = `https://` + domain + `/product/detail/` + name
	licenseURL = `https://` + domain + `/product/license/` + name
	versionURL = `https://` + domain + `/product/version/` + name
}

func SetProductDomain(domain string) {
	trackerURL = `https://` + domain + `/script/tracker.js`
	productURL = `https://` + domain + `/`
	licenseURL = `https://` + domain + `/license`
	versionURL = `https://` + domain + `/version`
}

func SetVersion(version string) {
	licenseVersion = version
}

func SetPackage(pkg string) {
	licensePackage = pkg
}

func Version() string {
	return licenseVersion
}

func Package() string {
	return licensePackage
}

func ProductURL() string {
	return productURL
}

func Domain() string {
	return domain
}

func SetDomain(_domain string) {
	if licenseMode != ModeDomain {
		licenseMode = ModeDomain
	}
	domain = _domain
}

func FullDomain() string {
	rootDomain := Domain()
	if len(rootDomain) == 0 {
		return rootDomain
	}
	rootDomain = strings.Trim(rootDomain, `.`)
	realDomain, _ := publicsuffix.EffectiveTLDPlusOne(rootDomain)
	if rootDomain == realDomain {
		return `www.` + realDomain
	}
	return rootDomain
}

func EqDomain(fullDomain string, rootDomain string) bool {
	rootDomain = strings.Trim(rootDomain, `.`)
	rootParts := strings.Split(rootDomain, `.`)
	fullParts := strings.Split(fullDomain, `.`)
	l := len(fullParts) - len(rootParts)
	if l < 0 {
		return false
	}
	//com.Dump(echo.H{`root`: rootParts, `full`: fullParts})
	for i, j := 0, len(rootParts); i < j; i++ {
		if rootParts[i] != fullParts[i+l] {
			return false
		}
	}
	return true
}

func ProductDetailURL() (url string) {
	url = ProductURL() + `?version=` + config.Version.Number
	switch licenseMode {
	case ModeMachineID:
		mid, err := MachineID()
		if err != nil {
			panic(err)
		}
		url += `&machineID=` + mid
	case ModeDomain:
		if len(Domain()) > 0 {
			url += `&domain=` + Domain()
		}
	default:
		panic(fmt.Sprintf(`unsupported license mode: %d`, licenseMode))
	}
	return
}

func TrackerURL() string {
	if trackerURL == `#` {
		return ``
	}
	return trackerURL + `?version=` + Version() + `&package=` + Package() + `&os=` + config.Version.BuildOS + `&arch=` + config.Version.BuildArch
}

func TrackerHTML() template.HTML {
	_trackerURL := TrackerURL()
	if len(_trackerURL) == 0 {
		return template.HTML(``)
	}
	return template.HTML(`<script type="text/javascript" async src="` + _trackerURL + `"></script>`)
}

func FilePath() string {
	return licenseFile
}

func FileName() string {
	return licenseFileName
}

func Exists() bool {
	return licenseExists
}

func Error() error {
	return licenseError
}

func License() lib.LicenseData {
	if licenseData == nil {
		return emptyLicense
	}
	return *licenseData
}

// MachineID 生成当前机器的机器码
func MachineID() (string, error) {
	if len(machineID) > 0 {
		return machineID, nil
	}
	addrs, err := lib.MACAddresses(false)
	if err != nil {
		return ``, err
	}
	if len(addrs) < 1 {
		return ``, lib.ErrorMachineID
	}
	cpuInfo, err := cpu.Info()
	if err != nil {
		return ``, err
	}
	var cpuID string
	if len(cpuInfo) > 0 {
		cpuID = cpuInfo[0].PhysicalID
		if len(cpuID) == 0 {
			cpuID = com.Md5(com.Dump(cpuInfo, false))
		}
	}
	machineID = com.MakePassword(lib.Hash(addrs[0])+`#`+cpuID, `coscms`, 3, 8, 19)
	return machineID, err
}

// FullLicenseURL 包含完整参数的授权网址
func FullLicenseURL(ctx echo.Context) string {
	return licenseURL + `?` + URLValues(ctx).Encode()
}

// URLValues 组装网址参数
func URLValues(ctx echo.Context) url.Values {
	v := url.Values{}
	v.Set(`os`, config.Version.BuildOS)
	v.Set(`arch`, config.Version.BuildArch)
	v.Set(`version`, licenseVersion)
	v.Set(`package`, licensePackage)
	if ctx != nil {
		v.Set(`source`, ctx.RequestURI())
	}
	switch licenseMode {
	case ModeMachineID:
		if len(machineID) == 0 {
			var err error
			machineID, err = MachineID()
			if err != nil {
				panic(fmt.Errorf(`failed to get machineID: %v`, err))
			}
		}
		v.Set(`machineID`, machineID)
	case ModeDomain:
		if len(Domain()) == 0 {
			panic(`license domain is required`)
		}
		v.Set(`domain`, Domain())
	default:
		panic(fmt.Sprintf(`unsupported license mode: %d`, licenseMode))
	}
	v.Set(`time`, time.Now().Format(`20060102-150405`))
	return v
}
