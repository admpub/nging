package license

import (
	"github.com/webx-top/echo/middleware/tplfunc"
)

func init() {
	tplfunc.TplFuncMap[`HasFeature`] = HasFeature
	tplfunc.TplFuncMap[`HasAnyFeature`] = HasAnyFeature
	tplfunc.TplFuncMap[`LicenseDomain`] = Domain
	tplfunc.TplFuncMap[`LicensePackage`] = Package
	tplfunc.TplFuncMap[`LicenseProductURL`] = ProductURL
	tplfunc.TplFuncMap[`LicenseSkipCheck`] = func() bool { return SkipLicenseCheck }
}
