package goforever

import (
	"log"
	"os"
)

//NewLog Create a new file for logging
func NewLog(path string) *os.File {
	if path == "" {
		return nil
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatalf("%s\n", err)
		return nil
	}
	return file
}
