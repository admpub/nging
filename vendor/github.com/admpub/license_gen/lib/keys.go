package lib

// Generating RSA keys
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

type RSAKeyData struct {
	Host       string
	ValidFrom  time.Time
	ValidFor   time.Duration
	IsCA       bool
	RSABits    int
	ECDSACurve string
}

// GenerateRSACertificate TODO: Find out why this method does not produce de-serializable keys (i.e.
// certificates produced by this method cannot be read back!)
func GenerateRSACertificate(certName string, keyName string, rsaKeyData *RSAKeyData) error {
	if len(rsaKeyData.Host) == 0 {
		return fmt.Errorf("Host parameter is required")
	}

	var priv interface{}
	var err error
	switch rsaKeyData.ECDSACurve {
	case "":
		priv, err = rsa.GenerateKey(rand.Reader, rsaKeyData.RSABits)
	case "P224":
		priv, err = ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case "P256":
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case "P384":
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case "P521":
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return fmt.Errorf("Unrecognized elliptic curve: %q", rsaKeyData.ECDSACurve)
	}
	if err != nil {
		return err
	}

	notBefore := rsaKeyData.ValidFrom
	notAfter := notBefore.Add(rsaKeyData.ValidFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return err
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Paraformance"},
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(rsaKeyData.Host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if rsaKeyData.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		return err
	}

	certOut, err := os.Create(certName)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", certName, err)
	}

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	certOut.Close()
	log.Printf("Written %s\n", certName)

	keyOut, err := os.OpenFile(keyName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", keyName, err)
	}
	pem.Encode(keyOut, pemBlockForKey(priv))
	keyOut.Close()
	log.Printf("Written %s\n", keyName)

	return nil
}

func GenerateCertificate(certName string, keyName string, rsaBits int) error {
	pubBytes, privBytes, err := GenerateCertificateData(rsaBits)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(certName, pubBytes, 0644); err != nil {
		return err
	}

	return ioutil.WriteFile(keyName, privBytes, 0644)
}

func GenerateCertificateData(rsaBits int) (pubBytes []byte, privBytes []byte, err error) {
	priv, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return nil, nil, err
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	pubBytes = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	privBytes = x509.MarshalPKCS1PrivateKey(priv)
	privBytes = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})

	return
}

func ReadPublicKey(r io.Reader) (*rsa.PublicKey, error) {
	keyBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, ErrorPubKey
	}

	pubkeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return nil, fmt.Errorf("Error parsing public key: %s", err)
	}

	pubkey, ok := pubkeyInterface.(*rsa.PublicKey)
	if !ok {
		log.Fatal("Fatal error")
	}

	return pubkey, nil
}

func ReadPublicKeyFromFile(key string) (*rsa.PublicKey, error) {
	file, err := os.Open(key)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	return ReadPublicKey(file)
}

func ReadPrivateKey(r io.Reader) (*rsa.PrivateKey, error) {
	keyBytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, ErrorPrivKey
	}
	privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privKey, nil
}

func ReadPrivateKeyFromFile(key string) (*rsa.PrivateKey, error) {
	file, err := os.Open(key)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	return ReadPrivateKey(file)
}
