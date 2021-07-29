package com

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"time"
)

var ErrExitedByContext = errors.New(`received an exit notification from the context`)

func Loop(ctx context.Context, exec func() error, duration time.Duration) error {
	if ctx == nil {
		ctx = context.Background()
	}
	check := func() <-chan struct{} {
		return ctx.Done()
	}
	for {
		select {
		case <-check():
			log.Println(CalledAtFileLine(2), ErrExitedByContext)
			return ErrExitedByContext
		default:
			if err := exec(); err != nil {
				return err
			}
			time.Sleep(duration)
		}
	}
}

// Notify 等待系统信号
// <-Notify()
func Notify(sig ...os.Signal) chan os.Signal {
	terminate := make(chan os.Signal, 1)
	if len(sig) == 0 {
		sig = []os.Signal{os.Interrupt}
	}
	signal.Notify(terminate, sig...)
	return terminate
}
