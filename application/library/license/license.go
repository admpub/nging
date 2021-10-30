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
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/admpub/license_gen/lib"
	"github.com/admpub/log"
	"github.com/admpub/once"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	"github.com/admpub/nging/v3/application/library/config"
	"github.com/admpub/nging/v3/application/library/restclient"
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
	mutex           sync.RWMutex
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

// Check 检查权限
func Check(ctx echo.Context) error {
	if SkipLicenseCheck {
		return nil
	}
	licenseError = validateFromOfficial(ctx)
	if licenseError != ErrConnectionFailed {
		return licenseError
	}
	//当官方服务器不可用时才验证本地许可证
	licenseError = Validate()
	return licenseError
}

func Ok(ctx echo.Context) bool {
	if SkipLicenseCheck {
		return true
	}
	switch licenseError {
	case nil:
		if licenseData == nil {
			licenseError = lib.UnlicensedVersion
			return false
		}
		if !licenseData.Info.Expiration.IsZero() && time.Now().After(licenseData.Info.Expiration) {
			licenseError = lib.ExpiredLicense
			return false
		}
		return true
	default:
		err := Check(ctx)
		if err == nil {
			licenseError = nil
			return true
		}
		log.Warn(err)
	}
	return false
}

// Validation 定义验证器
type Validation struct {
	NowVersions []string
}

// Validate 参数验证器
func (v *Validation) Validate(data *lib.LicenseData) error {
	if err := data.CheckExpiration(); err != nil {
		return err
	}
	if err := data.CheckVersion(v.NowVersions...); err != nil {
		return err
	}
	switch licenseMode {
	case ModeMachineID:
		mid, err := MachineID()
		if err != nil {
			return err
		}
		if data.Info.MachineID != mid {
			return lib.InvalidMachineID
		}
	case ModeDomain:
		return data.CheckDomain(Domain())
	default:
		panic(fmt.Sprintf(`unsupported license mode: %d`, licenseMode))
	}
	return nil
}

// Validate 验证授权
func Validate() error {
	licenseExists = com.FileExists(FilePath())
	if !licenseExists {
		licenseError = ErrLicenseNotFound
		return licenseError
	}
	b, err := ioutil.ReadFile(FilePath())
	if err != nil {
		return err
	}
	validator := &Validation{
		NowVersions: []string{licenseVersion},
	}
	licenseData, err = lib.CheckLicenseStringAndReturning(string(b), PublicKey(), validator)
	return err
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

func DownloadOnce(ctx echo.Context) error {
	downloadOnce.Do(func() {
		downloadTime = time.Now()
		downloadError = Download(ctx)
	})
	return downloadError
}

// Download 从官方服务器重新下载许可证
func Download(ctx echo.Context) error {
	operation := `获取授权证书失败：%v`
	client := restclient.Resty()
	client.SetHeader("Accept", "application/json")
	officialResponse := &OfficialResponse{}
	client.SetResult(officialResponse)
	fullURL := FullLicenseURL(ctx) + `&pipe=download`
	response, err := client.Get(fullURL)
	if err != nil {
		return fmt.Errorf(operation, err)
	}
	if response == nil {
		return ErrConnectionFailed
	}
	if response.IsError() {
		return fmt.Errorf(operation, string(response.Body()))
	}
	if officialResponse.Code != 1 {
		return fmt.Errorf(`%v`, officialResponse.Info)
	}
	if officialResponse.Data == nil {
		return ErrLicenseDownloadFailed
	}
	if com.FileExists(licenseFile) {
		err = os.Rename(licenseFile, licenseFile+`.`+time.Now().Format(`20060102150405`))
		if err != nil {
			return err
		}
	}
	licenseData = &officialResponse.Data.LicenseData
	b, err := com.JSONEncode(licenseData, `  `)
	if err != nil {
		b = []byte(err.Error())
	}
	err = ioutil.WriteFile(licenseFile, b, os.ModePerm)
	if err != nil {
		return fmt.Errorf(`保存授权证书失败：%v`, err)
	}
	return err
}
