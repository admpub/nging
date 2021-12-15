package log

import "time"

// Entry represents a log entry.
type Entry struct {
	Level     Leveler
	Category  string
	Message   string
	Time      time.Time
	CallStack string

	FormattedMessage string
}

// String returns the string representation of the log entry
func (e *Entry) String() string {
	return e.FormattedMessage
}
