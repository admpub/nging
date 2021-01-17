package image

import (
	"io/ioutil"
	"mime/multipart"
	"sync"

	"github.com/admpub/errors"
	godl "github.com/admpub/go-download/v2"
	_ "golang.org/x/image/webp"
)

var (
	cachedWatermarkFileData  = sync.Map{}
	cachedWatermarkFileIndex = []string{}
	cachedWatermarkFileMax   = 10
)

func SetCachedWatermarkFileMax(n int) {
	cachedWatermarkFileMax = n
}

func DeleteCachedWatermarkFileData(keys ...string) {
	for _, key := range keys {
		cachedWatermarkFileData.Delete(key)
		for i, k := range cachedWatermarkFileIndex {
			if key != k {
				continue
			}
			endIndex := len(cachedWatermarkFileIndex) - 1
			if endIndex == i {
				cachedWatermarkFileIndex = cachedWatermarkFileIndex[0:endIndex]
				break
			}
			if i == 0 {
				cachedWatermarkFileIndex = cachedWatermarkFileIndex[1:]
			} else {
				cachedWatermarkFileIndex = append(cachedWatermarkFileIndex[0:i], cachedWatermarkFileIndex[i+1:]...)
			}
			break
		}
	}
}

func ClearCachedWatermarkFileData() {
	cachedWatermarkFileData.Range(func(key, _ interface{}) bool {
		cachedWatermarkFileData.Delete(key)
		return true
	})
	cachedWatermarkFileIndex = cachedWatermarkFileIndex[0:0]
}

func StoreCachedWatermarkFileData(key string, value interface{}) {
	if len(cachedWatermarkFileIndex)+1 > cachedWatermarkFileMax {
		DeleteCachedWatermarkFileData(cachedWatermarkFileIndex[len(cachedWatermarkFileIndex)-1])
	}
	cachedWatermarkFileData.Store(key, value)
	cachedWatermarkFileIndex = append(cachedWatermarkFileIndex, key)
}

func LoadCachedWatermarkFileData(key string) (value interface{}, ok bool) {
	return cachedWatermarkFileData.Load(key)
}

// GetRemoteWatermarkFileData 获取远程水印图片文件数据
func GetRemoteWatermarkFileData(fileURL string) (FileReader, error) {
	value, ok := LoadCachedWatermarkFileData(fileURL)
	if ok {
		if f, y := value.(*WatermarkData); y {
			return f.File, nil
		}
	}
	file, err := ReadRemoteWatermarkFile(fileURL)
	if err != nil {
		return file, err
	}
	StoreCachedWatermarkFileData(fileURL, NewWatermarkData(file))
	return file, err
}

// ReadRemoteWatermarkFile 读取远程水印图片文件
func ReadRemoteWatermarkFile(fileURL string) (multipart.File, error) {
	file, err := godl.Open(fileURL, DefaultWatermarkFileDownloadOptions)
	if err != nil {
		return nil, errors.WithMessage(err, fileURL)
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, errors.WithMessage(err, fileURL)
	}
	return Bytes2file(b), err
}
