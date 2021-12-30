package monitor

import "errors"

var (
	ErrRunCompletedJob   = errors.New("error: try run completed job")
	ErrRunRunningJob     = errors.New("error: try run running job")
	ErrStopNonRunningJob = errors.New("error: imposible stop non running job")
)
