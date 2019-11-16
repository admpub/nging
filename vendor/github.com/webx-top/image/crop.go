package image

import (
	"image"
	_ "image/gif" // register gif format

	"github.com/admpub/imageproxy"
	"github.com/admpub/imaging"
)

// default compression quality of resized jpegs
const defaultQuality = 95

// resample filter used when resizing images
var resampleFilter = imaging.Lanczos

func Crop(m image.Image, opt imageproxy.Options) image.Image {
	if opt.Rotate < 0 {
		opt.Rotate = int(360 + opt.Rotate)
	} else {
		opt.Rotate = int(opt.Rotate)
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

	tmpW := opt.CropWidth
	tmpH := opt.CropHeight
	dstW := opt.Width
	dstH := opt.Height
	srcX := opt.CropX
	srcY := opt.CropY
	m = imaging.Crop(m, image.Rect(int(srcX), int(srcY), int(srcX+tmpW-1), int(srcY+tmpH-1)))
	m = imaging.Thumbnail(m, int(dstW), int(dstH), resampleFilter)

	return m
}
