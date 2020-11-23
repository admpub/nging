package codec

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/admpub/ccs-gm/sm2"
	"github.com/admpub/ccs-gm/sm3"
	"github.com/admpub/ccs-gm/utils"
	"github.com/admpub/ccs-gm/x509"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	defaultKey *sm2.PrivateKey
	defaultPwd []byte
	once       sync.Once
)

// Initialize 初始化默认私钥
func Initialize() {
	var err error
	keyFile := filepath.Join(echo.Wd(), `data`, `sm2`, `default.pem`)
	if !com.FileExists(keyFile) {
		defaultKey, err = SM2GenKey()
		if err != nil {
			panic(`SM2GenKey: ` + err.Error())
		}
		if err = SaveKey(defaultKey, keyFile); err != nil {
			panic(err)
		}
	} else {
		defaultKey, err = ReadKey(keyFile)
		if err != nil {
			panic(err)
		}
	}
}

// SaveKey 保存私钥公钥
func SaveKey(privateKey *sm2.PrivateKey, keyFile string, pwds ...[]byte) error {
	pwd := defaultPwd
	if len(pwds) > 0 {
		pwd = pwds[0]
	}
	// 保存私钥
	b, err := PrivateKeyToPEM(privateKey, pwd)
	if err != nil {
		return fmt.Errorf(`PrivateKeyToPEM: %w`, err)
	}
	os.MkdirAll(filepath.Dir(keyFile), os.ModePerm)
	if err = ioutil.WriteFile(keyFile, b, os.ModePerm); err != nil {
		return fmt.Errorf(`WriteFile `+keyFile+`: %w`, err)
	}
	// 保存公钥
	b, err = PublicKeyToPEM(&privateKey.PublicKey, pwd)
	if err != nil {
		return fmt.Errorf(`PublicKeyToPEM: %w`, err)
	}
	keyFile += `.pub`
	if err = ioutil.WriteFile(keyFile, b, os.ModePerm); err != nil {
		return fmt.Errorf(`WriteFile `+keyFile+`: %w`, err)
	}
	return nil
}

// ReadKey 读取私钥公钥
func ReadKey(keyFile string, pwds ...[]byte) (privateKey *sm2.PrivateKey, err error) {
	pwd := defaultPwd
	if len(pwds) > 0 {
		pwd = pwds[0]
	}
	var b []byte
	b, err = ioutil.ReadFile(keyFile)
	if err != nil {
		err = fmt.Errorf(`ReadFile `+keyFile+`: %w`, err)
		return
	}
	privateKey, err = PEMtoPrivateKey(b, pwd)
	if err != nil {
		err = fmt.Errorf(`PEMtoPrivateKey: %w`, err)
	}
	keyFile += `.pub`
	b, err = ioutil.ReadFile(keyFile)
	if err != nil {
		err = fmt.Errorf(`ReadFile `+keyFile+`: %w`, err)
		return
	}
	var publickKey *sm2.PublicKey
	publickKey, err = PEMtoPublicKey(b, pwd)
	if err != nil {
		err = fmt.Errorf(`PEMtoPrivateKey: %w`, err)
	} else {
		privateKey.PublicKey = *publickKey
	}
	return
}

// DefaultKey 默认私钥
func DefaultKey() *sm2.PrivateKey {
	once.Do(Initialize)
	return defaultKey
}

// SM2GenKey 生成私钥和公钥
func SM2GenKey() (*sm2.PrivateKey, error) {
	return sm2.GenerateKey(rand.Reader)
}

// SM2VerifySign 验签
func SM2VerifySign(priv *sm2.PrivateKey, msg []byte) bool {
	hash := sm3.SumSM3(msg)
	r, s, err := sm2.Sign(rand.Reader, priv, hash[:])
	if err != nil {
		return false
	}
	return sm2.Verify(&priv.PublicKey, hash[:], r, s)
}

// SM2Encrypt 加密
func SM2Encrypt(pubKey *sm2.PublicKey, msg []byte) ([]byte, error) {
	return sm2.Encrypt(rand.Reader, pubKey, msg)
}

// SM2Decrypt 解密
func SM2Decrypt(priv *sm2.PrivateKey, cipher []byte) ([]byte, error) {
	return sm2.Decrypt(cipher, priv)
}

// PrivateKeyToPEM 私钥转PEM文件内容
func PrivateKeyToPEM(privateKey *sm2.PrivateKey, pwd []byte) ([]byte, error) {
	return utils.PrivateKeyToPEM(privateKey, pwd)
}

// PublicKeyToPEM 公钥转PEM文件内容
func PublicKeyToPEM(publickKey *sm2.PublicKey, pwd []byte) ([]byte, error) {
	return utils.PublicKeyToPEM(publickKey, pwd)
}

// PEMtoPrivateKey PEM文件内容转私钥
func PEMtoPrivateKey(raw []byte, pwd []byte) (*sm2.PrivateKey, error) {
	return utils.PEMtoPrivateKey(raw, pwd)
}

// PEMtoPublicKey PEM文件内容转公钥
func PEMtoPublicKey(raw []byte, pwd []byte) (*sm2.PublicKey, error) {
	return utils.PEMtoPublicKey(raw, pwd)
}

// HexEncodeToString hex编码为字符串
func HexEncodeToString(b []byte) string {
	return hex.EncodeToString(b)
}

// PublicKeyToBytes marshals a public key to the bytes
func PublicKeyToBytes(publicKey *sm2.PublicKey) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("Invalid public key. It must be different from nil")
	}
	return x509.MarshalPKIXPublicKey(publicKey)
}

// PublicKeyToHexString 公钥转hex字符串
func PublicKeyToHexString(publicKey *sm2.PublicKey) string {
	pubASN1, err := PublicKeyToBytes(publicKey)
	if err != nil {
		return err.Error()
	}
	return HexEncodeToString(pubASN1)
}

// SM2DecryptHex 解密
func SM2DecryptHex(priv *sm2.PrivateKey, cipher string, noBase64 ...bool) (string, error) {
	if len(cipher) == 0 {
		return ``, nil
	}
	b, err := hex.DecodeString(cipher)
	if err != nil {
		return ``, err
	}
	plain, err := SM2Decrypt(priv, b)
	if err != nil {
		return ``, err
	}
	actual := string(plain)
	if len(noBase64) > 0 && noBase64[0] {
		return actual, nil
	}
	return com.Base64Decode(actual)
}

// DefaultPublicKeyHex 默认公钥hex字符串
func DefaultPublicKeyHex() string {
	return PublicKeyToHexString(&DefaultKey().PublicKey)
}

// DefaultSM2DecryptHex 默认密钥解密hex字符串
func DefaultSM2DecryptHex(cipher string, noBase64 ...bool) (string, error) {
	return SM2DecryptHex(DefaultKey(), cipher, noBase64...)
}
