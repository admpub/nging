package license

import (
	"os"
	"path/filepath"
	"time"

	"github.com/admpub/license_gen/lib"
	"github.com/webx-top/echo"
)

// Save 保存授权文件
func Save(b []byte) error {
	return os.WriteFile(licenseFile, b, os.ModePerm)
}

// Generate 生成演示版证书
func Generate(privBytes []byte, pemSaveDirs ...string) error {
	var err error
	if privBytes == nil {
		var pubBytes []byte
		pubBytes, privBytes, err = lib.GenerateCertificateData(2048)
		if err != nil {
			return err
		}
		publicKey = string(pubBytes)
		var pemSaveDir string
		if len(pemSaveDirs) > 0 {
			pemSaveDir = pemSaveDirs[0]
		} else {
			pemSaveDir = filepath.Join(echo.Wd(), `data`)
		}
		if len(pemSaveDir) > 0 {
			err = os.WriteFile(filepath.Join(pemSaveDir, `nging.pem.pub`), pubBytes, os.ModePerm)
			if err != nil {
				return err
			}
			err = os.WriteFile(filepath.Join(pemSaveDir, `nging.pem`), privBytes, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	info := &lib.LicenseInfo{
		Name:       `demo`,
		LicenseID:  `0`,
		Version:    Version(),
		Package:    Package(),
		Expiration: time.Now().Add(30 * 24 * time.Hour),
	}
	info.MachineID, err = MachineID()
	if err != nil {
		return err
	}
	licBytes, err := lib.GenerateLicense(info, string(privBytes))
	if err != nil {
		return err
	}
	return Save(licBytes)
}
