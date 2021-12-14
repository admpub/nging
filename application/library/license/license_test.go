package license

import (
	"runtime"
	"testing"

	"github.com/admpub/log"
	"github.com/admpub/nging/v4/application/library/config"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/publicsuffix"
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
	err := Download(nil)
	if err != nil {
		panic(err)
	}
}

func TestLicenseLatestVersion(t *testing.T) {
	defer log.Close()
	_, err := LatestVersion(nil, false)
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
