package codec

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/admpub/ccs-gm/sm2"
	"github.com/admpub/ccs-gm/sm3"
	"github.com/admpub/ccs-gm/utils"
	"github.com/admpub/ccs-gm/x509"
	"github.com/admpub/once"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var DefaultSM2 = NewSM2(`default`)

func NewSM2(name string) *SM2 {
	return &SM2{
		sm2Name: name,
	}
}

type SM2 struct {
	sm2Name               string
	defaultKey            *sm2.PrivateKey
	defaultPwd            []byte
	defaultPublicKeyBytes []byte
	defaultPublicKeyHex   string
	sm2once               once.Once
}

// Initialize 初始化默认私钥
func (s *SM2) init() {
	var err error
	keyFile := filepath.Join(echo.Wd(), `data`, `sm2`, s.sm2Name+`.pem`)
	if !com.FileExists(keyFile) {
		s.defaultKey, err = SM2GenKey()
		if err != nil {
			panic(`SM2GenKey: ` + err.Error())
		}
		if err = SaveKey(s.defaultKey, keyFile); err != nil {
			panic(err)
		}
	} else {
		s.defaultKey, err = ReadKey(keyFile)
		if err != nil {
			panic(err)
		}
	}
	err = s.initPublicKeyToMemory()
	if err != nil {
		panic(err)
	}
}

// SaveKey 保存私钥公钥
func (s *SM2) SaveKey(privateKey *sm2.PrivateKey, keyFile string, pwds ...[]byte) error {
	pwd := s.defaultPwd
	if len(pwds) > 0 {
		pwd = pwds[0]
	}
	// 保存私钥
	b, err := PrivateKeyToPEM(privateKey, pwd)
	if err != nil {
		return fmt.Errorf(`PrivateKeyToPEM: %w`, err)
	}
	if err := com.MkdirAll(filepath.Dir(keyFile), os.ModePerm); err != nil {
		return err
	}
	if err = os.WriteFile(keyFile, b, os.ModePerm); err != nil {
		return fmt.Errorf(`WriteFile `+keyFile+`: %w`, err)
	}
	os.Chmod(keyFile, os.ModePerm)
	// 保存公钥
	b, err = PublicKeyToPEM(&privateKey.PublicKey, pwd)
	if err != nil {
		return fmt.Errorf(`PublicKeyToPEM: %w`, err)
	}
	keyFile += `.pub`
	if err = os.WriteFile(keyFile, b, os.ModePerm); err != nil {
		return fmt.Errorf(`WriteFile `+keyFile+`: %w`, err)
	}
	os.Chmod(keyFile, os.ModePerm)
	return nil
}

// ReadKey 读取私钥公钥
func (s *SM2) ReadKey(keyFile string, pwds ...[]byte) (privateKey *sm2.PrivateKey, err error) {
	pwd := s.defaultPwd
	if len(pwds) > 0 {
		pwd = pwds[0]
	}
	var b []byte
	b, err = os.ReadFile(keyFile)
	if err != nil {
		err = fmt.Errorf(`ReadFile `+keyFile+`: %w`, err)
		return
	}
	privateKey, err = PEMtoPrivateKey(b, pwd)
	if err != nil {
		err = fmt.Errorf(`PEMtoPrivateKey: %w`, err)
		return
	}
	keyFile += `.pub`
	b, err = os.ReadFile(keyFile)
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

func (s *SM2) initPublicKeyToMemory() (err error) {
	s.defaultPublicKeyBytes, err = PublicKeyToBytes(&s.defaultKey.PublicKey)
	if err == nil {
		s.defaultPublicKeyHex = HexEncodeToString(s.defaultPublicKeyBytes)
	}
	return
}

// DefaultKey 默认私钥
func (s *SM2) DefaultKey() *sm2.PrivateKey {
	s.sm2once.Do(s.init)
	return s.defaultKey
}

// DefaultPublicKeyBytes 默认公钥
func (s *SM2) DefaultPublicKeyBytes() []byte {
	s.DefaultKey()
	return s.defaultPublicKeyBytes
}

// DefaultPublicKeyHex 默认公钥hex字符串
func (s *SM2) DefaultPublicKeyHex() string {
	s.DefaultKey()
	return s.defaultPublicKeyHex
}

// DefaultSM2DecryptHex 默认密钥解密hex字符串
func (s *SM2) DefaultDecryptHex(cipher string, noBase64 ...bool) (string, error) {
	return SM2DecryptHex(s.DefaultKey(), cipher, noBase64...)
}

func (s *SM2) Reset() {
	s.sm2once.Reset()
}

// ----------------

// SaveKey 保存私钥公钥
func SaveKey(privateKey *sm2.PrivateKey, keyFile string, pwds ...[]byte) error {
	return DefaultSM2.SaveKey(privateKey, keyFile, pwds...)
}

// ReadKey 读取私钥公钥
func ReadKey(keyFile string, pwds ...[]byte) (privateKey *sm2.PrivateKey, err error) {
	return DefaultSM2.ReadKey(keyFile, pwds...)
}

// DefaultKey 默认私钥
func DefaultKey() *sm2.PrivateKey {
	return DefaultSM2.DefaultKey()
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
		return nil, errors.New("invalid public key. It must be different from nil")
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
	base64 := true
	if len(noBase64) > 0 {
		base64 = !noBase64[0]
	} else if strings.HasPrefix(cipher, `-`) {
		base64 = false
		cipher = strings.TrimPrefix(cipher, `-`)
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
	if !base64 {
		return actual, nil
	}
	return com.Base64Decode(actual)
}

// DefaultPublicKeyBytes 默认公钥
func DefaultPublicKeyBytes() []byte {
	return DefaultSM2.DefaultPublicKeyBytes()
}

// DefaultPublicKeyHex 默认公钥hex字符串
func DefaultPublicKeyHex() string {
	return DefaultSM2.DefaultPublicKeyHex()
}

// DefaultSM2DecryptHex 默认密钥解密hex字符串
func DefaultSM2DecryptHex(cipher string, noBase64 ...bool) (string, error) {
	return DefaultSM2.DefaultDecryptHex(cipher, noBase64...)
}

func SM2Reset() {
	DefaultSM2.Reset()
}
