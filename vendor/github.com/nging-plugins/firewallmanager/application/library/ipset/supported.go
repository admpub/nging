package ipset

import "github.com/admpub/nging/v5/application/library/checkinstall"

var supported = checkinstall.New(`ipset`)

func IsSupported() bool {
	return supported.IsInstalled()
}

func ResetCheck() {
	supported.ResetCheck()
}
