package log

import "io"

// Target represents a target where the logger can send log messages to for further processing.
type Target interface {
	// Open prepares the target for processing log messages.
	// Open will be invoked when Logger.Open() is called.
	// If an error is returned, the target will be removed from the logger.
	// errWriter should be used to write errors found while processing log messages.
	Open(errWriter io.Writer) error
	// Process processes an incoming log message.
	Process(*Entry)
	// Close closes a target.
	// Close is called when Logger.Close() is called, which gives each target
	// a chance to flush the logged messages to their destination storage.
	Close()
	SetLevel(interface{})
	SetLevels(...Leveler)
}
