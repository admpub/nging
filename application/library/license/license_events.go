package license

import "github.com/admpub/license_gen/lib"

func OnSetLicense(fn func(*lib.LicenseData)) {
	if fn == nil {
		return
	}
	onSetLicenseHooks = append(onSetLicenseHooks, fn)
}

func FireSetLicense(data *lib.LicenseData) {
	for _, fn := range onSetLicenseHooks {
		fn(data)
	}
}
