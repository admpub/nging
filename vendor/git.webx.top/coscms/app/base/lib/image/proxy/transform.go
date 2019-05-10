// Copyright 2013 Google Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package imageproxy

import (
	"bytes"
	"image"
	_ "image/gif" // register gif format
	"image/jpeg"
	"image/png"
	"io"

	"github.com/disintegration/imaging"
)

// default compression quality of resized jpegs
const defaultQuality = 95

// resample filter used when resizing images
var resampleFilter = imaging.Lanczos

// Transform the provided image.  img should contain the raw bytes of an
// encoded image in one of the supported formats (gif, jpeg, or png).  The
// bytes of a similarly encoded image is returned.
func Transform(img []byte, opt Options) ([]byte, error) {
	if opt == emptyOptions {
		// bail if no transformation was requested
		return img, nil
	}
	return TransformFromReader(bytes.NewReader(img), opt)
}

func TransformFromReader(r io.Reader, opt Options) ([]byte, error) {
	if opt == emptyOptions {
		// bail if no transformation was requested
		return []byte{}, nil
	}

	// decode image
	m, format, err := image.Decode(r)
	if err != nil {
		return nil, err
	}

	// transform and encode image
	buf := new(bytes.Buffer)
	switch format {
	case "gif":
		fn := func(img image.Image) image.Image {
			return transformImage(img, opt)
		}
		err = GifProcess(buf, r, fn)
		if err != nil {
			return nil, err
		}
	case "jpeg":
		quality := opt.Quality
		if quality == 0 {
			quality = defaultQuality
		}

		m = transformImage(m, opt)
		err = jpeg.Encode(buf, m, &jpeg.Options{Quality: quality})
		if err != nil {
			return nil, err
		}
	case "png":
		m = transformImage(m, opt)
		err = png.Encode(buf, m)
		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// transformImage modifies the image m based on the transformations specified
// in opt.
func transformImage(m image.Image, opt Options) image.Image {
	if opt.CropOptions != nil {
		return CropImage(m, opt)
	}
	// convert percentage width and height values to absolute values
	imgW := m.Bounds().Max.X - m.Bounds().Min.X
	imgH := m.Bounds().Max.Y - m.Bounds().Min.Y
	var w, h int
	if 0 < opt.Width && opt.Width < 1 {
		w = int(float64(imgW) * opt.Width)
	} else if opt.Width < 0 {
		w = 0
	} else {
		w = int(opt.Width)
	}
	if 0 < opt.Height && opt.Height < 1 {
		h = int(float64(imgH) * opt.Height)
	} else if opt.Height < 0 {
		h = 0
	} else {
		h = int(opt.Height)
	}

	// never resize larger than the original image
	if !opt.ScaleUp {
		if w > imgW {
			w = imgW
		}
		if h > imgH {
			h = imgH
		}
	}

	// resize
	if w != 0 || h != 0 {
		if opt.Fit {
			m = imaging.Fit(m, w, h, resampleFilter)
		} else {
			if w == 0 || h == 0 {
				m = imaging.Resize(m, w, h, resampleFilter)
			} else {
				m = imaging.Thumbnail(m, w, h, resampleFilter)
			}
		}
	}

	// flip
	if opt.FlipVertical {
		m = imaging.FlipV(m)
	}
	if opt.FlipHorizontal {
		m = imaging.FlipH(m)
	}

	// rotate
	switch opt.Rotate {
	case 90:
		m = imaging.Rotate90(m)
	case 180:
		m = imaging.Rotate180(m)
	case 270:
		m = imaging.Rotate270(m)
	}

	return m
}

func CropImage(m image.Image, opt Options) image.Image {
	if opt.CropOptions.Rotate < 0 {
		opt.Rotate = int(360 + opt.CropOptions.Rotate)
	} else {
		opt.Rotate = int(opt.CropOptions.Rotate)
	}
	// rotate
	switch opt.Rotate {
	case 90:
		m = imaging.Rotate90(m)
	case 180:
		m = imaging.Rotate180(m)
	case 270:
		m = imaging.Rotate270(m)
	}

	tmp_img_w := opt.CropOptions.Width
	tmp_img_h := opt.CropOptions.Height
	dst_img_w := opt.Width
	dst_img_h := opt.Height
	src_x := opt.CropOptions.X
	src_y := opt.CropOptions.Y
	m = imaging.Crop(m, image.Rect(int(src_x), int(src_y), int(src_x+tmp_img_w-1), int(src_y+tmp_img_h-1)))
	m = imaging.Thumbnail(m, int(dst_img_w), int(dst_img_h), resampleFilter)

	return m
}
