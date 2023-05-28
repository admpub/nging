package ftp

import "errors"

var (
	ErrNotDirectory           = errors.New("Not a directory")
	ErrNotFile                = errors.New("Not a file")
	ErrDirectoryAlreadyExists = errors.New("A dir has the same name")
	ErrPutFile                = errors.New("Put File error")
)
