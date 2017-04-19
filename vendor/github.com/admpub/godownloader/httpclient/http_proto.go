package httpclient

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/admpub/godownloader/iotools"
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
	if resp.StatusCode != 206 {
		return false, errors.New("error: file not found or moved status: " + resp.Status)
	}
	if resp.ContentLength == 1 {
		log.Printf("info: multipart download support \n")
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
		log.Printf("error: file not found or moved status:", resp.StatusCode)
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

type PartialDownloader struct {
	dp     DownloadProgress
	client http.Client
	req    http.Response
	url    string
	file   *iotools.SafeFile
}

func CreatePartialDownloader(url string, file *iotools.SafeFile, from int64, pos int64, to int64) *PartialDownloader {
	var pd PartialDownloader
	pd.file = file
	pd.url = url
	pd.dp.From = from
	pd.dp.To = to
	pd.dp.Pos = pos
	return &pd
}

func (pd PartialDownloader) GetProgress() interface{} {
	return pd.dp
}

func (pd *PartialDownloader) BeforeDownload() error {
	//create new req
	r, err := http.NewRequest("GET", pd.url, nil)
	if err != nil {
		return err
	}

	r.Header.Add("Range", "bytes="+strconv.FormatInt(pd.dp.Pos, 10)+"-"+strconv.FormatInt(pd.dp.To, 10))
	/*
		f, _ := iotools.CreateSafeFile("test")
		r.Write(f)
		f.Close()
	*/
	resp, err := pd.client.Do(r)
	if err != nil {
		log.Printf("error: error download part file %v\n", err)
		return err
	}
	//check response
	if resp.StatusCode != 206 {
		log.Println("error: file not found or moved status:", resp.StatusCode)
		return errors.New("error: file not found or moved")
	}
	pd.req = *resp
	return nil
}

func (pd *PartialDownloader) AfterStopDownload() error {
	log.Println("info: try sync file")
	err := pd.file.Sync()
	pd.req.Body.Close()
	return err
}

func (pd *PartialDownloader) BeforeRun() error {
	return pd.BeforeDownload()
}

func (pd *PartialDownloader) AfterStop() error {
	return pd.AfterStopDownload()
}

func (pd *PartialDownloader) messureSpeed(realc int) {
	if time.Since(pd.dp.Lsmt).Seconds() > 0.5 {
		pd.dp.Speed = 2 * pd.dp.BytesInSecond
		pd.dp.Lsmt = time.Now()
		pd.dp.BytesInSecond = 0
	} else {
		pd.dp.BytesInSecond += int64(realc)
	}
}

func (pd *PartialDownloader) DownloadSergment() (bool, error) {
	//write flush data to disk
	buffer := make([]byte, FlushDiskSize, FlushDiskSize)

	count, err := pd.req.Body.Read(buffer)
	if (err != nil) && (err.Error() != "EOF") {
		pd.req.Body.Close()
		pd.file.Sync()
		return true, err
	}
	//log.Printf("returned from server %v bytes", count)
	if pd.dp.Pos+int64(count) > pd.dp.To {
		count = int(pd.dp.To - pd.dp.Pos)
		log.Printf("warning: server return to much for me i give only %v bytes", count)
	}

	realc, err := pd.file.WriteAt(buffer[:count], pd.dp.Pos)
	if err != nil {
		pd.file.Sync()
		pd.req.Body.Close()
		return true, err
	}
	pd.dp.Pos = pd.dp.Pos + int64(realc)
	pd.messureSpeed(realc)
	//log.Printf("writed %v pos %v to %v", realc, pd.dp.Pos, pd.dp.To)
	if pd.dp.Pos == pd.dp.To {
		//ok download part complete normal
		pd.file.Sync()
		pd.req.Body.Close()
		pd.dp.Speed = 0
		log.Printf("info: download complete normal")
		return true, nil
	}
	//not full download next segment
	return false, nil
}

func (pd *PartialDownloader) DoWork() (bool, error) {
	return pd.DownloadSergment()
}
