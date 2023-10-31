package htpasswd

import (
	"crypto/sha1"
	"encoding/base64"

	"github.com/GehirnInc/crypt/apr1_crypt"
	"golang.org/x/crypto/bcrypt"
)

func HashApr1(password string) (string, error) {
	return apr1_crypt.New().Generate([]byte(password), nil)
}

func HashBCrypt(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ``, err
	}
	return string(passwordBytes), nil
}

func HashSha(password string) string {
	s := sha1.New()
	s.Write([]byte(password))
	passwordSum := []byte(s.Sum(nil))
	return base64.StdEncoding.EncodeToString(passwordSum)
}

type Algorithm string

const (
	// AlgoAPR1 Apache MD5 crypt - legacy
	AlgoAPR1 Algorithm = "apr1"
	// AlgoBCrypt bcrypt - recommended
	AlgoBCrypt Algorithm = "bcrypt"
	// AlgoSHA sha5 insecure - do not use
	AlgoSHA Algorithm = "sha"
)
