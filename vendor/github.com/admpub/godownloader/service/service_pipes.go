package service

import (
	"context"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/admpub/godownloader/httpclient"
	"github.com/admpub/gohls-server/utils"
	"github.com/admpub/gohls/pkg"
)

func init() {
	label := `下载HLS视频`
	if utils.IsSupportedFFMPEG() {
		label += `(转换为mp4)`
	}
	PipeRegister(NewPipe(`dlhls`, label, func(d *httpclient.Downloader) error {
		ext := path.Ext(d.Fi.Url)
		if strings.ToLower(ext) != `.m3u8` {
			return nil
		}
		cfg := &pkg.Config{
			PlaylistURL: d.Fi.Url,
			OutputFile:  d.SafeFile().FilePath(),
		}
		p := strings.LastIndex(cfg.OutputFile, `.`)
		tsExt := `.ts`
		if p < 0 {
			cfg.OutputFile += tsExt
		} else {
			cfg.OutputFile = cfg.OutputFile[0:p] + tsExt
		}
		d.SetProgressGetter(func() (downloaded int64, total int64, percentProgress int64, speed int64) {
			prog := cfg.Progress()
			if prog.TotalNum > 0 {
				percentProgress = int64(prog.FinishedNum * 100 / prog.TotalNum)
			}

			return int64(prog.FinishedSize),
				-1, //int64(prog.TotalSize),//unknown total size
				percentProgress,
				int64(prog.SpeedInSecond)
		})
		var err error
		var done = make(chan struct{})
		ctx, cancelFunc := context.WithCancel(context.Background())
		go func() {
			err = d.SafeFile().ReOpen()
			if err == nil {
				err = cfg.Get(ctx, d.SafeFile().File)
				if err2 := d.SafeFile().Close(); err2 != nil {
					log.Println(err2)
				}
				if err == nil {
					tsFile := cfg.OutputFile
					mp4File := strings.TrimSuffix(tsFile, tsExt) + `.mp4`
					if err := utils.ConvertToMP4(tsFile, mp4File); err != nil {
						if !utils.IsUnsupported(err) {
							log.Println(`Conversion to mp4 file failed:`, err)
						}
					} else {
						if err := os.Remove(tsFile); err != nil {
							log.Println(`Deleting file "`+tsFile+`" failed:`, err)
						} else {
							d.SafeFile().SetFilePath(mp4File)
							cfg.OutputFile = mp4File
						}
					}
				}
			} else {
				log.Println(d.SafeFile().FilePath(), `reopen file failed:`, err)
			}
			done <- struct{}{}
		}()
		t := time.NewTicker(time.Second * 1)
		defer func() {
			t.Stop()
			log.Println(`m3u8 download completed:`, cfg.OutputFile)
		}()
		for {
			select {
			case <-done:
				d.Fi.Size = cfg.Progress().FinishedSize
				cancelFunc()
				return err
			case <-ctx.Done():
				log.Println(`m3u8 download canceled:`, cfg.OutputFile)
				return err
			case <-t.C:
				d.Fi.Size = cfg.Progress().FinishedSize
			}
		}
	}, `.m3u8`))
}

func NewPipe(ident string, label string, f func(*httpclient.Downloader) error, extensions ...string) *Pipe {
	return &Pipe{Ident: ident, Label: label, Extensions: extensions, function: f}
}

type Pipe struct {
	Ident      string
	Label      string
	Extensions []string
	function   func(*httpclient.Downloader) error
}

func (p *Pipe) Function() func(*httpclient.Downloader) error {
	return p.function
}

func (p *Pipe) SetFunction(f func(*httpclient.Downloader) error) *Pipe {
	p.function = f
	return p
}

var pipes = map[string]*Pipe{}

func PipeRegister(pipe *Pipe) {
	pipes[pipe.Ident] = pipe
}

func PipeList() map[string]*Pipe {
	return pipes
}

func PipeUnregister(ident string) {
	if _, ok := pipes[ident]; ok {
		delete(pipes, ident)
	}
}

func PipeGet(ident string) *Pipe {
	if pipe, ok := pipes[ident]; ok {
		return pipe
	}
	return nil
}

func PipeGetList(pipeNames ...string) []func(*httpclient.Downloader) error {
	pipes := []func(*httpclient.Downloader) error{}
	for _, pipeName := range pipeNames {
		pipe := PipeGet(pipeName)
		if pipe == nil {
			continue
		}
		pipes = append(pipes, pipe.function)
	}
	return pipes
}
