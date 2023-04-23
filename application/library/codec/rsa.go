package codec

import (
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"

	"github.com/admpub/ccs-gm/x509"
	"github.com/admpub/license_gen/lib"
	"github.com/admpub/once"
	"github.com/webx-top/codec"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var DefaultRSA = NewRSA(`default`)

func NewRSA(name string, bits ...int) *RSA {
	var rsaBits int
	if len(bits) > 0 {
		rsaBits = bits[0]
	}
	if rsaBits <= 0 {
		rsaBits = 2048
	}
	return &RSA{
		rsaBits: rsaBits,
		rsaName: name,
	}
}

type RSA struct {
	defaultKey      *codec.RSA
	publicKeyBytes  []byte
	publicKeyBase64 string
	rsaBits         int
	rsaName         string
	rsaOnce         once.Once
	keyFile         string
}

func (r *RSA) init() {
	if len(r.keyFile) == 0 {
		r.keyFile = filepath.Join(echo.Wd(), `data`, `rsa`, r.rsaName+`.pem`)
	}
	if !com.FileExists(r.keyFile) {
		if err := com.MkdirAll(filepath.Dir(r.keyFile), os.ModePerm); err != nil {
			panic(`RSAInitialize: MkdirAll: ` + err.Error())
		}
		err := lib.GenerateCertificate(r.keyFile+`.pub`, r.keyFile, r.rsaBits)
		if err != nil {
			panic(`RSAInitialize: GenerateCertificate: ` + err.Error())
		}
	}
	rsaKey, err := lib.ReadPrivateKeyFromFile(r.keyFile)
	if err != nil {
		panic(`RSAInitialize: ReadPrivateKeyFromFile(` + r.keyFile + `): ` + err.Error())
	}
	rsaPrivateKey, _ := codec.NewRSAPrivateKey(nil)
	rsaPrivateKey.SetPrivateKey(rsaKey)
	r.publicKeyBytes, err = RSAPublicKeyToBytes(&rsaKey.PublicKey)
	if err != nil {
		panic(`RSAInitialize: RSAPublicKeyToBytes: ` + err.Error())
	}
	r.publicKeyBase64 = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: r.publicKeyBytes,
	}))
	rsaPublicKey, _ := codec.NewRSAPublicKey(nil)
	rsaPublicKey.SetPublicKey(&rsaKey.PublicKey)
	r.defaultKey = codec.NewRSA()
	r.defaultKey.SetPrivateKey(rsaPrivateKey).SetPublicKey(rsaPublicKey)
}

// DefaultKey 默认私钥
func (r *RSA) DefaultKey() *codec.RSA {
	r.rsaOnce.Do(r.init)
	return r.defaultKey
}

// Encrypt 私钥加密
func (r *RSA) Encrypt(input []byte) ([]byte, error) {
	return r.DefaultKey().PublicKey().Encrypt(input)
}

// Decrypt  私钥解密
func (r *RSA) Decrypt(input []byte) ([]byte, error) {
	return r.DefaultKey().PrivateKey().Decrypt(input)
}

// SignMd5 使用RSAWithMD5算法签名
func (r *RSA) SignMd5(data []byte) ([]byte, error) {
	return r.DefaultKey().PrivateKey().SignMd5(data)
}

// SignSha1 使用RSAWithSHA1算法签名
func (r *RSA) SignSha1(data []byte) ([]byte, error) {
	return r.DefaultKey().PrivateKey().SignSha1(data)
}

// SignSha256 使用RSAWithSHA256算法签名
func (r *RSA) SignSha256(data []byte) ([]byte, error) {
	return r.DefaultKey().PrivateKey().SignSha256(data)
}

func (r *RSA) DefaultPublicKeyBytes() []byte {
	r.DefaultKey()
	return r.publicKeyBytes
}

func (r *RSA) DefaultPublicKeyBase64() string {
	r.DefaultKey()
	return r.publicKeyBase64
}

func (r *RSA) Reset() {
	r.rsaOnce.Reset()
}

// ----------------

// RSAKey 默认私钥
func RSADefaultKey() *codec.RSA {
	return DefaultRSA.DefaultKey()
}

// RSAEncrypt 私钥加密
func RSAEncrypt(input []byte) ([]byte, error) {
	return DefaultRSA.Encrypt(input)
}

// RSADecrypt  私钥解密
func RSADecrypt(input []byte) ([]byte, error) {
	return DefaultRSA.Decrypt(input)
}

// RSASignMd5 使用RSAWithMD5算法签名
func RSASignMd5(data []byte) ([]byte, error) {
	return DefaultRSA.SignMd5(data)
}

// RSASignSha1 使用RSAWithSHA1算法签名
func RSASignSha1(data []byte) ([]byte, error) {
	return DefaultRSA.SignSha1(data)
}

// RSASignSha256 使用RSAWithSHA256算法签名
func RSASignSha256(data []byte) ([]byte, error) {
	return DefaultRSA.SignSha256(data)
}

// RSAPublicKeyToBytes marshals a public key to the bytes
func RSAPublicKeyToBytes(publicKey *rsa.PublicKey) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("invalid public key. It must be different from nil")
	}
	return x509.MarshalPKIXPublicKey(publicKey)
}

func RSADefaultPublicKeyBytes() []byte {
	return DefaultRSA.DefaultPublicKeyBytes()
}

func RSADefaultPublicKeyBase64() string {
	return DefaultRSA.DefaultPublicKeyBase64()
}

func RSAReset() {
	DefaultRSA.Reset()
}
