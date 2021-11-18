/*
 * Copyright (c) 2015, Robert Bieber
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions
 * are met:
 *
 * 1. Redistributions of source code must retain the above copyright
 * notice, this list of conditions and the following disclaimer.
 *
 * 2. Redistributions in binary form must reproduce the above
 * copyright notice, this list of conditions and the following
 * disclaimer in the documentation and/or other materials provided
 * with the distribution.
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
 */

package decode

import (
	"image"
	"image/draw"
	"runtime"
	"unsafe"
)

// #cgo LDFLAGS: -lzbar
// #include <zbar.h>
import "C"

// TODO: Support multiple image formats, use type switch to select
// appropriate 4CC code
const y800 = 0x30303859

// Image wraps zbar's internal image type, and represents a
// 2-dimensional image suitable for scanning.
type Image struct {
	src       *image.Gray
	zbarImage *C.zbar_image_t
}

// NewImage creates a new zbar image by copying image data from an
// image.Image.
func NewImage(src image.Image) *Image {
	newImage := &Image{
		src:       image.NewGray(src.Bounds()),
		zbarImage: C.zbar_image_create(),
	}

	dims := newImage.src.Bounds().Size()
	C.zbar_image_set_size(newImage.zbarImage, C.uint(dims.X), C.uint(dims.Y))

	draw.Draw(
		newImage.src,
		newImage.src.Bounds(),
		src,
		image.Point{},
		draw.Over,
	)

	C.zbar_image_set_format(newImage.zbarImage, C.ulong(y800))
	C.zbar_image_set_data(
		newImage.zbarImage,
		unsafe.Pointer(&newImage.src.Pix[0]),
		C.ulong(len(newImage.src.Pix)),
		nil,
	)

	runtime.SetFinalizer(
		newImage,
		func(i *Image) {
			// The image data was allocated by the Go runtime, we
			// don't want zbar trying to free it
			C.zbar_image_set_data(newImage.zbarImage, nil, 0, nil)
			C.zbar_image_destroy(i.zbarImage)
		},
	)

	return newImage
}
