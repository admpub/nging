package httpclient

import (
	"log"
	"path/filepath"

	"github.com/admpub/godownloader/iotools"
	"github.com/admpub/godownloader/monitor"
)

type FileInfo struct {
	Size     int64  `json:"Size"`
	FileName string `json:"FileName"`
	Url      string `json:"Url"`
}

type Downloader struct {
	sf *iotools.SafeFile
	wp *monitor.WorkerPool
	Fi FileInfo
}

func (dl *Downloader) StopAll() []error {
	defer dl.sf.Close()
	return dl.wp.StopAll()
}

func (dl *Downloader) StartAll() []error {
	if err := dl.sf.ReOpen(); err != nil {
		return []error{err}
	}
	return dl.wp.StartAll()
}

func (dl *Downloader) GetProgress() []DownloadProgress {
	pr := dl.wp.GetAllProgress().([]interface{})
	re := make([]DownloadProgress, len(pr))
	for i, val := range pr {
		re[i] = val.(DownloadProgress)
	}
	return re
}

func (dl *Downloader) State() monitor.State {
	return dl.wp.State()
}

func CreateDownloader(url string, fp string, seg int64, getDown func() string) (dl *Downloader, err error) {
	c, err := GetSize(url)
	if err != nil {
		//can't get file size
		return nil, err
	}
	dfs := getDown() + fp
	dfs = filepath.Clean(dfs)
	sf, err := iotools.CreateSafeFile(dfs)
	if err != nil {
		//can't create file on path
		return nil, err
	}
	defer sf.Close()
	if err := sf.Truncate(c); err != nil {
		//can't truncate file
		return nil, err
	}
	//create part-downloader foreach segment
	ps := c / seg
	wp := new(monitor.WorkerPool)
	for i := int64(0); i < seg-int64(1); i++ {
		d := CreatePartialDownloader(url, sf, ps*i, ps*i, ps*i+ps)
		mv := monitor.MonitoredWorker{Itw: d}
		wp.AppendWork(&mv)
	}
	lastseg := int64(ps * (seg - 1))
	dow := CreatePartialDownloader(url, sf, lastseg, lastseg, c)
	mv := monitor.MonitoredWorker{Itw: dow}

	//add to worker pool
	wp.AppendWork(&mv)
	d := &Downloader{
		sf: sf,
		wp: wp,
		Fi: FileInfo{FileName: fp, Size: c, Url: url},
	}
	closeSafeFile(d)
	return d, nil
}

func RestoreDownloader(url string, fp string, dp []DownloadProgress, getDown func() string) (dl *Downloader, err error) {
	dfs := getDown() + fp
	sf, err := iotools.OpenSafeFile(dfs)
	if err != nil {
		//can't create file on path
		return nil, err
	}
	defer sf.Close()
	s, err := sf.Stat()
	if err != nil {
		return nil, err
	}
	wp := new(monitor.WorkerPool)
	for _, r := range dp {
		dow := CreatePartialDownloader(url, sf, r.From, r.Pos, r.To)
		mv := monitor.MonitoredWorker{Itw: dow}

		//add to worker pool
		wp.AppendWork(&mv)

	}
	d := &Downloader{
		sf: sf,
		wp: wp,
		Fi: FileInfo{FileName: fp, Size: s.Size(), Url: url},
	}
	closeSafeFile(d)
	return d, nil
}

func closeSafeFile(d *Downloader) {
	d.wp.AfterComplete(func() {
		log.Println(`info: close file`, d.Fi.FileName)
		d.sf.Close()
	})
}
