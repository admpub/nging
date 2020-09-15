package utils

import (
	"crypto/rand"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/admpub/ccs-gm/sm2"
	"github.com/admpub/ccs-gm/x509"
)

// PrivateKeyToPEM converts the private key to PEM format.
// EC private keys are converted to PKCS#8 format.
func PrivateKeyToPEM(privateKey *sm2.PrivateKey, pwd []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("Invalid sm2 private key. It must be different from nil.")
	}
	raw, err := x509.MarshalECPrivateKey(privateKey)

	if err != nil {
		return nil, err
	}

	var block *pem.Block
	if len(pwd) > 0 {
		block, err = x509.EncryptPEMBlock(
			rand.Reader,
			"PRIVATE KEY",
			raw,
			pwd,
			x509.PEMCipherAES256)
	} else {
		block = &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: raw,
		}
	}

	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(block), nil
}

// PEMtoPrivateKey unmarshals a pem to private key
func PEMtoPrivateKey(raw []byte, pwd []byte) (*sm2.PrivateKey, error) {
	if len(raw) == 0 {
		return nil, errors.New("Invalid PEM. It must be different from nil.")
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("Failed decoding PEM. Block must be different from nil. [% x]", raw)
	}

	if x509.IsEncryptedPEMBlock(block) {
		if len(pwd) == 0 {
			return nil, errors.New("Encrypted Key. Need a password")
		}

		decrypted, err := x509.DecryptPEMBlock(block, pwd)
		if err != nil {
			return nil, fmt.Errorf("Failed PEM decryption [%s]", err)
		}

		key, err := x509.ParseECPrivateKey(decrypted)
		if err != nil {
			return nil, err
		}
		sm2Key, ok := key.(*sm2.PrivateKey)
		if ok {
			return sm2Key, nil
		} else {
			return nil, errors.New("key type error")
		}
	}

	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	sm2Key, ok := key.(*sm2.PrivateKey)
	if ok {
		return sm2Key, nil
	} else {
		return nil, errors.New("key type error")
	}
}

// PublicKeyToPEM marshals a public key to the pem format
func PublicKeyToPEM(publicKey *sm2.PublicKey, pwd []byte) ([]byte, error) {
	if len(pwd) != 0 {
		return PublicKeyToEncryptedPEM(publicKey, pwd)
	}

	if publicKey == nil {
		return nil, errors.New("Invalid public key. It must be different from nil.")
	}

	PubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: PubASN1,
		},
	), nil

}

// PublicKeyToEncryptedPEM converts a public key to encrypted pem
func PublicKeyToEncryptedPEM(publicKey *sm2.PublicKey, pwd []byte) ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("Invalid public key. It must be different from nil.")
	}
	if len(pwd) == 0 {
		return nil, errors.New("Invalid password. It must be different from nil.")
	}

	raw, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	block, err := x509.EncryptPEMBlock(
		rand.Reader,
		"PUBLIC KEY",
		raw,
		pwd,
		x509.PEMCipherAES256)
	if err != nil {
		return nil, err
	}

	return pem.EncodeToMemory(block), nil
}

// PEMtoPublicKey unmarshals a pem to public key
func PEMtoPublicKey(raw []byte, pwd []byte) (*sm2.PublicKey, error) {
	if len(raw) == 0 {
		return nil, errors.New("Invalid PEM. It must be different from nil.")
	}
	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, fmt.Errorf("Failed decoding. Block must be different from nil. [% x]", raw)
	}

	// TODO: derive from header the type of the key
	if x509.IsEncryptedPEMBlock(block) {
		if len(pwd) == 0 {
			return nil, errors.New("Encrypted Key. Password must be different from nil")
		}

		decrypted, err := x509.DecryptPEMBlock(block, pwd)
		if err != nil {
			return nil, fmt.Errorf("Failed PEM decryption. [%s]", err)
		}

		key, err := x509.ParsePKIXPublicKey(decrypted)
		if err != nil {
			return nil, err
		}
		sm2Pk, ok := key.(*sm2.PublicKey)
		if ok {
			return sm2Pk, nil
		} else {
			return nil, errors.New("invalid public key format")
		}
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	sm2Pk, ok := key.(*sm2.PublicKey)
	if ok {
		return sm2Pk, nil
	} else {
		return nil, errors.New("invalid public key format")
	}
}
