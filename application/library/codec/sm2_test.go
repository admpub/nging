package codec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/webx-top/echo/testing/test"
)

var keyFile string

func TestMain(m *testing.M) {
	DefaultSM2.keyFile = filepath.Join(os.Getenv("GOPATH"), "src/github.com/admpub/nging/data/sm2/default.pem")
	keyFile = DefaultSM2.keyFile
	DefaultSM2.DefaultKey()

	m.Run()
}

func TestDecrypt(t *testing.T) {
	pk, err := ReadKey(keyFile)
	if err != nil {
		panic(err)
	}

	excepted := `best`
	b, err := SM2Encrypt(&pk.PublicKey, []byte(excepted))
	if err != nil {
		panic(err)
	}
	plain, err := SM2Decrypt(pk, b)
	if err != nil {
		panic(err)
	}
	actual := string(plain)
	test.Eq(t, excepted, actual)
}

func TestDecryptHex(t *testing.T) {
	pk, err := ReadKey(keyFile)
	if err != nil {
		panic(err)
	}

	excepted := `best`
	crypted, err := SM2EncryptHex(&pk.PublicKey, excepted)
	if err != nil {
		panic(err)
	}
	plain, err := SM2DecryptHex(pk, crypted)
	if err != nil {
		panic(err)
	}
	test.Eq(t, excepted, plain)
}

func TestDefaultDecryptHex(t *testing.T) {
	t.Logf(`sm2 publicKey: %v`, DefaultPublicKeyHex())
	excepted := `123`
	crypted, err := DefaultSM2EncryptHex(excepted)
	if err != nil {
		panic(err)
	}
	t.Logf(`sm2 crypted: %v`, crypted)
	plain, err := DefaultSM2DecryptHex(crypted)
	if err != nil {
		panic(err)
	}
	test.Eq(t, excepted, plain)
}
