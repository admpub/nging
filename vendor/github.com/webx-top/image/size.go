package image

const (
	inch2cm = 2.54
	cm2inch = 0.393700787
)

//Inch2Cm 英寸转厘米
func Inch2Cm(inch float32) float32 {
	return inch / cm2inch
}

//Cm2Inch 厘米转英寸
func Cm2Inch(cm float32) float32 {
	return cm / inch2cm
}

//Cm2Pix 厘米转像素
func Cm2Pix(cm float32, DPI int) int {
	return int(Cm2Inch(cm) * float32(DPI))
}

//Pix2Cm 像素转厘米
func Pix2Cm(pix float32, DPI int) float32 {
	return Inch2Cm(pix / float32(DPI))
}
