package shared

import "time"

// Defaults for sessions.Options
const (
	DefaultPath   = "/"
	DefaultMaxAge = 60 * 60 * 24 * 30 // 30days
)

// Defaults for store.Options
const (
	DefaultBucketName = "sessions"
)

// Defaults for reaper.Options
const (
	DefaultBatchSize     = 100
	DefaultCheckInterval = time.Minute
)
