package license

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"time"

	"github.com/admpub/license_gen/lib"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// Check 检查权限
func Check(ctx echo.Context, content ...[]byte) error {
	if SkipLicenseCheck {
		return nil
	}
	var validateRemote bool
	if licenseMode == ModeDomain && len(Domain()) > 0 {
		licenseError = validateFromOfficial(ctx)
		if licenseError != ErrConnectionFailed {
			return licenseError
		}
	} else {
		validateRemote = true
	}
	//当官方服务器不可用时才验证本地许可证
	licenseError = Validate(content...)
	if licenseError == nil && validateRemote {
		licenseError = validateFromOfficial(ctx)
		if licenseError != ErrConnectionFailed {
			return licenseError
		}
	}
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
		if len(Domain()) == 0 {
			SetDomain(data.Info.Domain)
			return nil
		}
		return data.CheckDomain(Domain())
	default:
		panic(fmt.Sprintf(`unsupported license mode: %d`, licenseMode))
	}
	return nil
}

// Validate 验证授权
func Validate(content ...[]byte) (err error) {
	var b []byte
	if len(content) > 0 && len(content[0]) > 0 {
		b = content[0]
	} else {
		licenseExists = com.FileExists(FilePath())
		if !licenseExists {
			licenseError = ErrLicenseNotFound
			return licenseError
		}
		b, err = ioutil.ReadFile(FilePath())
		if err != nil {
			return
		}
	}
	validator := &Validation{
		NowVersions: []string{licenseVersion},
	}
	licenseData, err = lib.CheckLicenseStringAndReturning(string(b), PublicKey(), validator)
	return
}

func CheckSiteURL(siteURL string) error {
	u, err := url.Parse(siteURL)
	if err != nil {
		err = fmt.Errorf(`%s: %w`, siteURL, err)
		return err
	}
	if SkipLicenseCheck || LicenseMode() != ModeDomain {
		return nil
	}
	rootDomain := Domain()
	if len(rootDomain) == 0 {
		err = errors.New(`please set up the license first`)
		return err
	}
	fullDomain := u.Hostname()
	if !EqDomain(fullDomain, rootDomain) {
		err = fmt.Errorf(`domain "%s" and licensed domain "%s" is mismatched`, fullDomain, rootDomain)
		return err
	}
	return err
}
