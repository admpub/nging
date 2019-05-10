package image

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"io/ioutil"
	"os"
)

//Base64ToFile base64 -> file
func Base64ToFile(base64Data string, file string) error {
	b, _ := base64.StdEncoding.DecodeString(base64Data)
	return ioutil.WriteFile(file, b, 0666)
}

//Base64ToBuffer base64 -> buffer
func Base64ToBuffer(base64Data string) *bytes.Buffer {
	b, _ := base64.StdEncoding.DecodeString(base64Data) //成图片文件并把文件写入到buffer
	return bytes.NewBuffer(b)                           // 必须加一个buffer 不然没有read方法就会报错
}

//转换成buffer之后里面就有Reader方法了。才能被图片API decode

//BufferToImageBuffer buffer-> ImageBuff（图片裁剪,代码接上面）
func BufferToImageBuffer(b *bytes.Buffer, x1 int, y1 int) *image.YCbCr {
	m, _, _ := image.Decode(b) // 图片文件解码
	rgbImg := m.(*image.YCbCr)
	subImg := rgbImg.SubImage(image.Rect(0, 0, x1, y1)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	return subImg
}

//ToFile img -> file(代码接上面)
func ToFile(subImg *image.YCbCr, file string) {
	f, _ := os.Create(file)     //创建文件
	defer f.Close()             //关闭文件
	jpeg.Encode(f, subImg, nil) //写入文件
}

//ToBase64 img -> base64(代码接上面)
func ToBase64(subImg *image.YCbCr) []byte {
	emptyBuff := bytes.NewBuffer(nil)                  //开辟一个新的空buff
	jpeg.Encode(emptyBuff, subImg, nil)                //img写入到buff
	dist := make([]byte, 50000)                        //开辟存储空间
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes()) //buff转成base64
	//fmt.Println(string(dist))                      //输出图片base64(type = []byte)
	//_ = ioutil.WriteFile("./base64pic.txt", dist, 0666) //buffer输出到jpg文件中（不做处理，直接写到文件）
	return dist
}

//FileToBase64 imgFile -> base64
func FileToBase64(srcFile string) []byte {
	ff, _ := ioutil.ReadFile(srcFile)   //我还是喜欢用这个快速读文件
	dist := make([]byte, 5000000)       //数据缓存
	base64.StdEncoding.Encode(dist, ff) // 文件转base64
	//_ = ioutil.WriteFile("./output2.jpg.txt", dist, 0666) //直接写入到文件就ok完活了。
	return dist
}
