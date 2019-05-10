/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

// Original source:
// https://code.google.com/p/zxing/source/browse/trunk/core/src/com/google/zxing/common/reedsolomon/ReedSolomonDecoder.java
//
// Copyright 2007 ZXing authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//      http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Implements Reed-Solomon decoding, as the name implies.
//
// The algorithm will not be explained here, but the following references were helpful
// in creating this implementation:
//
// - Bruce Maggs; "Decoding Reed-Solomon Codes" (see discussion of Forney's Formula)
//   http://www.cs.cmu.edu/afs/cs.cmu.edu/project/pscico-guyb/realworld/www/rs_decode.ps
// - J.I. Hall.; "Chapter 5. Generalized Reed-Solomon Codes" (see discussion of Euclidean algorithm)
//   www.mth.msu.edu/~jhall/classes/codenotes/GRS.pdf
//
// Much credit is due to William Rucklidge since portions of this code are an indirect
// port of his C++ Reed-Solomon implementation.
//
// Sean Owen, William Rucklidge

package rs

import (
	"errors"
	"fmt"
)

// Decoder can error correct data with the corresponding ECC codes.
type Decoder interface {
	// Decodes given set of received codewords, which include both data and
	// error-correction codewords. Really, this means it uses Reed-Solomon to
	// detect and correct errors, in-place, in the input.
	//
	// Returns the number of errors corrected or an error if decoding failed.
	Decode(data, ecc []byte) (int, error)
}

type rSDecoder struct {
	f *Field
}

// NewDecoder creates a decoder in the defined field.
func NewDecoder(f *Field) Decoder {
	return &rSDecoder{f}
}

func (d *rSDecoder) Decode(data, ecc []byte) (int, error) {
	// TODO(maruel): Temporary migration code.
	received := make([]byte, len(data)+len(ecc))
	copy(received, data)
	copy(received[len(data):], ecc)
	poly := &poly{d.f, received}
	syndromeCoeffs := make([]byte, len(ecc))
	noError := true
	for i := 0; i < len(ecc); i++ {
		eval := poly.evaluateAt(d.f.f.Exp(i))
		syndromeCoeffs[len(syndromeCoeffs)-1-i] = eval
		if eval != 0 {
			noError = false
		}
	}
	if noError {
		// Congrats! All the data was perfect.
		return 0, nil
	}
	// There was corruption found.
	syndrome := makePoly(d.f, syndromeCoeffs)
	sigma, omega, err := d.runEuclideanAlgorithm(buildMonomial(d.f, len(ecc), 1), syndrome, len(ecc))
	if err != nil {
		return 0, fmt.Errorf("runEuclidean() over %d bytes + %d ECC bytes failed: %s", len(data), len(ecc), err)
	}
	errorLocations := d.findErrorLocations(sigma)
	if errorLocations == nil {
		return 0, errors.New("Error locator degree does not match number of roots")
	}
	errorMagnitudes := d.findErrorMagnitudes(omega, errorLocations)
	for i := 0; i < len(errorLocations); i++ {
		position := len(received) - 1 - d.f.f.Log(errorLocations[i])
		if position < 0 {
			return 0, fmt.Errorf("Bad error location: %d", position)
		}
		// Calculate the original value.
		received[position] = d.f.f.Add(received[position], errorMagnitudes[i])
	}
	// Copy back.
	// TODO(maruel): Work in-place instead.
	copy(data, received)
	copy(ecc, received[len(data):])
	return len(errorLocations), nil
}

func (d *rSDecoder) runEuclideanAlgorithm(a, b *poly, R int) (sigma *poly, omega *poly, err error) {
	// Assume a's degree >= b's
	if a.degree() < b.degree() {
		a, b = b, a
	}
	rLast := a
	r := b
	tLast := getZero(d.f)
	t := getOne(d.f)

	// Run Euclidean algorithm until r's degree is less than R/2.
	for i := 0; r.degree() >= R/2; i++ {
		rLastLast := rLast
		tLastLast := tLast
		rLast = r
		tLast = t

		// Divide rLastLast by rLast, with quotient in q and remainder in r.
		if rLast.isZero() {
			// Oops, Euclidean algorithm already terminated?
			return nil, nil, fmt.Errorf("r_{i-1} was zero after %d iteration(s)", i)
		}
		r = rLastLast
		q := getZero(d.f)
		dltInverse := d.f.f.Inv(rLast.getCoefficient(rLast.degree()))
		for r.degree() >= rLast.degree() && !r.isZero() {
			// degreDiff is guaranteed to be >= 0
			degreeDiff := r.degree() - rLast.degree()
			scale := d.f.f.Mul(r.getCoefficient(r.degree()), dltInverse)
			q = q.add(buildMonomial(d.f, degreeDiff, scale))
			r = r.add(rLast.mulByMonomial(degreeDiff, scale))
		}
		t = q.mulPoly(tLast).add(tLastLast)
	}

	sigmaTildeAtZero := t.getCoefficient(0)
	if sigmaTildeAtZero == 0 {
		return nil, nil, errors.New("sigmaTilde(0) was zero")
	}

	inverse := d.f.f.Inv(sigmaTildeAtZero)
	return t.mulScalar(inverse), r.mulScalar(inverse), nil
}

// This is a direct application of Chien's search.
func (d *rSDecoder) findErrorLocations(errorLocator *poly) []byte {
	numErrors := errorLocator.degree()
	if numErrors == 1 {
		// Shortcut.
		return []byte{errorLocator.getCoefficient(1)}
	}
	result := make([]byte, numErrors)
	e := 0
	for i := 1; i < 256 && e < numErrors; i++ {
		if errorLocator.evaluateAt(byte(i)) == 0 {
			result[e] = d.f.f.Inv(byte(i))
			e++
		}
	}
	if e != numErrors {
		return nil
	}
	return result
}

// This is directly applying Forney's Formula.
func (d *rSDecoder) findErrorMagnitudes(errorEvaluator *poly, errorLocations []byte) []byte {
	s := len(errorLocations)
	result := make([]byte, s)
	for i := 0; i < s; i++ {
		xiInverse := d.f.f.Inv(errorLocations[i])
		denominator := byte(1)
		for j := 0; j < s; j++ {
			if i != j {
				denominator = d.f.f.Mul(denominator, d.f.f.Add(1, d.f.f.Mul(errorLocations[j], xiInverse)))
			}
		}
		result[i] = d.f.f.Mul(errorEvaluator.evaluateAt(xiInverse), d.f.f.Inv(denominator))
	}
	return result
}
