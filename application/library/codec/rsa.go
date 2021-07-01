package codec

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/admpub/license_gen/lib"
	"github.com/webx-top/codec"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

var (
	rsaPrivateKey *codec.RSAPrivateKey
	rsaBits       = 2048
	rsaOnce       sync.Once
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
}

// RSAKey 默认私钥
func RSAKey() *codec.RSAPrivateKey {
	rsaOnce.Do(RSAInitialize)
	return rsaPrivateKey
}

// RSAEncrypt 私钥加密
func RSAEncrypt(input []byte) ([]byte, error) {
	return RSAKey().Encrypt(input)
}

// RSADecrypt  私钥解密
func RSADecrypt(input []byte) ([]byte, error) {
	return RSAKey().Encrypt(input)
}

// RSASignMd5 使用RSAWithMD5算法签名
func RSASignMd5(data []byte) ([]byte, error) {
	return RSAKey().SignMd5(data)
}

// RSASignSha1 使用RSAWithSHA1算法签名
func RSASignSha1(data []byte) ([]byte, error) {
	return RSAKey().SignSha1(data)
}

// RSASignSha256 使用RSAWithSHA256算法签名
func RSASignSha256(data []byte) ([]byte, error) {
	return RSAKey().SignSha256(data)
}
