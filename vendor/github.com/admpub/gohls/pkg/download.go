package pkg

import (
	"context"
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
	progress     *Progress
}

func (cfg *Config) Get(ctx context.Context) error {
	return Get(ctx, cfg)
}

func (cfg *Config) Progress() *Progress {
	return cfg.progress
}

func Get(ctx context.Context, cfg *Config) error {
	cfg.progress = &Progress{}
	msChan := make(chan *Download, 1024)
	done := make(chan error)
	closeChan := func() {
		close(msChan)
		close(done)
	}
	var err error
	go func() {
		err = GetPlaylist(cfg.PlaylistURL, cfg.Duration, cfg.UseLocalTime, msChan, cfg.progress)
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
