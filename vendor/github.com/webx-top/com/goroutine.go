package com

import (
	"context"
	"errors"
	"log"
	"time"
)

var ErrExitedByContext = errors.New(`Received an exit notification from the context`)

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
