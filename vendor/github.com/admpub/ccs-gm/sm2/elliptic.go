// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package elliptic implements several standard elliptic curves over prime
// fields.
package sm2

import (
	"crypto/elliptic"
	"sync"
)

var initonce sync.Once

func initAll() {
	initP256()
}

// P256 returns a Curve which implements sm2 curve.
//
// The cryptographic operations are implemented using constant-time algorithms.
func P256() elliptic.Curve {
	initonce.Do(initAll)
	return p256
}