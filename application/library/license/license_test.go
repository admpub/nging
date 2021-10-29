package license

import "testing"

func init() {
	(&ServerURL{
		Tracker: `http://nging.coscms.com/product/script/nging/tracker.js`,
		Product: `http://nging.coscms.com/product/detail/nging`,
		License: `http://nging.coscms.com/product/license/nging`,
		Version: `http://nging.coscms.com/product/version/nging`,
	}).Apply()
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
	machineID, err := MachineID()
	if err != nil {
		panic(err)
	}
	err = latestVersion(machineID, nil)
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
