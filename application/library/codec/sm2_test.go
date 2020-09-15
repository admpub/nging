package codec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestDecrypt(t *testing.T) {
	keyFile := filepath.Join(os.Getenv("GOPATH"), "src/github.com/admpub/nging/data/sm2/default.pem")
	pk, err := ReadKey(keyFile)
	if err != nil {
		panic(err)
	}
	crypted := `047b7d20b23f8721c8d944a5e617707e474748dc6c3476942e075e42aa9024c2a3240182cfe003ee7c29b888657a0d178f4feadb5906d86191982246c0ca0d62fdd0716b6e7728845e5ed50df8ec6f6831d68ff80db0748192f5488138959691a120007bbb`
	actual, err := SM2DecryptHex(pk, crypted)
	if err != nil {
		panic(err)
	}
	test.Eq(t, `123`, actual)

	excepted := `best`
	b, err := SM2Encrypt(&pk.PublicKey, []byte(excepted))
	if err != nil {
		panic(err)
	}
	plain, err := SM2Decrypt(pk, b)
	if err != nil {
		panic(err)
	}
	actual = string(plain)
	test.Eq(t, excepted, actual)
}
