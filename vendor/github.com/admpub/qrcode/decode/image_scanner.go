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
	"errors"
	"image"
	"runtime"
)

// #include <zbar.h>
import "C"

// ImageScanner wraps zbar's internal image scanner type.  You can set
// configuration options on an ImageScanner (make sure to enable some
// or all symbologies first) and then use it to scan an Image for
// barcodes.
type ImageScanner struct {
	zbarScanner *C.zbar_image_scanner_t
}

// NewScanner creates a new ImageScanner and returns a pointer to it.
func NewScanner() *ImageScanner {
	scanner := &ImageScanner{
		zbarScanner: C.zbar_image_scanner_create(),
	}

	runtime.SetFinalizer(
		scanner,
		func(s *ImageScanner) {
			C.zbar_image_scanner_destroy(s.zbarScanner)
		},
	)

	return scanner
}

// ScanImage scans an Image and returns a slice of all Symbols found,
// or nil and an error if an error is encountered.
func (s *ImageScanner) ScanImage(img *Image) ([]*Symbol, error) {
	resultCode := C.zbar_scan_image(s.zbarScanner, img.zbarImage)

	if resultCode == 0 {
		return []*Symbol{}, nil
	} else if resultCode < 0 {
		// There doesn't seem to be an error code function for the
		// image scanner type
		return nil, errors.New("zbar: Error scanning image")
	}

	getFirst := func() *C.zbar_symbol_t {
		return C.zbar_image_first_symbol(img.zbarImage)
	}
	getNext := func(symbol *C.zbar_symbol_t) *C.zbar_symbol_t {
		return C.zbar_symbol_next(symbol)
	}

	symbols := []*Symbol{}
	for s := getFirst(); s != nil; s = getNext(s) {
		newSym := Symbol{
			Type:     SymbolType(C.zbar_symbol_get_type(s)),
			Data:     C.GoString(C.zbar_symbol_get_data(s)),
			Quality:  int(C.zbar_symbol_get_quality(s)),
			Boundary: []image.Point{},
		}

		for i := 0; i < int(C.zbar_symbol_get_loc_size(s)); i++ {
			newSym.Boundary = append(
				newSym.Boundary,
				image.Pt(
					int(C.zbar_symbol_get_loc_x(s, C.uint(i))),
					int(C.zbar_symbol_get_loc_y(s, C.uint(i))),
				),
			)
		}

		symbols = append(symbols, &newSym)
	}

	return symbols, nil
}

// SetEnabledAll turns the ImageScanner on or off for all symbologies.
func (s *ImageScanner) SetEnabledAll(enabled bool) *ImageScanner {
	return s.SetEnabledSymbology(None, enabled)
}

// SetEnabledSymbology turns the ImageScanner on or off for a specific
// SymbolType.
func (s *ImageScanner) SetEnabledSymbology(
	symbology SymbolType,
	enabled bool,
) *ImageScanner {
	s.setBooleanConfig(C.ZBAR_CFG_ENABLE, symbology, enabled)
	return s
}

// SetAddCheckAll enables or disables check digit when optional for
// all symbologies.
func (s *ImageScanner) SetAddCheckAll(enabled bool) *ImageScanner {
	return s.SetAddCheckSymbology(None, enabled)
}

// SetAddCheckSymbology enables or disables check digit when optional
// for a specific SymbolType.
func (s *ImageScanner) SetAddCheckSymbology(
	symbology SymbolType,
	enabled bool,
) *ImageScanner {
	s.setBooleanConfig(C.ZBAR_CFG_ADD_CHECK, symbology, enabled)
	return s
}

// SetEmitCheckAll enables or disables return of the check digit when
// present for all symbologies.
func (s *ImageScanner) SetEmitCheckAll(enabled bool) *ImageScanner {
	return s.SetEmitCheckSymbology(None, enabled)
}

// SetEmitCheckSymbology enables or disables return of the check digit
// when present for a specific SymbolType.
func (s *ImageScanner) SetEmitCheckSymbology(
	symbology SymbolType,
	enabled bool,
) *ImageScanner {
	s.setBooleanConfig(C.ZBAR_CFG_EMIT_CHECK, symbology, enabled)
	return s
}

// SetASCIIAll enables or disables the full ASCII character set for
// all symbologies.
func (s *ImageScanner) SetASCIIAll(enabled bool) *ImageScanner {
	return s.SetASCIISymbology(None, enabled)
}

// SetASCIISymbology enables or disables the full ASCII character set
// for a specific SymbolType.
func (s *ImageScanner) SetASCIISymbology(
	symbology SymbolType,
	enabled bool,
) *ImageScanner {
	s.setBooleanConfig(C.ZBAR_CFG_ASCII, symbology, enabled)
	return s
}

// SetMinLengthAll sets a minimum data length for all symbologies.
func (s *ImageScanner) SetMinLengthAll(length int) *ImageScanner {
	return s.SetMinLengthSymbology(None, length)
}

// SetMinLengthSymbology sets a minimum data length for a specific
// SymbolType.
func (s *ImageScanner) SetMinLengthSymbology(
	symbology SymbolType,
	length int,
) *ImageScanner {
	s.setIntConfig(C.ZBAR_CFG_MIN_LEN, symbology, length)
	return s
}

// SetMaxLengthAll sets a maximum data length for all symbologies.
func (s *ImageScanner) SetMaxLengthAll(length int) *ImageScanner {
	return s.SetMaxLengthSymbology(None, length)
}

// SetMaxLengthSymbology sets a maximum data length for a specific
// SymbolType.
func (s *ImageScanner) SetMaxLengthSymbology(
	symbology SymbolType,
	length int,
) *ImageScanner {
	s.setIntConfig(C.ZBAR_CFG_MAX_LEN, symbology, length)
	return s
}

// SetPositionEnabledAll enables or disables the collection of
// position data for all symbologies.
func (s *ImageScanner) SetPositionEnabledAll(enabled bool) *ImageScanner {
	return s.SetPositionEnabledSymbology(None, enabled)
}

// SetPositionEnabledSymbology enables or disables the collection of
// position data for a specific SymbolType.
func (s *ImageScanner) SetPositionEnabledSymbology(
	symbology SymbolType,
	enabled bool,
) *ImageScanner {
	s.setBooleanConfig(C.ZBAR_CFG_POSITION, symbology, enabled)
	return s
}

// SetXDensityAll sets the scanner's vertical scan density for all
// symbologies.
func (s *ImageScanner) SetXDensityAll(density int) *ImageScanner {
	return s.SetXDensitySymbology(None, density)
}

// SetXDensitySymbology sets the scanner's vertical scan density for a
// specific SymbolType.
func (s *ImageScanner) SetXDensitySymbology(
	symbology SymbolType,
	density int,
) *ImageScanner {
	s.setIntConfig(C.ZBAR_CFG_X_DENSITY, symbology, density)
	return s
}

// SetYDensityAll sets the scanner's horizontal scan density for all
// symbologies.
func (s *ImageScanner) SetYDensityAll(density int) *ImageScanner {
	return s.SetYDensitySymbology(None, density)
}

// SetYDensitySymbology sets the scanner's horizontal scan density for
// a specific SymbolType.
func (s *ImageScanner) SetYDensitySymbology(
	symbology SymbolType,
	density int,
) *ImageScanner {
	s.setIntConfig(C.ZBAR_CFG_Y_DENSITY, symbology, density)
	return s
}

func (s *ImageScanner) setBooleanConfig(
	option C.zbar_config_t,
	symbology SymbolType,
	enabled bool,
) {
	var e = 0
	if enabled {
		e = 1
	}

	C.zbar_image_scanner_set_config(
		s.zbarScanner,
		symbology.toEnum(),
		option,
		C.int(e),
	)
}

func (s *ImageScanner) setIntConfig(
	option C.zbar_config_t,
	symbology SymbolType,
	value int,
) {
	C.zbar_image_scanner_set_config(
		s.zbarScanner,
		symbology.toEnum(),
		option,
		C.int(value),
	)
}
