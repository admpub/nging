package server

import (
	"fmt"
	"log"
)

// Use an instance of this to log in a standard format
type Logger struct {
	sessionID string
}

func newLogger(id string) *Logger {
	l := new(Logger)
	l.sessionID = id
	return l
}

func (logger *Logger) Print(message interface{}) {
	log.Printf("%s   %s", logger.sessionID, message)
}

func (logger *Logger) Printf(format string, v ...interface{}) {
	logger.Print(fmt.Sprintf(format, v...))
}

func (logger *Logger) PrintCommand(command string, params string) {
	if command == "PASS" {
		log.Printf("%s > PASS ****", logger.sessionID)
	} else {
		log.Printf("%s > %s %s", logger.sessionID, command, params)
	}
}

func (logger *Logger) PrintResponse(code int, message string) {
	log.Printf("%s < %d %s", logger.sessionID, code, message)
}
