package lib

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/webx-top/com"
)

// License check errors
var (
	ErrorLicenseRead  = errors.New("Could not read license")
	ErrorPrivKeyRead  = errors.New("Could not read private key")
	ErrorPubKeyRead   = errors.New("Could not read public key")
	ErrorPrivKey      = errors.New("Invalid private key")
	ErrorPubKey       = errors.New("Invalid public key")
	ErrorMachineID    = errors.New("Could not read machine number")
	InvalidLicense    = errors.New("Invalid License file")
	UnlicensedVersion = errors.New("Unlicensed Version")
	InvalidMachineID  = errors.New("Invalid MachineID")
	InvalidLicenseID  = errors.New("Invalid LicenseID")
	ExpiredLicense    = errors.New("License expired")
)

type Validator interface {
	Validate(*LicenseData) error
}

// LicenseInfo - Core information about a license
type LicenseInfo struct {
	Name       string    `json:"name,omitempty"`
	LicenseID  string    `json:"licenseID,omitempty"`
	MachineID  string    `json:"machineID,omitempty"`
	Version    string    `json:"version,omitempty"`
	Expiration time.Time `json:"expiration"`
	Extra      Validator `json:"extra,omitempty"`
	validator  Validator
}

func (a *LicenseInfo) SetValidator(v Validator) {
	a.validator = v
}

func (a LicenseInfo) Remaining(langs ...string) *com.Durafmt {
	if a.Expiration.IsZero() {
		return nil
	}
	now := time.Now()
	duration := a.Expiration.Sub(now)
	//duration *= -1
	if len(langs) > 0 {
		return com.ParseDuration(duration, langs[0])
	}
	return com.ParseDuration(duration)
}

// LicenseData - This is the license data we serialise into a license file
type LicenseData struct {
	Info LicenseInfo `json:"info"`
	Key  string      `json:"key"`
}

// NewLicense from given info
func NewLicense(name string, expiry time.Time) *LicenseData {
	return &LicenseData{Info: LicenseInfo{Name: name, Expiration: expiry}}
}

func encodeKey(keyData []byte) string {
	return base64.StdEncoding.EncodeToString(keyData)
}

func decodeKey(keyStr string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(keyStr)
}

// Sign the License by updating the LicenseData.Key with given RSA private key
func (lic *LicenseData) Sign(pkey *rsa.PrivateKey) error {
	jsonLicInfo, err := json.Marshal(lic.Info)
	if err != nil {
		return err
	}

	signedData, err := Sign(pkey, jsonLicInfo)
	if err != nil {
		return err
	}

	lic.Key = encodeKey(signedData)

	return nil
}

// SignWithKey signs the License by updating the LicenseData.Key with given RSA
// private key read from a file
func (lic *LicenseData) SignWithKey(privKey string) error {
	rsaPrivKey, err := ReadPrivateKeyFromFile(privKey)
	if err != nil {
		return err
	}

	return lic.Sign(rsaPrivKey)
}

func (lic *LicenseData) ValidateLicenseKeyWithPublicKey(publicKey *rsa.PublicKey) error {
	signedData, err := decodeKey(lic.Key)
	if err != nil {
		return err
	}

	jsonLicInfo, err := json.Marshal(lic.Info)
	if err != nil {
		return err
	}

	return Unsign(publicKey, jsonLicInfo, signedData)
}

func (lic *LicenseData) ValidateLicenseKey(pubKey string) error {
	publicKey, err := ReadPublicKeyFromFile(pubKey)
	if err != nil {
		return err
	}

	return lic.ValidateLicenseKeyWithPublicKey(publicKey)
}

func (lic *LicenseData) CheckExpiration() error {
	if !lic.Info.Expiration.IsZero() && time.Now().After(lic.Info.Expiration) {
		return ExpiredLicense
	}
	return nil
}

func (lic *LicenseData) CheckVersion(versions ...string) error {
	if len(versions) > 0 && len(versions[0]) > 0 && len(lic.Info.Version) > 0 {
		if len(lic.Info.Version) > 1 {
			switch lic.Info.Version[0] {
			case '>':
				if len(lic.Info.Version) > 2 && lic.Info.Version[1] == '=' {
					if !com.VersionComparex(versions[0], lic.Info.Version[2:], `>=`) {
						return UnlicensedVersion
					}
					break
				}
				if !com.VersionComparex(versions[0], lic.Info.Version[1:], `>`) {
					return UnlicensedVersion
				}
			case '<':
				if len(lic.Info.Version) > 2 && lic.Info.Version[1] == '=' {
					if !com.VersionComparex(versions[0], lic.Info.Version[2:], `<=`) {
						return UnlicensedVersion
					}
					break
				}
				if !com.VersionComparex(versions[0], lic.Info.Version[1:], `<`) {
					return UnlicensedVersion
				}
			case '!':
				if len(lic.Info.Version) > 2 && lic.Info.Version[1] == '=' {
					if lic.Info.Version[2:] == versions[0] {
						return UnlicensedVersion
					}
					break
				}
				if lic.Info.Version[1:] == versions[0] {
					return UnlicensedVersion
				}
			default:
				if lic.Info.Version != versions[0] {
					return UnlicensedVersion
				}
			}
		} else {
			if lic.Info.Version != versions[0] {
				return UnlicensedVersion
			}
		}
	}
	return nil
}

func (lic *LicenseData) CheckMAC() error {
	addrs, err := MACAddresses(false)
	if err != nil {
		return err
	}
	var valid bool
	for _, addr := range addrs {
		if lic.Info.MachineID == Hash(addr) {
			valid = true
			break
		}
	}
	if !valid {
		return InvalidMachineID
	}
	return nil
}

func (lic *LicenseData) DefaultValidator(versions ...string) Validator {
	return &DefaultValidator{NowVersions: versions}
}

// CheckLicenseInfo checks license for logical errors such as for license expiry
func (lic *LicenseData) CheckLicenseInfo(versions ...string) error {
	if lic.Info.validator == nil {
		lic.Info.SetValidator(lic.DefaultValidator(versions...))
	}
	if err := lic.Info.validator.Validate(lic); err != nil {
		return err
	}
	if lic.Info.Extra != nil {
		return lic.Info.Extra.Validate(lic)
	}
	return nil
}

func (lic *LicenseData) WriteLicense(w io.Writer) error {
	jsonLic, err := json.MarshalIndent(lic, "", "  ")
	if err != nil {
		return err
	}

	_, werr := fmt.Fprintf(w, "%s", string(jsonLic))
	return werr
}

func (lic *LicenseData) SaveLicenseToFile(licName string) error {
	jsonLic, err := json.MarshalIndent(lic, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(licName, jsonLic, 0644)
}
