package pkg

import (
	"context"
	"io"
	"log"
	"time"
)

type Progress struct {
	TotalNum      int
	FinishedNum   int
	TotalSize     int64
	FinishedSize  int64
	SpeedInSecond float64
}

type Config struct {
	PlaylistURL  string
	OutputFile   string
	Duration     time.Duration
	UseLocalTime bool
	MaxRetries   int

	progress *Progress
}

func (cfg *Config) Get(ctx context.Context, reader ...io.Reader) error {
	return Get(ctx, cfg, reader...)
}

func (cfg *Config) Progress() *Progress {
	return cfg.progress
}

func Get(ctx context.Context, cfg *Config, reader ...io.Reader) error {
	cfg.progress = &Progress{}
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 5
	}
	msChan := make(chan *Download, 1024)
	defer func() {
		defer func() {
			recover()
		}()
		close(msChan)
	}()
	go func() {
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
	return DownloadSegment(ctx, cfg, msChan)
}
