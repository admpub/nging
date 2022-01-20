package file

import (
	"log"
	"time"

	"github.com/webx-top/echo/middleware/session/engine"
)

var (
	DefaultInterval = time.Minute * 30
)

// Cleanup runs a background goroutine every interval that deletes expired
// sessions from the database.
//
// The design is based on https://github.com/yosssi/boltstore
func (m *filesystemStore) Cleanup(interval time.Duration, maxAge int) (chan<- struct{}, <-chan struct{}) {
	if interval <= 0 {
		interval = DefaultInterval
	}
	if maxAge <= 0 {
		maxAge = engine.DefaultMaxAge
	}

	quit, done := make(chan struct{}), make(chan struct{})
	go m.cleanup(interval, maxAge, quit, done)
	return quit, done
}

// StopCleanup stops the background cleanup from running.
func (m *filesystemStore) StopCleanup(quit chan<- struct{}, done <-chan struct{}) {
	quit <- struct{}{}
	<-done
}

// cleanup deletes expired sessions at set intervals.
func (m *filesystemStore) cleanup(interval time.Duration, maxAge int, quit <-chan struct{}, done chan<- struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-quit:
			// Handle the quit signal.
			done <- struct{}{}
			return
		case <-ticker.C:
			// Delete expired sessions on each tick.
			err := m.DeleteExpired(float64(maxAge))
			if err != nil {
				log.Printf("sessions: filesystem: unable to delete expired sessions: %v", err)
			}
		}
	}
}
