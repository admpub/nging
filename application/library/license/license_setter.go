package license

import (
	"path/filepath"

	"github.com/admpub/license_gen/lib"
	"github.com/webx-top/echo"
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

func SetVersion(ver string) {
	version = ver
}

func SetPackage(pkg string) {
	packageName = pkg
}

func SetDomain(_domain string) {
	if licenseMode != ModeDomain {
		licenseMode = ModeDomain
	}
	domain = _domain
}

func SetError(err error) {
	lock4err.Lock()
	licenseError = err
	lock4err.Unlock()
}

func SetLicense(data *lib.LicenseData) {
	FireSetLicense(data)
	lock4data.Lock()
	licenseData = data
	lock4data.Unlock()
	switch licenseMode {
	case ModeDomain:
		if len(Domain()) == 0 {
			SetDomain(data.Info.Domain)
		}
	case ModeMachineID:
	}
}
