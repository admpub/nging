package nftables

import "github.com/admpub/nging/v5/application/library/checkinstall"

var supported = checkinstall.New(`nft`)

func IsSupported() bool {
	return supported.IsInstalled()
}

func ResetCheck() {
	supported.ResetCheck()
}
