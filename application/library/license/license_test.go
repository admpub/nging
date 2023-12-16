package license

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/admpub/log"
	"github.com/admpub/nging/v5/application/library/config"
	"github.com/admpub/pp/ppnocolor"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/defaults"
	"golang.org/x/net/publicsuffix"
)

func init() {
	config.Version.BuildOS = runtime.GOOS
	config.Version.BuildArch = runtime.GOARCH
	config.Version.Package = `free`
	config.Version.Number = `5.0.0`
	//*
	(&ServerURL{
		Tracker: `http://nging.coscms.com/product/script/nging/tracker.js`,
		Product: `http://nging.coscms.com/product/detail/nging`,
		License: `http://nging.coscms.com/product/license/nging`,
		Version: `http://nging.coscms.com/product/version/nging`,
	}).Apply()
	//*/
}

func TestLicenseDownload(t *testing.T) {
	err := Download(nil)
	if err != nil {
		panic(err)
	}
}

func TestLicenseLatestVersion(t *testing.T) {
	defer log.Close()
	ctx := defaults.NewMockContext()
	info, err := LatestVersion(ctx, ``, true)
	if err != nil {
		panic(err)
	}
	ppnocolor.Println(info)
	err = info.Extract()
	if err != nil {
		panic(err)
	}
	ppnocolor.Println(info.extractedDir)
	ppnocolor.Println(info.executable)
	ngingDir, _ := filepath.Abs(`./testdata`)
	if err != nil {
		panic(err)
	}
	echo.SetWorkDir(ngingDir)
	err = info.Upgrade(ctx, ngingDir)
	if err != nil {
		panic(err)
	}
}

func TestLicenseValidateFromOfficial(t *testing.T) {
	err := validateFromOfficial(nil)
	if err != nil {
		panic(err)
	}
}

func TestLicenseEqDomain(t *testing.T) {
	defer log.Close()
	assert.True(t, EqDomain(`www.webx.top`, `webx.top`))

	domain, err := publicsuffix.EffectiveTLDPlusOne(`www.webx.top`)
	assert.Nil(t, err)
	assert.Equal(t, `webx.top`, domain)

	domain, err = publicsuffix.EffectiveTLDPlusOne(`www.abc.com.cn`)
	assert.Nil(t, err)
	assert.Equal(t, `abc.com.cn`, domain)

	domain, err = publicsuffix.EffectiveTLDPlusOne(`com.cn`)
	assert.NotNil(t, err)
	assert.Equal(t, ``, domain)

	publicSuffix, icann := publicsuffix.PublicSuffix(`www.webx.top`)
	assert.True(t, icann)
	assert.Equal(t, `top`, publicSuffix)

	publicSuffix, icann = publicsuffix.PublicSuffix(`www.webx.x`)
	assert.False(t, icann)
	assert.Equal(t, `x`, publicSuffix)
}
