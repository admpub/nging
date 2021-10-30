package license

import (
	"runtime"
	"testing"

	"github.com/admpub/log"
	"github.com/admpub/nging/v3/application/library/config"
)

func init() {
	(&ServerURL{
		Tracker: `http://nging.coscms.com/product/script/nging/tracker.js`,
		Product: `http://nging.coscms.com/product/detail/nging`,
		License: `http://nging.coscms.com/product/license/nging`,
		Version: `http://nging.coscms.com/product/version/nging`,
	}).Apply()
	config.Version.BuildOS = runtime.GOOS
	config.Version.BuildArch = runtime.GOARCH
}

func TestLicenseDownload(t *testing.T) {
	machineID, err := MachineID()
	if err != nil {
		panic(err)
	}
	err = Download(machineID, nil)
	if err != nil {
		panic(err)
	}
}

func TestLicenseLatestVersion(t *testing.T) {
	defer log.Close()
	machineID, err := MachineID()
	if err != nil {
		panic(err)
	}
	_, err = LatestVersion(machineID, nil, true)
	if err != nil {
		panic(err)
	}
}

func TestLicenseValidateFromOfficial(t *testing.T) {
	machineID, err := MachineID()
	if err != nil {
		panic(err)
	}
	err = validateFromOfficial(machineID, nil)
	if err != nil {
		panic(err)
	}
}
