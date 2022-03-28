package pkg

import "errors"

var (
	ErrInvalidMediaPlaylist  = errors.New("not a valid media playlist")
	ErrInvalidMasterPlaylist = errors.New("not a valid master playlist")
	ErrExit                  = errors.New("exit")
	ErrContextCancelled      = errors.New("context cancelled")
	ErrPanic                 = errors.New("panic")
)
