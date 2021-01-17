package download

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/webx-top/com"
)

func Download(url, saveName string, options *Options) (int64, error) {
	f, err := Open(url, options)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}

	var output *os.File
	if len(saveName) == 0 {
		saveName = info.Name()
	} else {
		dir := filepath.Dir(saveName)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, os.ModePerm)
			if err != nil {
				return 0, err
			}
		}
	}
	output, err = os.OpenFile(saveName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return 0, err
	}
	defer output.Close()

	return io.Copy(output, f)
}

// NewHTTPClient New a client
func NewHTTPClient(options ...com.HTTPClientOptions) *http.Client {
	client := &http.Client{}
	for _, option := range options {
		option(client)
	}
	return client
}
