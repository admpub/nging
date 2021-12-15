/*
 * Copyright (c) 2015, Robert Bieber
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 *    notice, this list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above
 *    copyright notice, this list of conditions and the following
 *    disclaimer in the documentation and/or other materials provided
 *    with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
 * LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
 * FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
 * COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
 * INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package decode

import (
	"image"
)

// #cgo LDFLAGS: -lzbar
// #include <zbar.h>
import "C"

type SymbolType int

const (
	None    SymbolType = C.ZBAR_NONE
	Partial SymbolType = C.ZBAR_PARTIAL
	EAN8    SymbolType = C.ZBAR_EAN8
	UPCE    SymbolType = C.ZBAR_UPCE
	ISBN10  SymbolType = C.ZBAR_ISBN10
	UPCA    SymbolType = C.ZBAR_UPCA
	EAN13   SymbolType = C.ZBAR_EAN13
	ISBN13  SymbolType = C.ZBAR_ISBN13
	I25     SymbolType = C.ZBAR_I25
	Code39  SymbolType = C.ZBAR_CODE39
	PDF417  SymbolType = C.ZBAR_PDF417
	QRCode  SymbolType = C.ZBAR_QRCODE
	Code128 SymbolType = C.ZBAR_CODE128
)

// Name returns the name of a given symbol encoding type, or "UNKNOWN"
// if the encoding is not recognized.
func (s SymbolType) Name() string {
	return C.GoString(C.zbar_get_symbol_name(s.toEnum()))
}

// Quick conversion function to deal with C functions that want an
// enum type.
func (s SymbolType) toEnum() C.zbar_symbol_type_t {
	return C.zbar_symbol_type_t(s)
}

// Symbol represents a scanned barcode.
type Symbol struct {
	Type SymbolType
	Data string

	// Quality is an unscaled, relative quantity which expresses the
	// confidence of the match.  These values are currently meaningful
	// only in relation to each other: a larger value is more
	// confident than a smaller one.
	Quality int

	// Boundary is a set of image.Point which define a polygon
	// containing the scanned symbol in the image.
	Boundary []image.Point
}
