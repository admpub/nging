package httpclient

import (
	"log"
	"path/filepath"

	"github.com/admpub/godownloader/iotools"
	"github.com/admpub/godownloader/monitor"
)

type FileInfo struct {
	Size     int64    `json:"Size"`
	FileName string   `json:"FileName"`
	Url      string   `json:"Url"`
	Pipes    []string `json:"Pipes"`
}

type Downloader struct {
	sf             *iotools.SafeFile
	wp             *monitor.WorkerPool
	Fi             FileInfo
	pipes          []func(*Downloader) error
	progressGetter func() (downloaded int64, total int64, percentProgress int64, speed int64)
}

func (dl *Downloader) SafeFile() *iotools.SafeFile {
	return dl.sf
}

func (dl *Downloader) ProgressGetter() func() (downloaded int64, total int64, percentProgress int64, speed int64) {
	return dl.progressGetter
}

func (dl *Downloader) SetProgressGetter(f func() (downloaded int64, total int64, percentProgress int64, speed int64)) *Downloader {
	dl.progressGetter = f
	return dl
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

func CreateDownloader(url string, fp string, seg int64, getDown func() string, pipes ...func(*Downloader) error) (dl *Downloader, err error) {
	support, _ := CheckMultipart(url)
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
	if c > 0 {
		if err := sf.Truncate(c); err != nil {
			//can't truncate file
			return nil, err
		}
	}
	wp := new(monitor.WorkerPool)
	var dow monitor.DiscretWork
	if support {
		//create part-downloader foreach segment
		ps := c / seg
		for i := int64(0); i < seg-int64(1); i++ {
			d := CreatePartialDownloader(url, sf, ps*i, ps*i, ps*i+ps)
			mv := monitor.MonitoredWorker{Itw: d}
			wp.AppendWork(&mv)
		}
		lastseg := int64(ps * (seg - 1))
		dow = CreatePartialDownloader(url, sf, lastseg, lastseg, c)
	} else {
		dow = CreateDefaultDownloader(url, sf)
	}
	mv := monitor.MonitoredWorker{Itw: dow}

	//add to worker pool
	wp.AppendWork(&mv)
	d := &Downloader{
		sf:    sf,
		wp:    wp,
		Fi:    FileInfo{FileName: fp, Size: c, Url: url},
		pipes: pipes,
	}
	closeSafeFile(d)
	return d, nil
}

func RestoreDownloader(url string, fp string, dp []DownloadProgress, getDown func() string, pipes ...func(*Downloader) error) (dl *Downloader, err error) {
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
		sf:    sf,
		wp:    wp,
		Fi:    FileInfo{FileName: fp, Size: s.Size(), Url: url},
		pipes: pipes,
	}
	closeSafeFile(d)
	return d, nil
}

func closeSafeFile(d *Downloader) {
	d.wp.AfterComplete(func() {
		log.Println(`info: close file`, d.Fi.FileName)
		d.sf.Close()
		for _, pipe := range d.pipes {
			err := pipe(d)
			if err != nil {
				log.Println(`info: close file`, err)
			}
		}
	})
}
