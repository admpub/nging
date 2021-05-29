/* Copyright 2012 Marc-Antoine Ruel. Licensed under the Apache License, Version
2.0 (the "License"); you may not use this file except in compliance with the
License.  You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0. Unless required by applicable law or
agreed to in writing, software distributed under the License is distributed on
an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
or implied. See the License for the specific language governing permissions and
limitations under the License. */

// Original sources:
// https://code.google.com/p/zxing/source/browse/trunk/core/src/com/google/zxing/common/reedsolomon/GenericGFPoly.java
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
// Represents a polynomial whose coefficients are elements of a GF.
// Instances of this class are immutable.
//
// Much credit is due to William Rucklidge since portions of this code are an indirect
// port of his C++ Reed-Solomon implementation.
//
// @author Sean Owen

package rs

import (
	"fmt"
)

type poly struct {
	field        *Field
	coefficients []byte // In reverse order.
}

var zero = []byte{0}
var one = []byte{1}

// |coefficients| representing elements of GF(size), arranged from most
// significant (highest-power term) coefficient to least significant.
func makePoly(field *Field, coefficients []byte) *poly {
	if len(coefficients) == 0 {
		return nil
	}
	obj := &poly{field: field}
	if len(coefficients) > 1 && coefficients[0] == 0 {
		// Leading term must be non-zero for anything except the constant polynomial "0".
		firstNonZero := 1
		for coefficients[firstNonZero] == 0 && firstNonZero < len(coefficients) {
			firstNonZero++
		}
		if firstNonZero == len(coefficients) {
			obj.coefficients = zero
		} else {
			// Slice it.
			obj.coefficients = coefficients[firstNonZero:]
		}
	} else {
		obj.coefficients = coefficients
	}
	return obj
}

func getZero(field *Field) *poly {
	return &poly{field, zero}
}

func getOne(field *Field) *poly {
	return &poly{field, one}
}

func (p *poly) degree() int {
	return len(p.coefficients) - 1
}

// Returns the monomial representing coefficient * x^degree.
func buildMonomial(field *Field, degree int, coefficient byte) *poly {
	if degree < 0 {
		return nil
	}
	if coefficient == 0 {
		return getZero(field)
	}
	coefficients := make([]byte, degree+1)
	coefficients[0] = coefficient
	return &poly{field, coefficients}
}

// Returns true iff this polynomial is the monomial "0".
func (p *poly) isZero() bool {
	return p.coefficients[0] == 0
}

// Returns coefficient of x^degree term in this polynomial.
func (p *poly) getCoefficient(degree int) byte {
	return p.coefficients[len(p.coefficients)-1-degree]
}

// Returns evaluation of this polynomial at a given point.
func (p *poly) evaluateAt(a byte) byte {
	if a == 0 {
		// Just return the x^0 coefficient
		return p.getCoefficient(0)
	}
	if a == 1 {
		// Just the sum of the coefficients.
		result := byte(0)
		for _, v := range p.coefficients {
			result = p.field.f.Add(result, v)
		}
		return result
	}
	result := p.coefficients[0]
	for i := 1; i < len(p.coefficients); i++ {
		result = p.field.f.Add(p.field.f.Mul(a, result), p.coefficients[i])
	}
	return result
}

func (p *poly) add(other *poly) *poly {
	if p.isZero() {
		return other
	}
	if other.isZero() {
		return p
	}
	smaller := p.coefficients
	larger := other.coefficients
	if len(smaller) > len(larger) {
		smaller, larger = larger, smaller
	}
	sumDiff := make([]byte, len(larger))
	lengthDiff := len(larger) - len(smaller)
	// Copy high-order terms only found in higher-degree polynomial's coefficients
	copy(sumDiff, larger[:lengthDiff])

	for i := lengthDiff; i < len(larger); i++ {
		sumDiff[i] = p.field.f.Add(smaller[i-lengthDiff], larger[i])
	}
	return makePoly(p.field, sumDiff)
}

func (p *poly) mulPoly(other *poly) *poly {
	if p.isZero() || other.isZero() {
		return getZero(p.field)
	}
	aCoefficients := p.coefficients
	bCoefficients := other.coefficients
	product := make([]byte, len(aCoefficients)+len(bCoefficients)-1)
	for i := 0; i < len(aCoefficients); i++ {
		aCoeff := aCoefficients[i]
		for j := 0; j < len(bCoefficients); j++ {
			product[i+j] = p.field.f.Add(product[i+j], p.field.f.Mul(aCoeff, bCoefficients[j]))
		}
	}
	return makePoly(p.field, product)
}

func (p *poly) mulScalar(scalar byte) *poly {
	if scalar == 0 {
		return getZero(p.field)
	}
	if scalar == 1 {
		return p
	}
	product := make([]byte, len(p.coefficients))
	for i := 0; i < len(p.coefficients); i++ {
		product[i] = p.field.f.Mul(p.coefficients[i], scalar)
	}
	return makePoly(p.field, product)
}

func (p *poly) mulByMonomial(degree int, coefficient byte) *poly {
	if degree < 0 {
		return nil
	}
	if coefficient == 0 {
		return getZero(p.field)
	}
	size := len(p.coefficients)
	product := make([]byte, size+degree)
	for i := 0; i < size; i++ {
		product[i] = p.field.f.Mul(p.coefficients[i], coefficient)
	}
	return makePoly(p.field, product)
}

func (p *poly) divide(divisor *poly) (q *poly, r *poly) {
	if divisor.isZero() {
		// "Divide by 0".
		return nil, nil
	}
	quotient := getZero(p.field)
	remainder := p

	denominatorLeadingTerm := divisor.getCoefficient(divisor.degree())
	inverseDenominatorLeadingTerm := p.field.f.Inv(denominatorLeadingTerm)
	for remainder.degree() >= divisor.degree() && !remainder.isZero() {
		degreeDifference := remainder.degree() - divisor.degree()
		scale := p.field.f.Mul(remainder.getCoefficient(remainder.degree()), inverseDenominatorLeadingTerm)
		term := divisor.mulByMonomial(degreeDifference, scale)
		iterationQuotient := buildMonomial(p.field, degreeDifference, scale)
		quotient = quotient.add(iterationQuotient)
		remainder = remainder.add(term)
	}
	return quotient, remainder
}

func (p *poly) String() string {
	return fmt.Sprintf("poly{%v}", p.coefficients)
}
