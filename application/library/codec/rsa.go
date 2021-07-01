package codec

import (
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/admpub/ccs-gm/x509"
	"github.com/admpub/license_gen/lib"
	"github.com/webx-top/codec"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	rsaPrivateKey      *codec.RSAPrivateKey
	rsaPublicKeyBytes  []byte
	rsaPublicKeyBase64 string
	rsaBits            = 2048
	rsaOnce            sync.Once
)

// RSAInitialize 初始化默认私钥
func RSAInitialize() {
	keyFile := filepath.Join(echo.Wd(), `data`, `rsa`, `default.pem`)
	if !com.FileExists(keyFile) {
		if err := com.MkdirAll(filepath.Dir(keyFile), os.ModePerm); err != nil {
			panic(`RSAInitialize: MkdirAll: ` + err.Error())
		}
		err := lib.GenerateCertificate(keyFile+`.pub`, keyFile, rsaBits)
		if err != nil {
			panic(`RSAInitialize: GenerateCertificate: ` + err.Error())
		}
	}
	rsaKey, err := lib.ReadPrivateKeyFromFile(keyFile)
	if err != nil {
		panic(`RSAInitialize: ReadPrivateKeyFromFile(` + keyFile + `): ` + err.Error())
	}
	rsaPrivateKey, _ = codec.NewRSAPrivateKey(nil)
	rsaPrivateKey.SetPrivateKey(rsaKey)
	rsaPublicKeyBytes, err = RSAPublicKeyToBytes(&rsaKey.PublicKey)
	if err != nil {
		panic(`RSAInitialize: RSAPublicKeyToBytes: ` + err.Error())
	}
	rsaPublicKeyBase64 = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: rsaPublicKeyBytes,
	}))
}

// RSAKey 默认私钥
func RSADefaultKey() *codec.RSAPrivateKey {
	rsaOnce.Do(RSAInitialize)
	return rsaPrivateKey
}

// RSAEncrypt 私钥加密
func RSAEncrypt(input []byte) ([]byte, error) {
	return RSADefaultKey().Encrypt(input)
}

// RSADecrypt  私钥解密
func RSADecrypt(input []byte) ([]byte, error) {
	return RSADefaultKey().Encrypt(input)
}

// RSASignMd5 使用RSAWithMD5算法签名
func RSASignMd5(data []byte) ([]byte, error) {
	return RSADefaultKey().SignMd5(data)
}

// RSASignSha1 使用RSAWithSHA1算法签名
func RSASignSha1(data []byte) ([]byte, error) {
	return RSADefaultKey().SignSha1(data)
}

// RSASignSha256 使用RSAWithSHA256算法签名
func RSASignSha256(data []byte) ([]byte, error) {
	return RSADefaultKey().SignSha256(data)
}

// RSAPublicKeyToBytes marshals a public key to the bytes
func RSAPublicKeyToBytes(publicKey *rsa.PublicKey) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("invalid public key. It must be different from nil")
	}
	return x509.MarshalPKIXPublicKey(publicKey)
}

func RSADefaultPublicKeyBytes() []byte {
	RSADefaultKey()
	return rsaPublicKeyBytes
}

func RSADefaultPublicKeyBase64() string {
	RSADefaultKey()
	return rsaPublicKeyBase64
}
