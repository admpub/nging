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

	progress *Progress
	reader   io.Reader
}

func (cfg *Config) Get(ctx context.Context, reader ...io.Reader) error {
	return Get(ctx, cfg, reader...)
}

func (cfg *Config) Progress() *Progress {
	return cfg.progress
}

func Get(ctx context.Context, cfg *Config, reader ...io.Reader) error {
	cfg.progress = &Progress{}
	msChan := make(chan *Download, 1024)
	done := make(chan error)
	closeChan := func() {
		close(msChan)
		close(done)
	}
	var err error
	go func() {
		if len(reader) > 0 && reader[0] != nil {
			c, err := NewContext(cfg.PlaylistURL, 1024)
			if err != nil {
				log.Println(err)
				closeChan()
				return
			}
			err = cfg.GetPlaylistFromReader(c, reader[0], msChan)
		} else {
			err = cfg.GetPlaylist(cfg.PlaylistURL, msChan)
		}
		if err != nil {
			log.Println(err)
			closeChan()
		}
	}()
	go func() {
		done <- DownloadSegment(cfg.OutputFile, msChan, cfg.Duration, cfg.progress)
	}()
	for {
		select {
		case <-ctx.Done():
			closeChan()
			return nil
		case err = <-done:
			return err
		}
	}
}
