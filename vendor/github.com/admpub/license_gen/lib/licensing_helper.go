package lib

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"

	"github.com/webx-top/com"
)

func ReadLicense(r io.Reader) (*LicenseData, error) {
	ldata, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var license LicenseData
	if err := json.Unmarshal(ldata, &license); err != nil {
		return nil, err
	}

	return &license, nil
}

func ReadLicenseFromFile(licFile string) (*LicenseData, error) {
	file, err := os.Open(licFile)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	return ReadLicense(file)
}

func SignData(privKey string, data []byte) (string, error) {
	rsaPrivKey, err := ReadPrivateKey(strings.NewReader(privKey))
	if err != nil {
		fmt.Println("Error reading private key:", err)
		return ``, err
	}

	signedData, err := Sign(rsaPrivKey, data)
	if err != nil {
		fmt.Println("Error signing data:", err)
		return ``, err
	}

	return encodeKey(signedData), nil
}

// Sign signs data with rsa-sha256
func Sign(r *rsa.PrivateKey, data []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, r, crypto.SHA256, d)
}

func UnsignData(pubKey string, signature string, data []byte) error {
	publicKey, err := ReadPublicKey(strings.NewReader(pubKey))
	if err != nil {
		return err
	}

	signedData, err := decodeKey(signature)
	if err != nil {
		return err
	}

	return Unsign(publicKey, data, signedData)
}

// Unsign verifies the message using a rsa-sha256 signature
func Unsign(r *rsa.PublicKey, message []byte, sig []byte) error {
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	return rsa.VerifyPKCS1v15(r, crypto.SHA256, d, sig)
}

func MACAddresses(encoded bool) ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	hardwareAddrs := make([]string, 0)
	for _, inter := range interfaces {
		macAddr := fmt.Sprint(inter.HardwareAddr)
		if len(macAddr) == 0 {
			continue
		}
		if encoded {
			hardwareAddrs = append(hardwareAddrs, fmt.Sprintf(`%x`, macAddr))
			continue
		}
		hardwareAddrs = append(hardwareAddrs, macAddr)
	}
	return hardwareAddrs, err
}

func CheckAndReturningLicense(licenseReader, pubKeyReader io.Reader, validator Validator, versions ...string) (*LicenseData, error) {
	lic, err := ReadLicense(licenseReader)
	if err != nil {
		return lic, ErrorLicenseRead
	}

	publicKey, err := ReadPublicKey(pubKeyReader)
	if err != nil {
		return lic, ErrorPubKeyRead
	}

	if err := lic.ValidateLicenseKeyWithPublicKey(publicKey); err != nil {
		return lic, InvalidLicense // we have a key mismatch here meaning license data is tampered
	}
	lic.Info.SetValidator(validator)
	return lic, lic.CheckLicenseInfo(versions...)
}

// CheckLicense reads a license from licenseReader and then validate it against the
// public key read from pubKeyReader
func CheckLicense(licenseReader, pubKeyReader io.Reader, validator Validator, versions ...string) error {
	_, err := CheckAndReturningLicense(licenseReader, pubKeyReader, validator, versions...)
	return err
}

func CheckLicenseStringAndReturning(license, pubKey string, validator Validator, versions ...string) (*LicenseData, error) {
	return CheckAndReturningLicense(strings.NewReader(license), strings.NewReader(pubKey), validator, versions...)
}

// CheckLicenseString 检测授权文件是否有效
// license 为授权文件内容
// pubKey 为公钥内容
func CheckLicenseString(license, pubKey string, validator Validator, versions ...string) error {
	return CheckLicense(strings.NewReader(license), strings.NewReader(pubKey), validator, versions...)
}

func Hash(raw string) string {
	return strings.ToUpper(com.Hash(fmt.Sprintf(`%x`, raw)))
}

// GenerateLicense 生成授权文件内容
// privKey 为私钥内容
func GenerateLicense(info *LicenseInfo, privKey string) ([]byte, error) {
	data, err := BuildLicense(info, privKey)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(data, "", "  ")
}

func BuildLicense(info *LicenseInfo, privKey string) (*LicenseData, error) {
	if len(info.MachineID) == 0 {
		addrs, err := MACAddresses(true)
		if err != nil {
			return nil, err
		}
		if len(addrs) < 1 {
			return nil, ErrorMachineID
		}
		info.MachineID = strings.ToUpper(com.Hash(addrs[0]))
	}
	data := &LicenseData{
		Info: *info,
	}
	rsaPrivKey, err := ReadPrivateKey(strings.NewReader(privKey))
	if err != nil {
		return nil, err
	}
	err = data.Sign(rsaPrivKey)
	if err != nil {
		return nil, err
	}
	return data, err
}
