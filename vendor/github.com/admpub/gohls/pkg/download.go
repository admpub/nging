package pkg

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"
)

type Progress struct {
	TotalNum      int
	FinishedNum   int
	TotalSize     int64
	FinishedSize  int64
	SpeedInSecond float64
}

func (p *Progress) ResetExecute() {
	p.FinishedNum = 0
	p.FinishedSize = 0
	p.SpeedInSecond = 0
}

func (p *Progress) JSONBytes() []byte {
	b, _ := json.Marshal(p)
	return b
}

func (p *Progress) FromJSONBytes(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return json.Unmarshal(b, p)
}

type Config struct {
	PlaylistURL     string
	OutputFile      string
	Duration        time.Duration
	UseLocalTime    bool
	ForceRedownload bool

	progress *Progress
}

func (cfg *Config) Get(ctx context.Context, reader ...io.Reader) error {
	return Get(ctx, cfg, reader...)
}

func (cfg *Config) Progress() *Progress {
	return cfg.progress
}

func (cfg *Config) RemoveHistory() {
	os.Remove(cfg.OutputFile + `._prog_`)
	os.Remove(cfg.OutputFile + `._seg_`)
}

func Get(ctx context.Context, cfg *Config, reader ...io.Reader) error {
	cfg.progress = &Progress{}
	if !cfg.ForceRedownload {
		if b, err := os.ReadFile(cfg.OutputFile + `._prog_`); err == nil {
			cfg.progress.FromJSONBytes(b)
		}
	}
	msChan := make(chan *Download, 1024)
	defer func() {
		defer func() {
			recover()
		}()
		close(msChan)
	}()
	go func() {
		if !cfg.ForceRedownload {
			if b, err := os.ReadFile(cfg.OutputFile + `._seg_`); err == nil {
				dll := []*Download{}
				err = json.Unmarshal(b, &dll)
				if err == nil {
					defer func() {
						recover()
					}()
					for _, d := range dll {
						copyD := *d
						msChan <- &copyD
					}
					return
				}
			}
		}
		var err error
		if len(reader) > 0 && reader[0] != nil {
			var c *Context
			c, err = NewContext(cfg.PlaylistURL, 1024)
			if err != nil {
				log.Println(err)
				return
			}
			err = cfg.GetPlaylistFromReader(c, reader[0], msChan)
		} else {
			err = cfg.GetPlaylist(cfg.PlaylistURL, msChan)
		}
		if err != nil {
			log.Println(err)
		}
	}()
	err := DownloadSegment(ctx, cfg, msChan)
	if err == nil {
		cfg.RemoveHistory()
	}
	return err
}
