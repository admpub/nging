package charset

import (
	"github.com/webx-top/chardet"
)

var textDetector = chardet.NewTextDetector(chardet.GB18030, chardet.UTF8)

func DetectText(content []byte) (string, error) {
	r, err := textDetector.DetectBest(content)
	if err != nil {
		return ``, err
	}
	return r.Charset, nil
}
