package license

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/admpub/license_gen/lib"
	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

// Check 检查权限
func Check(ctx echo.Context) error {
	if SkipLicenseCheck {
		return nil
	}
	var validateRemote bool
	if licenseMode != ModeDomain || len(Domain()) > 0 {
		err := validateFromOfficial(ctx)
		if err != ErrConnectionFailed {
			return err
		}
	} else {
		validateRemote = true
	}
	//当官方服务器不可用时才验证本地许可证
	err := Validate()
	if err == nil && validateRemote {
		err = validateFromOfficial(ctx)
		if err == ErrConnectionFailed {
			err = nil
		}
	}
	return err
}

// VerifyPostLicenseContent 验证提交的证书内容
func VerifyPostLicenseContent(ctx echo.Context, content []byte) error {
	err := Validate(content)
	if err == nil {
		err = validateFromOfficial(ctx)
		if err == ErrConnectionFailed {
			err = nil
		}
	}
	SetError(err)
	return err
}

func Ok(ctx echo.Context) bool {
	if SkipLicenseCheck {
		return true
	}
	switch Error() {
	case nil:
		data := License()
		if data == emptyLicense {
			SetError(lib.UnlicensedVersion)
			return false
		}
		if !data.Info.Expiration.IsZero() && time.Now().After(data.Info.Expiration) {
			fi, err := os.Stat(FilePath())
			if err == nil {
				if fi.ModTime().After(licenseModTime) {
					licenseModTime = fi.ModTime()
					goto CHECK
				}
			}
			SetError(lib.ExpiredLicense)
			return false
		}
		return true

	CHECK:
		fallthrough

	default:
		err := Check(ctx)
		SetError(err)
		if err == nil {
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

func ReadLicenseKeyFile() ([]byte, error) {
	return ioutil.ReadFile(FilePath())
}

// Validate 验证授权
func Validate(content ...[]byte) (err error) {
	var b []byte
	if len(content) > 0 && len(content[0]) > 0 {
		b = content[0]
	} else {
		licenseExists = com.FileExists(FilePath())
		if !licenseExists {
			return ErrLicenseNotFound
		}
		b, err = ReadLicenseKeyFile()
		if err != nil {
			return
		}
	}
	validator := &Validation{
		NowVersions: []string{strings.SplitN(licenseVersion, `-`, 2)[0]},
	}
	var pubKey string
	b, pubKey = LicenseDecode(b)
	if len(pubKey) > 0 {
		if publicKey != pubKey {
			SetPublicKey(pubKey)
		}
	} else {
		pubKey = publicKey
	}
	var data *lib.LicenseData
	data, err = lib.CheckLicenseStringAndReturning(com.Bytes2str(b), pubKey, validator)
	if err == nil {
		SetLicense(data)
	}
	return
}

func CheckSiteURL(siteURL string) error {
	if SkipLicenseCheck || LicenseMode() != ModeDomain {
		return nil
	}
	if len(siteURL) == 0 {
		return nil
	}
	u, err := url.Parse(siteURL)
	if err != nil {
		err = fmt.Errorf(`%s: %w`, siteURL, err)
		return err
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
