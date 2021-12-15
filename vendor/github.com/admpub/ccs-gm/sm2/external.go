// Copyright 2020 cetc-30. All rights reserved.
// SPDX-License-Identifier: Apache-2.0
// license that can be found in the LICENSE file.

package sm2

import (
	"crypto"
	"crypto/rand"
	"encoding/asn1"
	"io"
	"math/big"
)

type Sm2PrivateKey struct {
	D *big.Int //sk
}

type Sm2PublicKey struct {
	X *big.Int //pk.X
	Y *big.Int //pk.Y
}

type sm2Signature struct {
	R, S *big.Int
}

func (priv *PrivateKey) Sign(rand io.Reader, msg []byte, opt crypto.SignerOpts) ([]byte, error) {
	r, s, err := Sign(rand, priv, msg)
	if err != nil {
		return nil, err
	}
	return asn1.Marshal(sm2Signature{r, s})
}

func (pub *PublicKey) Verify(msg []byte, sign []byte) bool {
	var sm2Sign sm2Signature
	_, err := asn1.Unmarshal(sign, &sm2Sign)
	if err != nil {
		return false
	}
	return Verify(pub, msg, sm2Sign.R, sm2Sign.S)
}

func Sm2KeyGen(rand io.Reader) (sk, pk []byte, err error) {
	priv, _ := GenerateKey(rand)
	var sm2SK Sm2PrivateKey
	var sm2PK Sm2PublicKey

	sm2SK.D = priv.D
	sm2PK.X = priv.X
	sm2PK.Y = priv.Y

	sk, _ = asn1.Marshal(sm2SK)
	pk, _ = asn1.Marshal(sm2PK)
	return
}

func Sm2Sign(sk, pk, msg []byte) ([]byte, error) {
	var sm2SK Sm2PrivateKey
	var sm2PK Sm2PublicKey
	_, err := asn1.Unmarshal(sk, &sm2SK)
	if err != nil {
		return nil, err
	}

	_, err = asn1.Unmarshal(pk, &sm2PK)
	if err != nil {
		return nil, err
	}

	var priv PrivateKey
	priv.Curve = P256()
	priv.D = sm2SK.D
	priv.X = sm2PK.X
	priv.Y = sm2PK.Y

	r, s, err := Sign(rand.Reader, &priv, msg)
	if err != nil {
		return nil, err
	}

	return asn1.Marshal(sm2Signature{r, s})
}

func Sm2Verify(sign, pk, msg []byte) bool {
	var sm2Sign sm2Signature
	var sm2PK Sm2PublicKey

	_, err := asn1.Unmarshal(sign, &sm2Sign)
	if err != nil {
		return false
	}

	_, err = asn1.Unmarshal(pk, &sm2PK)
	if err != nil {
		return false
	}

	var PK PublicKey
	PK.Curve = P256()
	PK.X = sm2PK.X
	PK.Y = sm2PK.Y

	return PK.Verify(msg, sign)
}
