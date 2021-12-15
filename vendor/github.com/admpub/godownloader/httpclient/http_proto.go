package httpclient

import (
	"errors"
	"log"
	"net/http"
	"time"
)

const FlushDiskSize = 1024 * 1024

func CheckMultipart(urls string) (bool, error) {
	r, err := http.NewRequest("GET", urls, nil)
	if err != nil {
		return false, err
	}
	r.Header.Add("Range", "bytes=0-0")
	cl := http.Client{}
	resp, err := cl.Do(r)
	if err != nil {
		log.Printf("error: can't check multipart support assume no %v \n", err)
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 206 {
		if resp.StatusCode != 200 {
			err = errors.New("error: file not found or moved status: " + resp.Status)
		} else {
			log.Println("info: multipart download is unsupport")
		}

		return false, nil
	}
	if resp.ContentLength == 1 {
		log.Println("info: multipart download support")
		return true, nil
	}
	return false, nil
}

func GetSize(urls string) (int64, error) {
	cl := http.Client{}
	resp, err := cl.Head(urls)
	if err != nil {
		log.Printf("error: when try get file size %v \n", err)
		return 0, err
	}
	if resp.StatusCode != 200 {
		log.Println("error: file not found or moved status:", resp.StatusCode)
		return 0, errors.New("error: file not found or moved")
	}
	log.Printf("info: file size is %d bytes \n", resp.ContentLength)
	return resp.ContentLength, nil
}

type DownloadProgress struct {
	From          int64 `json:"From"`
	To            int64 `json:"To"`
	Pos           int64 `json:"Pos"`
	BytesInSecond int64
	Speed         int64
	Lsmt          time.Time
}
